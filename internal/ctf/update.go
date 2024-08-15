package ctf

import (
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/utils"

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

	err = c.VerifyChallengeHints(id, challenge)
	if err != nil {
		logrus.WithField("challenge_name", challenge.Name).Error("Failed to verify hints")
		return err
	}

	return nil
}

func (c *Client) VerifyChallengeFlag(id int, challenge *v1.ChallengeSpec) error {
	apiFlags, err := c.GetFlags(&api.GetFlagsParams{
		ChallengeID: utils.Optional(id),
	})
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
