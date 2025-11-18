package webServer

import (
	"client-runaway-zenoti/internal/cerbo"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/services/auth"
	"client-runaway-zenoti/internal/services/automator"
	"client-runaway-zenoti/internal/services/svc_attribution"
	"client-runaway-zenoti/internal/services/svc_config"
	"client-runaway-zenoti/internal/services/svc_ghl"
	svc_jpmreport "client-runaway-zenoti/internal/services/svc_jpmReport"
	"client-runaway-zenoti/internal/services/svc_zenoti"

	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	// GHL Trigger
	router.POST("/hl/trigger", runway.TriggerSubscriptionsHandler)

	// Cerbo webhooks
	router.POST("/cerbo/webhook/:secret", cerbo.WebhookHandler)

	// JPM Report
	router.GET("/jpm/report", svc_jpmreport.GetReport)

	// Auth routes
	router.POST("/register", auth.Register)
	router.GET("/login", auth.Auth, auth.Login)

	// Automations routes
	auto := router.Group("/auto")
	auto.GET("/catalog", auth.Auth, automator.GetCatalog)

	auto.GET("/:locationId", auth.Auth, automator.GetAutomations)
	auto.POST("/:locationId", auth.Auth, automator.CreateAutomation)
	auto.PATCH("/:automationId", auth.Auth, automator.UpdateAutomation)
	auto.DELETE("/:automationId", auth.Auth, automator.DeleteAutomation)
	auto.POST("/duplicate/:automationId", auth.Auth, automator.DuplicateAutomation)

	auto.GET("/runs/:automationId", auth.Auth, automator.GetAutomationRuns)
	auto.GET("/run-details/:runId", auth.Auth, automator.GetAutomationRunDetails)
	auto.POST("/run/:runId/restart", auth.Auth, automator.StartFromAutomationRun)

	// Batch runs
	auto.POST("/batch-run", auth.Auth, automator.StartBatchRun)
	auto.GET("/batch-runs/:locationId", auth.Auth, automator.GetBatchRuns)
	auto.GET("/batch-run/:batchRunId", auth.Auth, automator.GetBatchRunDetails)
	auto.PATCH("/batch-run/:batchRunId/cancel", auth.Auth, automator.CancelBatchRun)

	// GHL Routes
	ghl := router.Group("/ghl")
	ghl.GET("/:locationId/pipelines", auth.Auth, svc_ghl.GetPipelines)

	// Settings routes
	settings := router.Group("/settings")
	settings.GET("/zenoti/apis", auth.Auth, svc_config.GetZenotiApis)
	settings.POST("/zenoti/apis", auth.Auth, svc_config.CreateZenotiApi)
	settings.PATCH("/zenoti/apis/:zenotiApiId", auth.Auth, svc_config.UpdateZenotiApi)
	settings.DELETE("/zenoti/apis/:zenotiApiId", auth.Auth, svc_config.DeleteZenotiApi)

	settings.PATCH("/locations/:locationId", auth.Auth, svc_config.UpdateLocation)

	settings.POST("/flows", auth.Auth, svc_attribution.CreateAttributionFlow)
	settings.GET("/flows", auth.Auth, svc_attribution.GetAttributionFlows)
	settings.DELETE("/flows/:flowId", auth.Auth, svc_attribution.DeleteAttributionFlow)

	// Integrations routes
	integrations := router.Group("/integrations")
	integrations.GET("/zenoti/centers/:zenotiApiId", auth.Auth, svc_zenoti.GetZenotiCenters)

}
