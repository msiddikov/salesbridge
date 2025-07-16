package webServer

import (
	"client-runaway-zenoti/internal/cerbo"
	"client-runaway-zenoti/internal/runway"

	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	// GHL Trigger
	router.POST("/hl/trigger", runway.TriggerSubscriptionsHandler)

	// Cerbo webhooks
	router.POST("/cerbo/webhook/:secret", cerbo.WebhookHandler)
}
