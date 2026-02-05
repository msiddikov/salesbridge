package svc_internal_assistant

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	openaiv1 "client-runaway-zenoti/packages/openaiV1"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mcpclient "github.com/mark3labs/mcp-go/client"
	mcptransport "github.com/mark3labs/mcp-go/client/transport"
	mcptypes "github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
)

// getOrchestratorAssistantID returns the orchestrator assistant ID from config
func getOrchestratorAssistantID() (string, error) {
	settings := config.Confs.Settings

	// Try orchestrator ID first
	id := strings.TrimSpace(settings.InternalOrchestratorID)

	// Fallback to legacy InternalAssistantID
	if id == "" {
		id = strings.TrimSpace(settings.InternalAssistantID)
	}

	if id == "" {
		return "", fmt.Errorf("orchestrator assistant ID not configured")
	}

	return id, nil
}

type streamRequest struct {
	ThreadID     string        `json:"thread_id,omitempty"`
	Message      string        `json:"message"`
	Instructions string        `json:"instructions,omitempty"`
	Context      []contextItem `json:"context,omitempty"`
}

type streamEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type contextItem struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}

type threadResponse struct {
	ID          uint   `json:"id"`
	ThreadID    string `json:"thread_id"`
	AssistantID string `json:"assistant_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type threadMessageResponse struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type agentInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListAgents returns the configured agent assistant IDs for the orchestrator to use
func ListAgents(c *gin.Context) {
	settings := config.Confs.Settings

	agents := []agentInfo{
		{
			ID:          strings.TrimSpace(settings.InternalOrchestratorID),
			Name:        "Orchestrator",
			Description: "Routes requests to the appropriate agent based on the task",
		},
		{
			ID:          strings.TrimSpace(settings.InternalBuilderID),
			Name:        "Builder",
			Description: "Helps build and configure automations",
		},
		{
			ID:          strings.TrimSpace(settings.InternalHelperID),
			Name:        "Helper",
			Description: "Answers questions and finds information in the workspace",
		},
	}

	c.JSON(http.StatusOK, agents)
}

func Stream(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	req := streamRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	// All requests go to the orchestrator, which decides which agent to use
	assistantID, err := getOrchestratorAssistantID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	mcpKey, err := getInternalMCPKey(user.ProfileID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load internal MCP key"})
		return
	}
	if mcpKey == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "mcp internal api key is not configured"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	sendEvent(c, "status", streamEvent{Type: "starting"})

	openaiClient, err := newOpenAIClient()
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "openai_client_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	threadID, err := resolveInternalAssistantThreadID(openaiClient, user.ProfileID, assistantID, req.ThreadID)
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "thread_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	sendEvent(c, "thread", streamEvent{
		Type: "thread",
		Data: map[string]any{
			"thread_id": threadID,
		},
	})

	mcpURL, err := resolveMCPURL()
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "mcp_url_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	mcpClient, err := newMCPClient(ctx, mcpURL, mcpKey)
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "mcp_client_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	tools, err := listMCPTools(ctx, mcpClient)
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "mcp_tools_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	runInstructions := buildRunInstructions(req.Instructions, req.Context)
	message, err := runAssistant(ctx, openaiClient, mcpClient, assistantID, threadID, req.Message, runInstructions, tools, func(eventType string, payload any) {
		sendEvent(c, eventType, payload)
	})
	if err != nil {
		sendEvent(c, "error", streamEvent{Type: "run_error", Data: err.Error()})
		sendEvent(c, "done", streamEvent{Type: "done"})
		return
	}

	sendEvent(c, "message", streamEvent{Type: "final", Data: message})
	sendEvent(c, "done", streamEvent{Type: "done"})
}

func getInternalMCPKey(profileID uint) (string, error) {
	if profileID == 0 {
		return "", fmt.Errorf("profile id is required")
	}

	var key models.MCPApiKey
	err := db.DB.
		Where("profile_id = ? AND is_internal = ? AND is_active = ?", profileID, true, true).
		Order("id DESC").
		First(&key).Error

	if err == nil {
		plain := strings.TrimSpace(key.PlainKey)
		if plain != "" {
			return plain, nil
		}
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	plainKey, keyHash, keyPrefix, err := models.GenerateMCPApiKey()
	if err != nil {
		return "", err
	}

	if key.ID != 0 {
		err = db.DB.Model(&models.MCPApiKey{}).
			Where("id = ?", key.ID).
			Updates(map[string]any{
				"key_hash":    keyHash,
				"key_prefix":  keyPrefix,
				"plain_key":   plainKey,
				"is_active":   true,
				"is_internal": true,
			}).Error
		if err != nil {
			return "", err
		}
	} else {
		newKey := models.MCPApiKey{
			Name:       fmt.Sprintf("internal-profile-%d", profileID),
			PlainKey:   plainKey,
			KeyHash:    keyHash,
			KeyPrefix:  keyPrefix,
			ProfileID:  profileID,
			IsActive:   true,
			IsInternal: true,
		}
		if err := db.DB.Create(&newKey).Error; err != nil {
			return "", err
		}
	}

	return plainKey, nil
}

func sendEvent(c *gin.Context, event string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	fmt.Fprintf(c.Writer, "event: %s\n", event)
	for _, line := range strings.Split(string(data), "\n") {
		fmt.Fprintf(c.Writer, "data: %s\n", line)
	}
	fmt.Fprint(c.Writer, "\n")

	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

func resolveInternalAssistantThreadID(client openaiv1.Client, profileID uint, assistantID string, requestedThreadID string) (string, error) {
	if profileID == 0 || assistantID == "" {
		return "", fmt.Errorf("profile id and assistant id are required")
	}

	requestedThreadID = strings.TrimSpace(requestedThreadID)
	if requestedThreadID != "" {
		var existing models.InternalAssistantThread
		err := db.DB.
			Where("profile_id = ? AND assistant_id = ? AND thread_id = ?", profileID, assistantID, requestedThreadID).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", fmt.Errorf("thread not found")
			}
			return "", err
		}
		_ = db.DB.Model(&existing).Update("updated_at", time.Now()).Error
		return existing.ThreadID, nil
	}

	created, err := client.ThreadsCreate(openaiv1.ThreadRequest{})
	if err != nil {
		return "", err
	}

	newThread := models.InternalAssistantThread{
		ProfileID:   profileID,
		AssistantID: assistantID,
		ThreadID:    created.ID,
	}
	if err := db.DB.Create(&newThread).Error; err != nil {
		return "", err
	}

	return created.ID, nil
}

func ListThreads(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var threads []models.InternalAssistantThread
	if err := db.DB.
		Where("profile_id = ?", user.ProfileID).
		Order("updated_at desc").
		Find(&threads).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to list threads"})
		return
	}

	resp := make([]threadResponse, 0, len(threads))
	for _, t := range threads {
		resp = append(resp, threadResponse{
			ID:          t.ID,
			ThreadID:    t.ThreadID,
			AssistantID: t.AssistantID,
			CreatedAt:   t.CreatedAt.Unix(),
			UpdatedAt:   t.UpdatedAt.Unix(),
		})
	}

	c.JSON(http.StatusOK, resp)
}

func ListThreadMessages(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	threadID := strings.TrimSpace(c.Param("threadId"))
	if threadID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "thread_id is required"})
		return
	}

	var thread models.InternalAssistantThread
	if err := db.DB.
		Where("profile_id = ? AND thread_id = ?", user.ProfileID, threadID).
		First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "thread not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load thread"})
		return
	}

	limit := 0
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed < 1 || parsed > 200 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 200"})
			return
		}
		limit = parsed
	}

	openaiClient, err := newOpenAIClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "openai client error"})
		return
	}

	messages, err := openaiClient.ThreadsMessagesList(threadID, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load messages"})
		return
	}

	resp := make([]threadMessageResponse, 0, len(messages.Data))
	for i := len(messages.Data) - 1; i >= 0; i-- {
		msg := messages.Data[i]
		resp = append(resp, threadMessageResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   threadMessageText(msg),
			CreatedAt: msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func newOpenAIClient() (openaiv1.Client, error) {
	settings := config.Confs.Settings
	if settings.OpenAIInternalAPIKey == "" {
		return openaiv1.Client{}, fmt.Errorf("internal ai api key is required")
	}

	service := openaiv1.Service{
		APIKey:       settings.OpenAIInternalAPIKey,
		BaseURL:      settings.OpenAIBaseURL,
		Organization: settings.OpenAIOrganization,
		Project:      settings.OpenAIProject,
	}

	return service.NewClient("")
}

func resolveMCPURL() (string, error) {
	raw := strings.TrimSpace(config.Confs.Settings.MCPURL)
	if raw == "" {
		return "", fmt.Errorf("mcp url must be set in config")
	}
	return raw, nil
}

func newMCPClient(ctx context.Context, mcpURL, apiKey string) (*mcpclient.Client, error) {
	headers := map[string]string{}
	if apiKey != "" {
		headers["Authorization"] = "Bearer " + apiKey
	}

	transport, err := mcptransport.NewStreamableHTTP(
		mcpURL,
		mcptransport.WithHTTPHeaders(headers),
		mcptransport.WithHTTPTimeout(60*time.Second),
	)
	if err != nil {
		return nil, err
	}

	client := mcpclient.NewClient(transport)
	if err := client.Start(ctx); err != nil {
		return nil, err
	}

	_, err = client.Initialize(ctx, mcptypes.InitializeRequest{
		Params: mcptypes.InitializeParams{
			ProtocolVersion: mcptypes.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcptypes.Implementation{
				Name:    "internal-assistant",
				Version: "1.0.0",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func listMCPTools(ctx context.Context, client *mcpclient.Client) ([]mcptypes.Tool, error) {
	result, err := client.ListTools(ctx, mcptypes.ListToolsRequest{})
	if err != nil {
		return nil, err
	}
	return result.Tools, nil
}

func buildRunInstructions(extra string, ctxItems []contextItem) string {
	base := strings.TrimSpace(``)
	contextBlock := formatContext(ctxItems)
	extra = strings.TrimSpace(extra)

	parts := []string{base}
	if contextBlock != "" {
		parts = append(parts, "Context:\n"+contextBlock)
	}
	if extra != "" {
		parts = append(parts, extra)
	}
	return strings.Join(parts, "\n\n")
}

func formatContext(items []contextItem) string {
	if len(items) == 0 {
		return ""
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		value := "null"
		if item.Value != nil {
			if raw, err := json.Marshal(item.Value); err == nil {
				value = string(raw)
			} else {
				value = fmt.Sprint(item.Value)
			}
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", name, value))
	}
	return strings.Join(lines, "\n")
}

func runAssistant(
	ctx context.Context,
	client openaiv1.Client,
	mcpClient *mcpclient.Client,
	assistantID string,
	threadID string,
	message string,
	instructions string,
	tools []mcptypes.Tool,
	onEvent func(eventType string, payload any),
) (string, error) {
	openaiTools, err := toOpenAITools(tools)
	if err != nil {
		return "", err
	}

	threadID = strings.TrimSpace(threadID)
	if threadID == "" {
		thread, err := client.ThreadsCreate(openaiv1.ThreadRequest{})
		if err != nil {
			return "", err
		}
		threadID = thread.ID
	}

	_, err = client.ThreadsMessagesCreate(threadID, openaiv1.ThreadMessageRequest{
		Role:    "user",
		Content: message,
	})
	if err != nil {
		return "", err
	}

	run, err := client.RunsCreate(threadID, openaiv1.RunRequest{
		AssistantID:  assistantID,
		Instructions: instructions,
		Tools:        openaiTools,
	})
	if err != nil {
		return "", err
	}

	run, err = waitForRun(ctx, client, mcpClient, threadID, run, onEvent)
	if err != nil {
		return "", err
	}

	messages, err := client.ThreadsMessagesList(threadID, 20)
	if err != nil {
		return "", err
	}

	return firstAssistantMessage(messages)
}

func waitForRun(
	ctx context.Context,
	client openaiv1.Client,
	mcpClient *mcpclient.Client,
	threadID string,
	run openaiv1.Run,
	onEvent func(eventType string, payload any),
) (openaiv1.Run, error) {
	for i := 0; i < 60; i++ {
		switch run.Status {
		case "completed":
			return run, nil
		case "failed", "cancelled", "expired":
			return run, fmt.Errorf("run %s", run.Status)
		case "requires_action":
			if run.RequiredAction == nil || run.RequiredAction.Type != "submit_tool_outputs" || run.RequiredAction.SubmitToolOutputs == nil {
				return run, fmt.Errorf("unsupported required_action")
			}
			outputs, err := handleToolCalls(ctx, mcpClient, run.RequiredAction.SubmitToolOutputs.ToolCalls, onEvent)
			if err != nil {
				return run, err
			}
			run, err = client.RunsSubmitToolOutputs(threadID, run.ID, openaiv1.RunSubmitToolOutputsRequest{
				ToolOutputs: outputs,
			})
			if err != nil {
				return run, err
			}
		default:
			time.Sleep(2 * time.Second)
			updated, err := client.RunsGet(threadID, run.ID)
			if err != nil {
				return run, err
			}
			run = updated
		}
	}

	return run, fmt.Errorf("run timed out")
}

func handleToolCalls(
	ctx context.Context,
	mcpClient *mcpclient.Client,
	toolCalls []openaiv1.RunToolCall,
	onEvent func(eventType string, payload any),
) ([]openaiv1.RunToolOutput, error) {
	outputs := make([]openaiv1.RunToolOutput, 0, len(toolCalls))
	for _, call := range toolCalls {
		if call.Type != "function" || call.Function == nil {
			return nil, fmt.Errorf("unsupported tool call")
		}

		args := map[string]any{}
		if strings.TrimSpace(call.Function.Arguments) != "" {
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				return nil, fmt.Errorf("invalid tool arguments: %w", err)
			}
		}

		onEvent("tool_call", streamEvent{
			Type: "tool_call",
			Data: map[string]any{
				"name":      call.Function.Name,
				"arguments": args,
			},
		})

		result, err := mcpClient.CallTool(ctx, mcptypes.CallToolRequest{
			Params: mcptypes.CallToolParams{
				Name:      call.Function.Name,
				Arguments: args,
			},
		})
		if err != nil {
			outputs = append(outputs, openaiv1.RunToolOutput{
				ToolCallID: call.ID,
				Output:     fmt.Sprintf(`{"error":"%s"}`, escapeJSON(err.Error())),
			})
			continue
		}

		payload, err := toolResultPayload(result)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, openaiv1.RunToolOutput{
			ToolCallID: call.ID,
			Output:     payload,
		})
	}

	return outputs, nil
}

func toolResultPayload(result *mcptypes.CallToolResult) (string, error) {
	if result == nil {
		return "", errors.New("empty tool result")
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func threadMessageText(msg openaiv1.ThreadMessageResponse) string {
	if len(msg.Content) == 0 {
		return ""
	}
	parts := make([]string, 0, len(msg.Content))
	for _, item := range msg.Content {
		if item.Type != "text" || item.Text == nil {
			continue
		}
		text := strings.TrimSpace(item.Text.Value)
		if text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n")
}

func toOpenAITools(tools []mcptypes.Tool) ([]openaiv1.Tool, error) {
	openaiTools := make([]openaiv1.Tool, 0, len(tools))
	for _, tool := range tools {
		schema, err := toolSchema(tool)
		if err != nil {
			return nil, err
		}
		openaiTools = append(openaiTools, openaiv1.Tool{
			Type: "function",
			Function: &openaiv1.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  schema,
			},
		})
	}
	return openaiTools, nil
}

func toolSchema(tool mcptypes.Tool) (map[string]interface{}, error) {
	if tool.RawInputSchema != nil {
		schema := map[string]interface{}{}
		if err := json.Unmarshal(tool.RawInputSchema, &schema); err != nil {
			return nil, err
		}
		ensureObjectSchemaProperties(schema)
		return schema, nil
	}

	raw, err := json.Marshal(tool.InputSchema)
	if err != nil {
		return nil, err
	}

	schema := map[string]interface{}{}
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil, err
	}

	ensureObjectSchemaProperties(schema)
	return schema, nil
}

func ensureObjectSchemaProperties(schema map[string]interface{}) {
	if schema == nil {
		return
	}
	if schema["type"] != "object" {
		return
	}
	if _, ok := schema["properties"]; ok {
		return
	}
	schema["properties"] = map[string]interface{}{}
}

func firstAssistantMessage(messages openaiv1.ThreadMessagesListResponse) (string, error) {
	for _, msg := range messages.Data {
		if msg.Role != "assistant" {
			continue
		}
		for _, item := range msg.Content {
			if item.Type == "text" && item.Text != nil {
				return item.Text.Value, nil
			}
		}
	}
	return "", fmt.Errorf("assistant message not found")
}

func escapeJSON(input string) string {
	encoded, _ := json.Marshal(input)
	escaped := string(encoded)
	if len(escaped) >= 2 {
		return escaped[1 : len(escaped)-1]
	}
	return input
}

// RunSubAssistant runs a sub-assistant with MCP tools and returns the response.
// This is used by the orchestrator to delegate to builder/validator/helper assistants.
func RunSubAssistant(ctx context.Context, mcpAPIKey string, assistantID string, message string) (string, error) {
	if assistantID == "" {
		return "", fmt.Errorf("assistant ID is required")
	}
	if message == "" {
		return "", fmt.Errorf("message is required")
	}

	// Create OpenAI client
	openaiClient, err := newOpenAIClient()
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI client: %v", err)
	}

	// Resolve MCP URL
	mcpURL, err := resolveMCPURL()
	if err != nil {
		return "", fmt.Errorf("failed to resolve MCP URL: %v", err)
	}

	// Create MCP client with the provided API key
	mcpClient, err := newMCPClient(ctx, mcpURL, mcpAPIKey)
	if err != nil {
		return "", fmt.Errorf("failed to create MCP client: %v", err)
	}

	// List available MCP tools
	tools, err := listMCPTools(ctx, mcpClient)
	if err != nil {
		return "", fmt.Errorf("failed to list MCP tools: %v", err)
	}

	// Run the assistant (creates new thread, no events callback)
	response, err := runAssistant(
		ctx,
		openaiClient,
		mcpClient,
		assistantID,
		"", // empty threadID = create new thread
		message,
		"", // no additional instructions
		tools,
		func(eventType string, payload any) {}, // no-op event handler
	)
	if err != nil {
		return "", fmt.Errorf("assistant run failed: %v", err)
	}

	return response, nil
}
