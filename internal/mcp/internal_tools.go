package mcp

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/automator"
	"client-runaway-zenoti/internal/services/svc_config"
	"client-runaway-zenoti/internal/services/svc_internal_assistant"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ============================================================================
// Internal Tool Registry
// ============================================================================

// internalToolNames contains the names of tools only visible to internal API keys
var internalToolNames = map[string]bool{
	"list_locations":           true,
	"get_location_details":     true,
	"list_automations":         true,
	"get_automation":           true,
	"list_automation_runs":     true,
	"get_automation_run":       true,
	"list_batch_runs":          true,
	"list_integrations":        true,
	"get_catalog":              true,
	"get_catalog_list":         true,
	"run_builder_assistant":    true,
	"run_validator_assistant":  true,
	"run_helper_assistant":     true,
}

// IsInternalTool checks if a tool name is internal-only
func IsInternalTool(toolName string) bool {
	return internalToolNames[toolName]
}

// ============================================================================
// Tool Filter - Complete Separation
// ============================================================================

// createToolFilter returns a ToolFilterFunc that provides complete separation:
// - Internal keys see ONLY internal tools
// - Customer keys see ONLY customer tools
func (m *MCPServer) createToolFilter() server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		apiKey, err := m.auth.ValidateAndGetKey(ctx)
		if err != nil {
			// No valid key - return only customer tools (they'll still fail auth)
			return filterOutInternalTools(tools)
		}

		if apiKey.IsInternal {
			// Internal key: return ONLY internal tools
			return filterToInternalTools(tools)
		}

		// Customer key: return ONLY customer tools
		return filterOutInternalTools(tools)
	}
}

