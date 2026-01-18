package openaiv1

import "encoding/json"

type (
	Tool struct {
		Type     string        `json:"type"`
		Function *ToolFunction `json:"function,omitempty"`
	}

	ToolFunction struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description,omitempty"`
		Parameters  map[string]interface{} `json:"parameters,omitempty"`
	}

	ResponseUsage struct {
		InputTokens  int `json:"input_tokens,omitempty"`
		OutputTokens int `json:"output_tokens,omitempty"`
		TotalTokens  int `json:"total_tokens,omitempty"`
	}

	ResponseError struct {
		Type    string `json:"type,omitempty"`
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}

	ResponseOutputContent struct {
		Type string          `json:"type"`
		Text string          `json:"text,omitempty"`
		Data json.RawMessage `json:"data,omitempty"`
	}

	ResponseOutput struct {
		ID      string                  `json:"id,omitempty"`
		Type    string                  `json:"type,omitempty"`
		Status  string                  `json:"status,omitempty"`
		Role    string                  `json:"role,omitempty"`
		Content []ResponseOutputContent `json:"content,omitempty"`
	}

	Assistant struct {
		ID           string            `json:"id,omitempty"`
		Object       string            `json:"object,omitempty"`
		Name         string            `json:"name,omitempty"`
		Model        string            `json:"model,omitempty"`
		Instructions string            `json:"instructions,omitempty"`
		Tools        []Tool            `json:"tools,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		CreatedAt    int64             `json:"created_at,omitempty"`
	}

	Thread struct {
		ID        string            `json:"id,omitempty"`
		Object    string            `json:"object,omitempty"`
		Metadata  map[string]string `json:"metadata,omitempty"`
		CreatedAt int64             `json:"created_at,omitempty"`
	}

	ThreadMessage struct {
		ID        string            `json:"id,omitempty"`
		Object    string            `json:"object,omitempty"`
		Role      string            `json:"role,omitempty"`
		Content   json.RawMessage   `json:"content,omitempty"`
		Metadata  map[string]string `json:"metadata,omitempty"`
		CreatedAt int64             `json:"created_at,omitempty"`
	}

	Run struct {
		ID           string            `json:"id,omitempty"`
		Object       string            `json:"object,omitempty"`
		Status       string            `json:"status,omitempty"`
		AssistantID  string            `json:"assistant_id,omitempty"`
		ThreadID     string            `json:"thread_id,omitempty"`
		Model        string            `json:"model,omitempty"`
		Instructions string            `json:"instructions,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		Usage        ResponseUsage     `json:"usage,omitempty"`
		CreatedAt    int64             `json:"created_at,omitempty"`
	}
)
