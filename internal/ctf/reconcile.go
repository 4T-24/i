package ctf

import (
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/utils"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/sirupsen/logrus"
)

func (c *Client) ReconcileChallenge(challenges []*v1.ChallengeSpec) error {
	// Get all challenges
	ctfChallenges, err := c.GetChallenges(&api.GetChallengesParams{
		View: utils.Optional("admin"),
	})
	if err != nil {
		logrus.Error("Failed to get challenges from CTFd")
		return err
	}

	// Create a map of challenges
	challengesToDelete := make(map[int]bool)
	for _, challenge := range ctfChallenges {
		challengesToDelete[challenge.ID] = true
	}

	// Iterate over the challenges
	for _, challenge := range challenges {
		var found bool
		var foundID int
		for _, ctfChallenge := range ctfChallenges {
			if ctfChallenge.Name == challenge.Name && ctfChallenge.Type == challenge.Type {
				found = true
				foundID = ctfChallenge.ID
				delete(challengesToDelete, ctfChallenge.ID)
				break
			}
		}

		// If the challenge exists, update it
		if found {
			err := c.UpdateChallenge(foundID, challenge)
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
	for id := range challengesToDelete {
		err := c.DeleteChallenge(id)
		if err != nil {
			logrus.WithField("challenge_id", id).Error("Failed to delete challenge")
			return err
		}
	}

	return nil
}

func (c *Client) InsertChallenge(challenge *v1.ChallengeSpec) error {
	// If the challenge does not exist, create it
	var requirements *api.Requirements
	if len(challenge.Requirements.Prerequisites) > 0 {
		requirements = &api.Requirements{
			Anonymize:     challenge.Requirements.Anonymize,
			Prerequisites: challenge.Requirements.Prerequisites,
		}
	}

	apiChallenge, err := c.PostChallenges(&api.PostChallengesParams{
		Name:         challenge.Name,
		Category:     challenge.Category,
		Value:        challenge.Value,
		Description:  challenge.Description,
		Initial:      challenge.Initial,
		Decay:        challenge.Decay,
		Function:     challenge.DecayFunction,
		Minimum:      challenge.Minimum,
		MaxAttempts:  challenge.MaxAttempts,
		NextID:       challenge.NextID,
		Requirements: requirements,
		State:        challenge.State,
		Type:         challenge.Type,
	})
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create challenge")
		return err
	}

	// We need to create or update the flag
	_, err = c.PostFlags(&api.PostFlagsParams{
		Challenge: apiChallenge.ID,
		Content:   challenge.Flag,
		Type:      "static",
	})
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create flag")
		return err
	}

	// We need to create or update the hints

	for _, hint := range challenge.Hints {
		var requirements api.Requirements
		if hint.Requirements != nil {
			requirements = api.Requirements(*hint.Requirements)
		}

		_, err = c.PostHints(&api.PostHintsParams{
			ChallengeID:  apiChallenge.ID,
			Content:      hint.Content,
			Cost:         hint.Cost,
			Requirements: requirements,
		})
		if err != nil {
			logrus.WithField("challenge_name", challenge.Name).Error("Failed to create hint")
			return err
		}
	}

	return nil
}

func (c *Client) UpdateChallenge(id int, challenge *v1.ChallengeSpec) error {
	// If the challenge exists, update it
	var requirements *api.Requirements
	if len(challenge.Requirements.Prerequisites) > 0 {
		requirements = &api.Requirements{
			Anonymize:     challenge.Requirements.Anonymize,
			Prerequisites: challenge.Requirements.Prerequisites,
		}
	}

	_, err := c.PatchChallenge(id, &api.PatchChallengeParams{
		Name:         challenge.Name,
		Category:     challenge.Category,
		Value:        utils.Optional(challenge.Value),
		Description:  challenge.Description,
		Initial:      challenge.Initial,
		Decay:        challenge.Decay,
		Function:     challenge.DecayFunction,
		Minimum:      challenge.Minimum,
		MaxAttempts:  challenge.MaxAttempts,
		NextID:       challenge.NextID,
		Requirements: requirements,
		State:        challenge.State,
	})
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to update challenge")
		return err
	}

	// Delete the flags and hints and recreate them
	apiFlags, err := c.GetFlags(&api.GetFlagsParams{
		ChallengeID: utils.Optional(id),
	})
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to get flags")
		return err
	}

	for _, flag := range apiFlags {
		err = c.DeleteFlag(fmt.Sprint(flag.ID))
		if err != nil {
			logrus.WithField("flag_id", flag.ID).Error("Failed to delete flag")
			return err
		}
	}

	apiHints, err := c.GetHints(&api.GetHintsParams{
		ChallengeID: utils.Optional(id),
	})
	if err != nil {
		logrus.WithField("challenge_id", id).Error("Failed to get hints")
		return err
	}

	for _, hint := range apiHints {
		err = c.DeleteHint(fmt.Sprint(hint.ID))
		if err != nil {
			logrus.WithField("hint_id", hint.ID).Error("Failed to delete hint")
			return err
		}
	}

	// We need to create or update the flag and hints
	err = c.CreateFlagAndHints(id, challenge)
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create flag and hints")
		return err
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

	// Delete hints
	hints, err := c.GetHints(&api.GetHintsParams{
		ChallengeID: utils.Optional(id),
	})
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

func (c *Client) CreateFlagAndHints(challengeID int, challenge *v1.ChallengeSpec) error {
	// We need to create or update the flag
	_, err := c.PostFlags(&api.PostFlagsParams{
		Challenge: challengeID,
		Content:   challenge.Flag,
		Type:      "static",
	})
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to create flag")
		return err
	}

	// We need to create or update the hints
	for _, hint := range challenge.Hints {
		_, err = c.PostHints(&api.PostHintsParams{
			ChallengeID: challengeID,
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