// filterToInternalTools returns only internal tools
func filterToInternalTools(tools []mcp.Tool) []mcp.Tool {
	filtered := make([]mcp.Tool, 0)
	for _, tool := range tools {
		if IsInternalTool(tool.Name) {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// filterOutInternalTools returns only customer tools (non-internal)
func filterOutInternalTools(tools []mcp.Tool) []mcp.Tool {
	filtered := make([]mcp.Tool, 0)
	for _, tool := range tools {
		if !IsInternalTool(tool.Name) {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// ============================================================================
// Tool Handler Middleware - Execution Protection
// ============================================================================

// createToolMiddleware returns middleware that blocks cross-access:
// - Internal tools can only be called by internal keys
// - Customer tools can only be called by customer keys
func (m *MCPServer) createToolMiddleware() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			toolName := request.Params.Name
			isInternal := IsInternalTool(toolName)

			apiKey, err := m.auth.ValidateAndGetKey(ctx)
			if err != nil {
				return mcp.NewToolResultError("authentication required"), nil
			}

			// Block cross-access
			if isInternal && !apiKey.IsInternal {
				return mcp.NewToolResultError("access denied: this tool is not available"), nil
			}
			if !isInternal && apiKey.IsInternal {
				return mcp.NewToolResultError("access denied: this tool is not available for internal keys"), nil
			}

			return next(ctx, request)
		}
	}
}

// ============================================================================
// Internal Tool Registration
// ============================================================================

func (m *MCPServer) registerInternalTools() {
	// list_locations
	m.server.AddTool(
		mcp.NewTool("list_locations",
			mcp.WithDescription("List all locations for the current profile with their integration status"),
		),
		m.handleListLocations,
	)

	// get_location_details
	m.server.AddTool(
		mcp.NewTool("get_location_details",
			mcp.WithDescription("Get detailed information about a specific location including all integrations"),
			mcp.WithString("location_id", mcp.Required(), mcp.Description("The location ID to get details for")),
		),
		m.handleGetLocationDetails,
	)

	// list_automations
	m.server.AddTool(
		mcp.NewTool("list_automations",
			mcp.WithDescription("List automations for a location or all locations"),
			mcp.WithString("location_id", mcp.Description("Optional: filter by location ID")),
			mcp.WithString("state", mcp.Description("Optional: filter by state (draft, active, archived)")),
		),
		m.handleListAutomations,
	)

	// get_automation
	m.server.AddTool(
		mcp.NewTool("get_automation",
			mcp.WithDescription("Get full automation details including nodes and edges"),
			mcp.WithString("automation_id", mcp.Required(), mcp.Description("The automation ID to retrieve")),
		),
		m.handleGetAutomation,
	)

	// list_automation_runs
	m.server.AddTool(
		mcp.NewTool("list_automation_runs",
			mcp.WithDescription("List automation runs with optional filtering"),
			mcp.WithString("automation_id", mcp.Description("Optional: filter by automation ID")),
			mcp.WithString("status", mcp.Description("Optional: filter by status (pending, running, success, failed, with_errors, canceled)")),
			mcp.WithNumber("limit", mcp.Description("Optional: max results to return (default 20, max 100)")),
		),
		m.handleListAutomationRuns,
	)

	// get_automation_run
	m.server.AddTool(
		mcp.NewTool("get_automation_run",
			mcp.WithDescription("Get detailed information about a specific automation run including node executions"),
			mcp.WithString("run_id", mcp.Required(), mcp.Description("The automation run ID to retrieve")),
		),
		m.handleGetAutomationRun,
	)

	// list_batch_runs
	m.server.AddTool(
		mcp.NewTool("list_batch_runs",
			mcp.WithDescription("List batch runs for a location"),
			mcp.WithString("location_id", mcp.Required(), mcp.Description("The location ID to list batch runs for")),
			mcp.WithString("status", mcp.Description("Optional: filter by status")),
		),
		m.handleListBatchRuns,
	)

	// list_integrations
	m.server.AddTool(
		mcp.NewTool("list_integrations",
			mcp.WithDescription("List all configured integrations (Zenoti, Cerbo, Google Ads) for the profile"),
		),
		m.handleListIntegrations,
	)

	// get_catalog
	m.server.AddTool(
		mcp.NewTool("get_catalog",
			mcp.WithDescription("Get the full automation nodes catalog with all available triggers, actions, and conditions"),
		),
		m.handleGetCatalog,
	)

	// get_catalog_list
	m.server.AddTool(
		mcp.NewTool("get_catalog_list",
			mcp.WithDescription("Get dynamic list values for node fields (e.g., cerboUsers, aiAssistants, ghlPipelines)"),
			mcp.WithString("location_id", mcp.Required(), mcp.Description("The location ID for context")),
			mcp.WithString("list_name", mcp.Required(), mcp.Description("The list name to retrieve (e.g., cerboUsers, aiAssistants, googleAdsActions, ghlPipelines)")),
		),
		m.handleGetCatalogList,
	)

	// run_builder_assistant - Delegate to builder assistant for creating automations
	m.server.AddTool(
		mcp.NewTool("run_builder_assistant",
			mcp.WithDescription("Run the builder assistant to help create or modify automations. Returns the builder's response."),
			mcp.WithString("location_id", mcp.Required(), mcp.Description("The location ID where the automation will be created")),
			mcp.WithString("user_intent", mcp.Required(), mcp.Description("What the user wants to achieve with the automation")),
			mcp.WithString("mentioned_nodes", mcp.Description("Specific nodes or actions the user mentioned wanting to use")),
			mcp.WithString("user_acceptance_criteria", mcp.Description("Criteria for what makes the automation successful")),
		),
		m.handleRunBuilderAssistant,
	)

	// run_validator_assistant - Delegate to validator assistant for reviewing automations
	m.server.AddTool(
		mcp.NewTool("run_validator_assistant",
			mcp.WithDescription("Run the validator assistant to review and validate an automation. Returns validation results and suggestions."),
			mcp.WithString("user_intent", mcp.Required(), mcp.Description("The original intent behind the automation")),
			mcp.WithString("user_acceptance_criteria", mcp.Description("Criteria for what makes the automation successful")),
			mcp.WithObject("automation", mcp.Required(), mcp.Description("The automation object to validate")),
		),
		m.handleRunValidatorAssistant,
	)

	// run_helper_assistant - Delegate to helper assistant for answering questions
	m.server.AddTool(
		mcp.NewTool("run_helper_assistant",
			mcp.WithDescription("Run the helper assistant to answer questions and find information. Returns the helper's response."),
			mcp.WithString("user_intent", mcp.Required(), mcp.Description("What the user is trying to understand or accomplish")),
			mcp.WithString("user_problem", mcp.Description("The specific problem or question the user has")),
			mcp.WithString("context_json", mcp.Description("JSON string of context items array, e.g. [{\"name\":\"key\",\"value\":\"val\"}]")),
		),
		m.handleRunHelperAssistant,
	)
}

// ============================================================================
// Internal Tool Handlers
// ============================================================================

// authenticateInternal validates API key and ensures it's an internal key
func (m *MCPServer) authenticateInternal(ctx context.Context) (*models.MCPApiKey, error) {
	apiKey, err := m.auth.ValidateAndGetKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	if !apiKey.IsInternal {
		return nil, fmt.Errorf("access denied: internal key required")
	}
	return apiKey, nil
}

func (m *MCPServer) handleListLocations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locations, err := svc_config.GetLocationsForProfile(apiKey.ProfileID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch locations: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(locations, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleGetLocationDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locationID, err := request.RequireString("location_id")
	if err != nil {
		return mcp.NewToolResultError("location_id is required"), nil
	}

	location, err := svc_config.GetLocationForProfile(apiKey.ProfileID, locationID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Location not found: %v", err)), nil
	}

	// Build detailed response (excluding sensitive fields)
	response := map[string]any{
		"id":               location.Id,
		"name":             location.Name,
		"profileId":        location.ProfileID,
		"zenotiCenterId":   location.ZenotiCenterId,
		"zenotiCenterName": location.ZenotiCenterName,
		"hasZenoti":        location.ZenotiCenterId != "",
		"hasCerbo":         location.CerboApiObjId != nil && *location.CerboApiObjId > 0,
	}

	if location.CerboApiObjId != nil && *location.CerboApiObjId > 0 {
		response["cerboSubdomain"] = location.CerboApiObj.Subdomain
	}

	jsonBytes, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleListAutomations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locationID, _ := request.RequireString("location_id")
	state, _ := request.RequireString("state")

	automations, err := automator.GetAutomationsForProfile(apiKey.ProfileID, locationID, state)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch automations: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(automations, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleGetAutomation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	automationID, err := request.RequireString("automation_id")
	if err != nil {
		return mcp.NewToolResultError("automation_id is required"), nil
	}

	automation, err := automator.GetAutomationForProfile(apiKey.ProfileID, automationID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Automation not found: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(automation, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleListAutomationRuns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	automationID, _ := request.RequireString("automation_id")
	status, _ := request.RequireString("status")
	limitFloat, _ := request.RequireFloat("limit")
	limit := int(limitFloat)

	runs, err := automator.GetAutomationRunsForProfile(apiKey.ProfileID, automationID, status, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch runs: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(runs, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleGetAutomationRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	runID, err := request.RequireString("run_id")
	if err != nil {
		return mcp.NewToolResultError("run_id is required"), nil
	}

	run, err := automator.GetAutomationRunForProfile(apiKey.ProfileID, runID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Run not found: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(run, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleListBatchRuns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locationID, err := request.RequireString("location_id")
	if err != nil {
		return mcp.NewToolResultError("location_id is required"), nil
	}

	status, _ := request.RequireString("status")

	batchRuns, err := automator.GetBatchRunsForLocation(apiKey.ProfileID, locationID, status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch batch runs: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(batchRuns, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleListIntegrations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	integrations := svc_config.GetIntegrationsForProfile(apiKey.ProfileID)

	jsonBytes, _ := json.MarshalIndent(integrations, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleGetCatalog(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get the catalog from automator package
	catalog := automator.GetCatalogData()

	jsonBytes, _ := json.MarshalIndent(catalog, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (m *MCPServer) handleGetCatalogList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locationID, err := request.RequireString("location_id")
	if err != nil {
		return mcp.NewToolResultError("location_id is required"), nil
	}

	listName, err := request.RequireString("list_name")
	if err != nil {
		return mcp.NewToolResultError("list_name is required"), nil
	}

	// Verify location belongs to profile and get with preloaded relations
	location, err := svc_config.GetLocationForProfile(apiKey.ProfileID, locationID)
	if err != nil {
		return mcp.NewToolResultError("Location not found or access denied"), nil
	}

	// Get the list data from automator
	listData, err := automator.GetCatalogListData(listName, *location)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get list data: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(listData, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// ============================================================================
// Assistant Runner Tools
// ============================================================================

// runSubAssistant runs a sub-assistant with MCP tools and returns the response.
// It uses the internal MCP API key from the authenticated request.
func (m *MCPServer) runSubAssistant(ctx context.Context, mcpAPIKey string, assistantID string, message string) (string, error) {
	return svc_internal_assistant.RunSubAssistant(ctx, mcpAPIKey, assistantID, message)
}

func (m *MCPServer) handleRunBuilderAssistant(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get the builder assistant ID
	assistantID := strings.TrimSpace(config.Confs.Settings.InternalBuilderID)
	if assistantID == "" {
		return mcp.NewToolResultError("builder assistant ID not configured"), nil
	}

	// Extract parameters
	locationID, _ := request.RequireString("location_id")
	userIntent, _ := request.RequireString("user_intent")
	mentionedNodes, _ := request.RequireString("mentioned_nodes")
	acceptanceCriteria, _ := request.RequireString("user_acceptance_criteria")

	// Build the input message as JSON
	input := map[string]string{
		"locationId":             locationID,
		"userIntent":             userIntent,
		"mentionedNodes":         mentionedNodes,
		"userAcceptanceCriteria": acceptanceCriteria,
	}

	inputJSON, _ := json.Marshal(input)
	response, err := m.runSubAssistant(ctx, apiKey.PlainKey, assistantID, string(inputJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Builder assistant error: %v", err)), nil
	}

	return mcp.NewToolResultText(response), nil
}

func (m *MCPServer) handleRunValidatorAssistant(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get the validator assistant ID
	assistantID := strings.TrimSpace(config.Confs.Settings.InternalValidatorID)
	if assistantID == "" {
		return mcp.NewToolResultError("validator assistant ID not configured"), nil
	}

	// Extract parameters
	userIntent, _ := request.RequireString("user_intent")
	acceptanceCriteria, _ := request.RequireString("user_acceptance_criteria")

	// Get automation object from arguments
	var automation any
	if args, ok := request.Params.Arguments.(map[string]any); ok {
		automation = args["automation"]
	}

	// Build the input message as JSON
	input := map[string]any{
		"userIntent":             userIntent,
		"userAcceptanceCriteria": acceptanceCriteria,
		"automation":             automation,
	}

	inputJSON, _ := json.Marshal(input)
	response, err := m.runSubAssistant(ctx, apiKey.PlainKey, assistantID, string(inputJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Validator assistant error: %v", err)), nil
	}

	return mcp.NewToolResultText(response), nil
}

func (m *MCPServer) handleRunHelperAssistant(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := m.authenticateInternal(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get the helper assistant ID
	assistantID := strings.TrimSpace(config.Confs.Settings.InternalHelperID)
	if assistantID == "" {
		return mcp.NewToolResultError("helper assistant ID not configured"), nil
	}

	// Extract parameters
	userIntent, _ := request.RequireString("user_intent")
	userProblem, _ := request.RequireString("user_problem")
	contextJSON, _ := request.RequireString("context_json")

	// Parse context JSON if provided
	var contextItems any
	if contextJSON != "" {
		if err := json.Unmarshal([]byte(contextJSON), &contextItems); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid context_json: %v", err)), nil
		}
	}

	// Build the input message as JSON
	input := map[string]any{
		"userIntent":  userIntent,
		"userProblem": userProblem,
		"context":     contextItems,
	}

	inputJSON, _ := json.Marshal(input)
	response, err := m.runSubAssistant(ctx, apiKey.PlainKey, assistantID, string(inputJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Helper assistant error: %v", err)), nil
	}

	return mcp.NewToolResultText(response), nil
}
