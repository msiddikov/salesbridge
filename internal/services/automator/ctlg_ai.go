package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_openai"
	"context"
	"fmt"
	"strings"
)

var (
	aiCategory = Category{
		Id:    "ai",
		Name:  "AI",
		Icon:  "ri:robot-line",
		Color: "#10B981",
		Nodes: []Node{
			aiAssistantAction,
		},
	}

	aiAssistantAction = Node{
		Id:          "ai.assistant",
		Title:       "AI Assistant",
		Description: "Runs a message through a selected AI assistant.",
		ExecFunc:    aiAssistantExecute,
		Type:        NodeTypeAction,
		Icon:        "ri:robot-line",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(aiAssistantOutputFields),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "assistant_id", Type: "string", Required: true},
			{Key: "text", Type: "string", Required: true},
		},
	}

	aiAssistantOutputFields = []NodeField{
		{Key: "message", Type: "string"},
	}
)

func aiAssistantExecute(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	text := strings.TrimSpace(fmt.Sprint(fields["text"]))
	if text == "" {
		return errorPayload(nil, "text is required")
	}

	rawID := strings.TrimSpace(fmt.Sprint(fields["assistant_id"]))
	if rawID == "" || rawID == "<nil>" {
		return errorPayload(nil, "assistant_id is required")
	}

	responseText, err := svc_openai.RunThroughAssistant(rawID, text)
	if err != nil {
		return errorPayload(err, "assistant run failed")
	}

	return successPayload(map[string]interface{}{"message": responseText})
}
