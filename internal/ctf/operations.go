package ctf

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/utils"
	"net/http"

	"github.com/ctfer-io/go-ctfd/api"
)

type CustomPostChallengeParams struct {
	Slug string `json:"slug"`

	IsInstanced bool `json:"is_instanced"`
	HasOracle   bool `json:"has_oracle"`

	api.PostChallengesParams
}

type CustomPatchChallengeParams struct {
	Slug string `json:"slug"`

	IsInstanced bool `json:"is_instanced"`
	HasOracle   bool `json:"has_oracle"`

	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Description    string  `json:"description"`
	Function       string  `json:"function"`
	ConnectionInfo *string `json:"connection_info,omitempty"`
	Value          *int    `json:"value,omitempty"`
	Initial        *string `json:"initial,omitempty"`
	Decay          *string `json:"decay,omitempty"`
	Minimum        *string `json:"minimum,omitempty"`
	MaxAttempts    *string `json:"max_attempts,omitempty"`
	NextID         *int    `json:"next_id,omitempty"`
	// Requirements can update the challenge's behavior and prerequisites i.e.
	// the other challenges the team/user must have solved before.
	// WARNING: it won't return those in the response body, so updating this
	// field requires you to do it manually through *Client.GetChallengeRequirements
	Requirements *api.Requirements `json:"requirements,omitempty"`
	State        string            `json:"state"`
}

type CustomChallenge struct {
	Slug string `json:"slug"`

	IsInstanced bool `json:"is_instanced"`
	HasOracle   bool `json:"has_oracle"`

	api.Challenge
}

func (c *Client) GetChallenge(id int) (*CustomChallenge, error) {
	var ctfChallenge CustomChallenge
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/challenges/%d", id), nil)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	resp := api.Response{
		Data: &ctfChallenge,
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &ctfChallenge, nil
}

func (c *Client) GetChallenges() ([]*CustomChallenge, error) {
	var ctfChallenges []*CustomChallenge
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/challenges?view=admin", nil)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	resp := api.Response{
		Data: &ctfChallenges,
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return ctfChallenges, nil
}

func (c *Client) PostCustomChallenge(challenge *v1.ChallengeSpec) (*api.Challenge, error) {
	var apiChallenge api.Challenge

	body, err := json.Marshal(&CustomPostChallengeParams{
		Slug:        challenge.Slug,
		IsInstanced: challenge.IsInstanced,
		HasOracle:   challenge.HasOracle,
		PostChallengesParams: api.PostChallengesParams{
			Name:        challenge.Name,
			Category:    challenge.Category,
			Value:       challenge.Value,
			Description: challenge.Description,
			Initial:     challenge.Initial,
			Decay:       challenge.Decay,
			Function:    challenge.DecayFunction,
			Minimum:     challenge.Minimum,
			MaxAttempts: challenge.MaxAttempts,
			State:       challenge.State,
			Type:        challenge.Type,
		},
	})
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBuffer(body))
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	resp := api.Response{
		Data: &apiChallenge,
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &apiChallenge, nil
}

func (c *Client) PatchCustomChallenge(id int, challenge *v1.ChallengeSpec) (*api.Challenge, error) {
	var apiChallenge api.Challenge

	params := &CustomPatchChallengeParams{
		Slug:        challenge.Slug,
		IsInstanced: challenge.IsInstanced,
		HasOracle:   challenge.HasOracle,
		Name:        challenge.Name,
		Category:    challenge.Category,
		Description: challenge.Description,
		Initial:     utils.Optional(utils.SprintPtr(challenge.Initial)),
		Decay:       utils.Optional(utils.SprintPtr(challenge.Decay)),
		Minimum:     utils.Optional(utils.SprintPtr(challenge.Minimum)),
		MaxAttempts: utils.Optional(utils.SprintPtr(challenge.MaxAttempts)),
		Function:    challenge.DecayFunction,
		State:       challenge.State,
	}
	if challenge.Value != 0 {
		params.Value = utils.Optional(challenge.Value)
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/challenges/%d", id), bytes.NewBuffer(body))
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	resp := api.Response{
		Data: &apiChallenge,
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &apiChallenge, nil
}
