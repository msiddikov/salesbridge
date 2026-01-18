package openaiv1

import (
	"encoding/json"
)

type (
	AssistantRequest struct {
		Model        string            `json:"model"`
		Name         string            `json:"name,omitempty"`
		Description  string            `json:"description,omitempty"`
		Instructions string            `json:"instructions,omitempty"`
		Tools        []Tool            `json:"tools,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
	}

	AssistantUpdateRequest struct {
		Model        *string           `json:"model,omitempty"`
		Name         *string           `json:"name,omitempty"`
		Description  *string           `json:"description,omitempty"`
		Instructions *string           `json:"instructions,omitempty"`
		Tools        []Tool            `json:"tools,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
	}

	AssistantDeleteResponse struct {
		ID      string `json:"id,omitempty"`
		Object  string `json:"object,omitempty"`
		Deleted bool   `json:"deleted,omitempty"`
	}
)

func (c *Client) AssistantsCreate(req AssistantRequest) (Assistant, error) {
	res := Assistant{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/assistants",
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) AssistantsGet(assistantID string) (Assistant, error) {
	res := Assistant{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/assistants/" + assistantID,
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) AssistantsUpdate(assistantID string, req AssistantUpdateRequest) (Assistant, error) {
	res := Assistant{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/assistants/" + assistantID,
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) AssistantsDelete(assistantID string) (AssistantDeleteResponse, error) {
	res := AssistantDeleteResponse{}

	_, _, err := c.fetch(reqParams{
		Method:   "DELETE",
		Endpoint: "/assistants/" + assistantID,
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}
