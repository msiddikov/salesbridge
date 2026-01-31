package mcp

import (
	"context"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer wraps the mcp-go server with custom auth
type MCPServer struct {
	server     *server.MCPServer
	httpServer *server.StreamableHTTPServer
	auth       *AuthHandler
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(addr string) *MCPServer {
	s := server.NewMCPServer(
		"SalesBridge MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	mcpServer := &MCPServer{
		server: s,
		auth:   NewAuthHandler(),
	}

	// Register tools
	mcpServer.registerTools()

	// Create streamable HTTP server with auth context extraction
	mcpServer.httpServer = server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(true),
		server.WithHeartbeatInterval(30*time.Second),
		server.WithHTTPContextFunc(extractAuthContext),
	)

	return mcpServer
}

// extractAuthContext extracts the API key from HTTP headers and adds it to the context
func extractAuthContext(ctx context.Context, r *http.Request) context.Context {
	authHeader := r.Header.Get("Authorization")
	apiKey := ExtractAPIKeyFromHeader(authHeader)
	if apiKey != "" {
		return ContextWithAPIKey(ctx, apiKey)
	}
	return ctx
}

func (m *MCPServer) registerTools() {
	// get_location tool
	getLocationTool := mcp.NewTool("get_location",
		mcp.WithDescription("Get location information by ID. Returns full details about the location including name, integrations, and configuration. Requires API key with access to the specified location."),
		mcp.WithString("location_id",
			mcp.Required(),
			mcp.Description("The unique identifier of the location to retrieve"),
		),
	)
	m.server.AddTool(getLocationTool, m.handleGetLocation)

	// get_first_or_create_zenoti_guest tool
	getFirstOrCreateGuestTool := mcp.NewTool("get_first_or_create_zenoti_guest",
		mcp.WithDescription("Search for a Zenoti guest by email and phone. If found, returns the guest. If not found, creates a new guest with the provided information."),
		mcp.WithString("location_id",
			mcp.Required(),
			mcp.Description("The location ID to use for Zenoti integration"),
		),
		mcp.WithString("first_name",
			mcp.Description("Guest's first name (required for creating new guest)"),
		),
		mcp.WithString("last_name",
			mcp.Description("Guest's last name (required for creating new guest)"),
		),
		mcp.WithString("email",
			mcp.Description("Guest's email address (used for search and creation)"),
		),
		mcp.WithString("phone",
			mcp.Description("Guest's phone number without special characters and country code (used for search and creation)"),
		),
		mcp.WithString("date_of_birth",
			mcp.Description("Guest's date of birth in YYYY-MM-DD format"),
		),
	)
	m.server.AddTool(getFirstOrCreateGuestTool, m.handleGetFirstOrCreateZenotiGuest)

	// list_zenoti_services tool
	listServicesTool := mcp.NewTool("list_zenoti_services",
		mcp.WithDescription("List all available services for a Zenoti location. Returns service ID, name, description, and duration."),
		mcp.WithString("location_id",
			mcp.Required(),
			mcp.Description("The location ID to list services for"),
		),
	)
	m.server.AddTool(listServicesTool, m.handleListZenotiServices)
}

// Run starts the MCP server with streamable HTTP transport
func (m *MCPServer) Run(addr string) error {
	return m.httpServer.Start(addr)
}

// Shutdown gracefully shuts down the server
func (m *MCPServer) Shutdown(ctx context.Context) error {
	return m.httpServer.Shutdown(ctx)
}
