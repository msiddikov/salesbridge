package webServer

import (
	"bytes"
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/services/automator"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"context"
	"encoding/json"
	"io"
	"net/http"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func ghlWebhookHandler(c *gin.Context) {

	genericPayload := runwayv2.WebhookGenericPayload{}
	xWhSignature := c.GetHeader("x-wh-signature")

	body, _ := io.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &genericPayload)
	if err != nil {
		c.Data(lvn.Res(400, "", "invalid payload"))
		return
	}

	switch genericPayload.Type {
	case "AppointmentUpdate":
		automator.GhlTriggerAppointmentUpdated(context.Background(), body)
		transferToDevServer(body, xWhSignature)
	case "ContactCreate":
	}

	// respond with 200 OK to GHL
	c.Data(lvn.Res(200, "Success", "ok"))
}

func transferToDevServer(body []byte, xWhSignature string) {
	if config.Confs.Settings.SrvDomain == "https://salesbridge-api.lavina.tech" {
		// send the body with x-wh-signature to the dev server
		req, _ := http.NewRequest("POST", "https://mason.lavina.uz/hl/webhookv2", bytes.NewBuffer(body))
		req.Header.Set("x-wh-signature", xWhSignature)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		_, _ = client.Do(req)
	}
}
