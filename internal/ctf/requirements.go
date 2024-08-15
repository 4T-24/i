package ctf

import (
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/utils"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/sirupsen/logrus"
)

func (c *Client) ReconcileRequirements(challenges map[string]*v1.ChallengeSpec) error {
	// Challenges are mapped my their slug, and requirements references the slug of the required challenge
	apiChallenges, err := c.GetChallenges()
	if err != nil {
		return fmt.Errorf("failed to get challenges: %w", err)
	}

	// Create a map of challenges
	mappedApiChallenges := make(map[string]*CustomChallenge)
	for _, challenge := range apiChallenges {
		if challenge.Type == "i_static" || challenge.Type == "i_dynamic" {
			challenge, err := c.GetChallenge(challenge.ID)
			if err != nil {
				return fmt.Errorf("failed to get challenge %d: %w", challenge.ID, err)
			}

			mappedApiChallenges[challenge.Slug] = challenge
		}
	}

	for _, challenge := range challenges {
		apiChallenge, found := mappedApiChallenges[challenge.Slug]
		if !found {
			return fmt.Errorf("challenge %s does not exist", challenge.Name)
		}

		var prerequisites []int

		for _, prerequisite := range challenge.Requirements.Prerequisites {
			prerequisiteChallenge, found := challenges[prerequisite]
			if !found {
				return fmt.Errorf("challenge %s has a prerequisite %s that does not exist", challenge.Name, prerequisite)
			}

			prerequisiteApiChallenge, found := mappedApiChallenges[prerequisiteChallenge.Slug]
			if !found {
				return fmt.Errorf("challenge %s has a prerequisite %s that does not exist", challenge.Name, prerequisite)
			}

			prerequisites = append(prerequisites, prerequisiteApiChallenge.ID)
		}

		c.PatchChallenge(apiChallenge.ID, &api.PatchChallengeParams{
			Name:        apiChallenge.Name,
			Category:    apiChallenge.Category,
			Description: apiChallenge.Description,
			Function:    apiChallenge.Function,
			Requirements: &api.Requirements{
				Prerequisites: prerequisites,
			},
			State: apiChallenge.State,
		})

		// Now proceed to do requirements for hints
		apiHints, err := c.GetChallengeHints(apiChallenge.ID)
		if err != nil {
			return err
		}

		// Create a map of hints for the challenge
		mappedApiHints := make(map[string]*api.Hint)
		for _, hint := range apiHints {
			mappedApiHints[*hint.Content] = hint
		}

		for _, hint := range challenge.Hints {
			if hint.Requirements == nil {
				hint.Requirements = &v1.HintRequirements{}
			}

			apiHint, found := mappedApiHints[hint.Content]
			if !found {
				return fmt.Errorf("hint %s does not exist", hint.Content)
			}

			prerequisites = []int{}

			for _, prerequisite := range hint.Requirements.Prerequisites {
				if prerequisite < 0 || prerequisite >= len(challenge.Hints) {
					return fmt.Errorf("hint %s has a prerequisite %d that does not exist", hint.Content, prerequisite)
				}

				prerequisiteHint, found := mappedApiHints[challenge.Hints[prerequisite].Content]
				if !found {
					return fmt.Errorf("hint %s has a prerequisite %s that does not exist", hint.Content, challenge.Hints[prerequisite].Content)
				}

				prerequisites = append(prerequisites, prerequisiteHint.ID)
			}

			c.PatchHint(fmt.Sprint(apiHint.ID), &api.PatchHintsParams{
				ChallengeID: apiChallenge.ID,
				Content:     hint.Content,
				Cost:        hint.Cost,
				Requirements: api.Requirements{
					Prerequisites: prerequisites,
				},
			})
		}
	}

	for _, challenge := range challenges {
		apiChallenge, found := mappedApiChallenges[challenge.Slug]
		if !found {
			return fmt.Errorf("challenge %s does not exist", challenge.Name)
		}

		var nextId *int
		nextChallenge, found := mappedApiChallenges[challenge.NextSlug]
		if !found {
			if challenge.NextSlug != "" {
				logrus.WithField("slug", challenge.NextSlug).Warn("next challenge does not exist")
			}
			continue
		}
		nextId = utils.Optional(nextChallenge.ID)

		c.PatchChallenge(apiChallenge.ID, &api.PatchChallengeParams{
			Name:        apiChallenge.Name,
			Category:    apiChallenge.Category,
			Description: apiChallenge.Description,
			Function:    apiChallenge.Function,
			State:       apiChallenge.State,
			NextID:      nextId,
		})
	}

	return nil
}
