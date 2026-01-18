package openaiv1

import "encoding/json"

type (
	ResponseRequest struct {
		Model           string         `json:"model"`
		Input           interface{}    `json:"input,omitempty"`
		Instructions    string         `json:"instructions,omitempty"`
		Tools           []Tool         `json:"tools,omitempty"`
		ToolChoice      interface{}    `json:"tool_choice,omitempty"`
		Temperature     *float64       `json:"temperature,omitempty"`
		MaxOutputTokens int            `json:"max_output_tokens,omitempty"`
		Metadata        map[string]any `json:"metadata,omitempty"`
		Stream          bool           `json:"stream,omitempty"`
	}

	Response struct {
		ID        string           `json:"id,omitempty"`
		Object    string           `json:"object,omitempty"`
		Status    string           `json:"status,omitempty"`
		Model     string           `json:"model,omitempty"`
		CreatedAt int64            `json:"created_at,omitempty"`
		Output    []ResponseOutput `json:"output,omitempty"`
		Usage     ResponseUsage    `json:"usage,omitempty"`
		Error     *ResponseError   `json:"error,omitempty"`
	}
)

func (c *Client) ResponsesCreate(req ResponseRequest) (Response, error) {
	res := Response{}
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/responses",
		Body:     string(payload),
	}, &res)

	return res, err
}
