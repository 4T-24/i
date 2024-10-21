package ctf

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/files"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/sirupsen/logrus"
)

func (c *Client) UpdateChallenge(id int, challenge *v1.ChallengeSpec) error {
	_, err := c.PatchCustomChallenge(id, challenge)
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to update challenge")
		return err
	}

	err = c.VerifyChallengeFlag(id, challenge)
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to verify flag")
		return err
	}

	if len(challenge.Files) > 0 {
		err = c.VerifyChallengeFiles(id, challenge)
		if err != nil {
			logrus.WithField("challenge_name", challenge.Name).Error("Failed to verify files")
		}
	}

	if len(challenge.Hints) > 0 {
		err = c.VerifyChallengeHints(id, challenge)
		if err != nil {
			logrus.WithField("challenge_name", challenge.Name).Error("Failed to verify hints")
			return err
		}
	}

	return nil
}

func (c *Client) VerifyChallengeFlag(id int, challenge *v1.ChallengeSpec) error {
	apiFlags, err := c.GetChallengeFlags(id)
	if err != nil {
		return err
	}

	for _, flag := range apiFlags {
		if flag.Content == challenge.Flag {
			return nil
		}

		err = c.DeleteFlag(fmt.Sprint(flag.ID))
		if err != nil {
			return err
		}
	}

	_, err = c.PostFlags(&api.PostFlagsParams{
		Challenge: id,
		Content:   challenge.Flag,
		Type:      "static",
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) VerifyChallengeFiles(id int, challenge *v1.ChallengeSpec) error {
	apiFiles, err := c.GetChallengeFiles(id)
	if err != nil {
		return err
	}

	var filePaths []string
	for _, f := range challenge.Files {
		filePaths = append(filePaths, f.Path)
	}

	fileDatas, err := files.GetFiles(challenge.Repository, filePaths)
	if err != nil {
		return err
	}
	var fileSums = make(map[string]string)
	var foundFiles = make(map[string]bool)
	for i, f := range challenge.Files {
		sha := sha1.Sum(fileDatas[i])
		fileSums[hex.EncodeToString(sha[:])] = f.Name
		foundFiles[f.Name] = false
	}

	for _, file := range apiFiles {
		if name, found := fileSums[file.SHA1sum]; found {
			foundFiles[name] = true
			continue
		}

		err = c.DeleteFile(fmt.Sprint(file.ID))
		if err != nil {
			return err
		}
	}

	for i, f := range challenge.Files {
		if foundFiles[f.Name] {
			continue
		}

		_, err = c.PostFiles(&api.PostFilesParams{
			Files: []*api.InputFile{
				{
					Name:    f.Name,
					Content: fileDatas[i],
				},
			},
			Challenge: &id,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) VerifyChallengeHints(id int, challenge *v1.ChallengeSpec) error {
	apiHints, err := c.GetChallengeHints(id)
	if err != nil {
		return err
	}

	for _, hint := range apiHints {
		// Check if we should delete the hint or keep it
		var hintExists bool
		for _, newHint := range challenge.Hints {
			if *hint.Content == newHint.Content {
				hintExists = true
				continue
			}
		}

		if hintExists {
			continue
		}

		err = c.DeleteHint(fmt.Sprint(hint.ID))
		if err != nil {
			return err
		}
	}

	for _, newHint := range challenge.Hints {
		var hintExists bool
		for _, hint := range apiHints {
			if *hint.Content == newHint.Content {
				hintExists = true
				continue
			}
		}

		if hintExists {
			continue
		}

		_, err = c.PostHints(&api.PostHintsParams{
			ChallengeID: id,
			Content:     newHint.Content,
			Cost:        newHint.Cost,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
