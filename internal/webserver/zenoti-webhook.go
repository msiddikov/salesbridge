package webServer

import (
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/tgbot"
	"client-runaway-zenoti/internal/zenoti"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"encoding/json"
	"fmt"
	"io"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type (
	Event_type string
)

func setZenotiWebhooksRoutes(router *gin.Engine) {
	router.POST("/zenoti/appointment", newAppointment)
	router.POST("/zenoti/sales", newAppointment)

	router.POST("/zenoti/webhook", zenotiWebhook)
}

func newAppointment(c *gin.Context) {
	body := runway.SurveyForm{}
	c.BindJSON(&body)
	err := runway.SurveyPost(body, c.Param("locationId"), c.Param("workflowId"))
	if err != nil {
		panic(err)
	}
}

func zenotiWebhook(c *gin.Context) {

	// get body as a string
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	if len(bodyBytes) == 0 {
		c.Writer.WriteHeader(200)
		return
	}

	// Unmarshal the body
	body := zenotiv1.WebhookData{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		panic(err)
	}

	switch body.Event_type {
	case "Invoice.Closed":

		err = runway.UpdateCollectionFromWebhook(body)
		if err != nil {
			tgbot.Notify("Webhook Data", fmt.Sprintf("%s\nDATA: %s", err.Error(), string(bodyBytes)), true)
			panic(err)
		}
		c.Data(lvn.Res(200, nil, "Success"))
	case "Guest.Created":
		err = zenoti.GuestCreatedWebhookHandler(bodyBytes)
		if err != nil {
			tgbot.Notify("Webhook Data", fmt.Sprintf("%s\nDATA: %s", err.Error(), string(bodyBytes)), true)
		}
		c.Data(lvn.Res(200, nil, "Success"))
	case "Guest.Updated":
		err = zenoti.GuestCreatedWebhookHandler(bodyBytes)
		if err != nil {
			tgbot.Notify("Webhook Data", fmt.Sprintf("%s\nDATA: %s", err.Error(), string(bodyBytes)), true)
		}
		c.Data(lvn.Res(200, nil, "Success"))
	default:
		c.Data(lvn.Res(200, nil, "Success"))
	}
}

func webhookAuth(c *gin.Context) {

	// check for key header
	key := c.GetHeader("api-key")
	if key != "JFM3oiu098&&83Dfg56)567" {
		tgbot.Notify("Webhook Auth", "Unauthorized", true)
		return
		c.JSON(401, "")
		c.Abort()
	}
}
