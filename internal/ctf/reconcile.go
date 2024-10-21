package ctf

import (
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/files"
	"instancer/internal/utils"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/sirupsen/logrus"
)

func (c *Client) ReconcileChallenge(challenges map[string]*v1.ChallengeSpec) error {
	ctfChallenges, err := c.GetChallenges()
	if err != nil {
		logrus.WithError(err).Error("Failed to get challenges")
		return err
	}

	// Create a map of challenges
	challengesToDelete := make(map[string]*CustomChallenge)
	for _, challenge := range ctfChallenges {
		if challenge.Type == "i_static" || challenge.Type == "i_dynamic" {
			ctfChallenge, err := c.GetChallenge(challenge.ID)
			if err != nil {
				logrus.WithError(err).WithField("challenge_id", challenge.ID).Error("Failed to get challenge")
				return err
			}

			challengesToDelete[ctfChallenge.Slug] = ctfChallenge
		}
	}

	// Iterate over the challenges
	for _, challenge := range challenges {
		ctfChallenge, found := challengesToDelete[challenge.Slug]

		// If the challenge exists, update it
		if found {
			delete(challengesToDelete, ctfChallenge.Slug)
			err := c.UpdateChallenge(ctfChallenge.ID, challenge)
			if err != nil {
				logrus.WithField("challenge_name", challenge.Name).Error("Failed to update challenge")
				return err
			}
		} else {
			// If the challenge does not exist, create it
			err := c.InsertChallenge(challenge)
			if err != nil {
				logrus.WithField("challenge_name", challenge.Name).Error("Failed to create challenge")
				return err
			}
		}
	}

	// Delete challenges that are not in the list
	for _, challenge := range challengesToDelete {
		err := c.DeleteChallenge(challenge.ID)
		if err != nil {
			logrus.WithField("challenge_id", challenge.ID).Error("Failed to delete challenge")
			return err
		}
	}

	// Reconcile requirements
	err = c.ReconcileRequirements(challenges)
	if err != nil {
		logrus.WithError(err).Error("Failed to reconcile requirements")
		return err
	}

	return nil
}

func (c *Client) InsertChallenge(challenge *v1.ChallengeSpec) error {
	apiChallenge, err := c.PostCustomChallenge(challenge)
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create challenge")
		return err
	}

	// We need to create the flag
	_, err = c.PostFlags(&api.PostFlagsParams{
		Challenge: apiChallenge.ID,
		Content:   challenge.Flag,
		Type:      "static",
	})
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create flag")
		return err
	}

	// We need to create files
	if len(challenge.Files) > 0 {
		var filePaths []string
		for _, f := range challenge.Files {
			filePaths = append(filePaths, f.Path)
		}

		fileDatas, err := files.GetFiles(challenge.Repository, filePaths)
		if err != nil {
			logrus.WithError(err).WithField("challenge_name", challenge.Name).Error("Failed to download files from git")
		}

		var apiFiles []*api.InputFile
		for i, f := range challenge.Files {
			apiFiles = append(apiFiles, &api.InputFile{
				Name:    f.Name,
				Content: fileDatas[i],
			})
		}

		c.PostFiles(&api.PostFilesParams{
			Files:     apiFiles,
			Challenge: &apiChallenge.ID,
		})
	}

	// We need to create hints
	for _, hint := range challenge.Hints {
		_, err = c.PostHints(&api.PostHintsParams{
			ChallengeID: apiChallenge.ID,
			Content:     hint.Content,
			Cost:        hint.Cost,
		})
		if err != nil {
			logrus.WithField("challenge_name", challenge.Name).Error("Failed to create hint")
			return err
		}
	}

	return nil
}

func (c *Client) DeleteChallenge(id int) error {
	// Delete flags
	flags, err := c.GetFlags(&api.GetFlagsParams{
		ChallengeID: utils.Optional(id),
	})
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to get flags")
		return err
	}
	for _, flag := range flags {
		err = c.DeleteFlag(fmt.Sprint(flag.ID))
		if err != nil {
			logrus.WithField("flag_id", flag.ID).Error("Failed to delete flag")
			return err
		}
	}

	// Delete files
	files, err := c.GetChallengeFiles(id)
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to get files")
		return err
	}
	for _, file := range files {
		err = c.DeleteFile(fmt.Sprint(file.ID))
		if err != nil {
			logrus.WithField("flag_id", file.ID).Error("Failed to delete file")
			return err
		}
	}

	// Delete hints
	hints, err := c.GetChallengeHints(id)
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to get hints")
		return err
	}
	for _, hint := range hints {
		err = c.DeleteHint(fmt.Sprint(hint.ID))
		if err != nil {
			logrus.WithField("hint_id", hint.ID).Error("Failed to delete hint")
			return err
		}
	}

	err = c.Client.DeleteChallenge(id)
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to delete challenge")
		return err
	}

	return nil
}
