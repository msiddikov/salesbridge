package webServer

import (
	"client-runaway-zenoti/internal/cerbo"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/services/auth"
	"client-runaway-zenoti/internal/services/automator"
	"client-runaway-zenoti/internal/services/svc_attribution"
	"client-runaway-zenoti/internal/services/svc_cerbo"
	"client-runaway-zenoti/internal/services/svc_config"
	"client-runaway-zenoti/internal/services/svc_ghl"
	"client-runaway-zenoti/internal/services/svc_googleads"
	svc_jpmreport "client-runaway-zenoti/internal/services/svc_jpmReport"
	"client-runaway-zenoti/internal/services/svc_mcp"
	"client-runaway-zenoti/internal/services/svc_openai"
	"client-runaway-zenoti/internal/services/svc_zenoti"

	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	// GHL Trigger
	router.POST("/hl/trigger", runway.TriggerSubscriptionsHandler)
	router.POST("/hl/webhookv2", svc_ghl.WebhookAuthMiddle, ghlWebhookHandler)

	// Cerbo webhooks
	router.POST("/cerbo/webhook/:secret", cerbo.WebhookHandler)

	// JPM Report
	router.GET("/jpm/report", svc_jpmreport.GetReport)

	// Auth routes
	router.POST("/register", auth.Register)
	router.GET("/login", auth.Auth, auth.Login)
	router.GET("/auth/hl", svc_ghl.AddLocationAuthHandler)

	// Automations routes
	auto := router.Group("/auto")
	auto.GET("/catalog", auth.Auth, automator.GetCatalog)
	auto.GET("/catalog/:locationId/list/:listName", auth.Auth, automator.GetLists)

	auto.GET("/:locationId", auth.Auth, automator.GetAutomations)
	auto.POST("/:locationId", auth.Auth, automator.CreateAutomation)
	auto.PATCH("/:automationId", auth.Auth, automator.UpdateAutomation)
	auto.DELETE("/:automationId", auth.Auth, automator.DeleteAutomation)
	auto.POST("/duplicate/:automationId", auth.Auth, automator.DuplicateAutomation)

	auto.GET("/runs", auth.Auth, automator.GetAutomationRuns)
	auto.GET("/runs/export", auth.Auth, automator.ExportAutomationRuns)
	auto.GET("/run-details/:runId", auth.Auth, automator.GetAutomationRunDetails)
	auto.POST("/run/:runId/restart", auth.Auth, automator.StartFromAutomationRun)
	auto.POST("/trigger/:automationId", auth.Auth, automator.StartTriggerForAutomation)

	// Batch runs
	auto.POST("/batch-run", auth.Auth, automator.StartBatchRun)
	auto.GET("/batch-runs/:locationId", auth.Auth, automator.GetBatchRuns)
	auto.GET("/batch-run/:batchRunId", auth.Auth, automator.GetBatchRunDetails)
	auto.PATCH("/batch-run/:batchRunId/cancel", auth.Auth, automator.CancelBatchRun)

	// GHL Routes
	ghl := router.Group("/ghl")
	ghl.GET("/:locationId/pipelines", auth.Auth, svc_ghl.GetPipelines)

	// Zenoti API Settings
	settings := router.Group("/settings")
	settings.GET("/zenoti/apis", auth.Auth, svc_config.GetZenotiApis)
	settings.POST("/zenoti/apis", auth.Auth, svc_config.CreateZenotiApi)
	settings.PATCH("/zenoti/apis/:zenotiApiId", auth.Auth, svc_config.UpdateZenotiApi)
	settings.DELETE("/zenoti/apis/:zenotiApiId", auth.Auth, svc_config.DeleteZenotiApi)

	// Cerbo API Settings
	settings.GET("/cerbo/apis", auth.Auth, svc_config.GetCerboApis)
	settings.POST("/cerbo/apis", auth.Auth, svc_config.CreateCerboApi)
	settings.PATCH("/cerbo/apis/:cerboApiId", auth.Auth, svc_config.UpdateCerboApi)
	settings.DELETE("/cerbo/apis/:cerboApiId", auth.Auth, svc_config.DeleteCerboApi)

	// Locations Settings
	settings.GET("/locations/list", auth.Auth, svc_config.ListLocations)
	settings.PATCH("/locations/:locationId", auth.Auth, svc_config.UpdateLocation)
	settings.DELETE("/locations/:locationId", auth.Auth, svc_config.DeleteLocation)
	settings.GET("/locations/:locationId/delete-preview", auth.Auth, svc_config.DeleteLocationDryRun)
	settings.GET("/oauth/link", auth.Auth, svc_ghl.GetGhlOauthLink)

	// Attribution Flows settings
	settings.POST("/flows", auth.Auth, svc_attribution.CreateAttributionFlow)
	settings.GET("/flows", auth.Auth, svc_attribution.GetAttributionFlows)
	settings.DELETE("/flows/:flowId", auth.Auth, svc_attribution.DeleteAttributionFlow)

	// Integrations routes
	integrations := router.Group("/integrations")
	integrations.GET("/zenoti/centers/:zenotiApiId", auth.Auth, svc_zenoti.GetZenotiCenters)
	integrations.GET("/cerbo/encounter-types/:locationId", auth.Auth, svc_cerbo.GetEncounterTypesForLocation)

	// Google Ads OAuth (new flow; legacy jpmReport remains unchanged)
	ga := router.Group("/google-ads")
	ga.GET("/auth-url", auth.Auth, svc_googleads.GetAuthURL)
	ga.POST("/callback", svc_googleads.OAuthCallback)
	ga.GET("/callback", svc_googleads.OAuthCallback)
	ga.GET("/accounts", auth.Auth, svc_googleads.ListConnections)
	ga.DELETE("/accounts/:accountId", auth.Auth, svc_googleads.DeleteConnection)
	ga.GET("/accounts/:accountId/hierarchy", auth.Auth, svc_googleads.ListAccountHierarchy)
	ga.POST("/locations/:locationId/settings", auth.Auth, svc_googleads.SaveLocationSetting)
	ga.GET("/locations/:locationId/settings", auth.Auth, svc_googleads.GetLocationSetting)
	ga.GET("/locations/:locationId/conversion-actions", auth.Auth, svc_googleads.GetLocationConversionActions)

	// Legacy aliases to avoid breaking existing references
	legacyGA := router.Group("/googleads")
	legacyGA.GET("/oauth/url", auth.Auth, svc_googleads.GetAuthURL)
	legacyGA.GET("/connections", auth.Auth, svc_googleads.ListConnections)
	legacyGA.DELETE("/connections/:accountId", auth.Auth, svc_googleads.DeleteConnection)
	router.GET("/googleads/oauth/callback", svc_googleads.OAuthCallback)

	// AI Assistants
	ai := router.Group("/ai")
	ai.GET("/models", auth.Auth, svc_openai.ListModels)
	ai.GET("/pricing", auth.Auth, svc_openai.ListPricing)
	ai.PUT("/pricing", auth.Auth, svc_openai.UpsertPricingBatch)
	ai.PUT("/pricing/:model", auth.Auth, svc_openai.UpsertPricing)

	ai.GET("/usage", auth.Auth, svc_openai.ListUsage)
	ai.GET("/assistants", auth.Auth, svc_openai.ListAssistants)
	ai.GET("/assistants/:assistantId", auth.Auth, svc_openai.GetAssistant)
	ai.POST("/assistants", auth.Auth, svc_openai.CreateAssistant)
	ai.PATCH("/assistants/:assistantId", auth.Auth, svc_openai.UpdateAssistant)
	ai.DELETE("/assistants/:assistantId", auth.Auth, svc_openai.DeleteAssistant)

	// MCP API Key management
	mcpKeys := router.Group("/mcp/keys")
	mcpKeys.POST("", auth.Auth, svc_mcp.CreateAPIKey)
	mcpKeys.GET("", auth.Auth, svc_mcp.ListAPIKeys)
	mcpKeys.PATCH("/:keyId", auth.Auth, svc_mcp.UpdateAPIKey)
	mcpKeys.PATCH("/:keyId/revoke", auth.Auth, svc_mcp.RevokeAPIKey)
	mcpKeys.POST("/:keyId/regenerate", auth.Auth, svc_mcp.RegenerateAPIKey)
	mcpKeys.DELETE("/:keyId", auth.Auth, svc_mcp.DeleteAPIKey)

}
