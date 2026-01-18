package openaiv1

import (
	"encoding/json"
	"strconv"
)

type (
	ThreadRequest struct {
		Metadata map[string]string `json:"metadata,omitempty"`
	}

	ThreadMessageRequest struct {
		Role     string            `json:"role"`
		Content  interface{}       `json:"content"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}

	RunRequest struct {
		AssistantID  string            `json:"assistant_id"`
		Model        string            `json:"model,omitempty"`
		Instructions string            `json:"instructions,omitempty"`
		Tools        []Tool            `json:"tools,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
	}

	RunToolOutput struct {
		ToolCallID string `json:"tool_call_id"`
		Output     string `json:"output"`
	}

	RunSubmitToolOutputsRequest struct {
		ToolOutputs []RunToolOutput `json:"tool_outputs"`
		Stream      bool            `json:"stream,omitempty"`
	}

	ThreadMessageContentText struct {
		Value string `json:"value"`
	}

	ThreadMessageContent struct {
		Type string                    `json:"type"`
		Text *ThreadMessageContentText `json:"text,omitempty"`
	}

	ThreadMessageResponse struct {
		ID        string                 `json:"id,omitempty"`
		Object    string                 `json:"object,omitempty"`
		Role      string                 `json:"role,omitempty"`
		Content   []ThreadMessageContent `json:"content,omitempty"`
		CreatedAt int64                  `json:"created_at,omitempty"`
	}

	ThreadMessagesListResponse struct {
		Data []ThreadMessageResponse `json:"data,omitempty"`
	}
)

func (c *Client) ThreadsCreate(req ThreadRequest) (Thread, error) {
	res := Thread{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/threads",
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) ThreadsMessagesCreate(threadID string, req ThreadMessageRequest) (ThreadMessage, error) {
	res := ThreadMessage{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/threads/" + threadID + "/messages",
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) RunsCreate(threadID string, req RunRequest) (Run, error) {
	res := Run{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/threads/" + threadID + "/runs",
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) RunsGet(threadID, runID string) (Run, error) {
	res := Run{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/threads/" + threadID + "/runs/" + runID,
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) RunsSubmitToolOutputs(threadID, runID string, req RunSubmitToolOutputsRequest) (Run, error) {
	res := Run{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/threads/" + threadID + "/runs/" + runID + "/submit_tool_outputs",
		Body:     string(payload),
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}

func (c *Client) ThreadsMessagesList(threadID string, limit int) (ThreadMessagesListResponse, error) {
	res := ThreadMessagesListResponse{}

	params := []queryParam{{Key: "order", Value: "desc"}}
	if limit > 0 {
		params = append(params, queryParam{Key: "limit", Value: strconv.Itoa(limit)})
	}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/threads/" + threadID + "/messages",
		QParams:  params,
		Headers: map[string]string{
			"OpenAI-Beta": "assistants=v2",
		},
	}, &res)

	return res, err
}
