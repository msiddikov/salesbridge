package webServer

import (
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/services/automator"
	"context"
	"encoding/json"
	"fmt"
	"io"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type (
	commonFields struct {
		Type       string
		LocationId string
	}
)

func setRunwayRoutes(router *gin.Engine) {
	router.POST("/gohighlevel/webhook", webhook)
	router.POST("/gohighlevel/webhook/new_appt", appointmentWebhook)
}

func webhook(c *gin.Context) {

	bodyBytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		panic(err)
	}

	common := commonFields{}
	err = json.Unmarshal(bodyBytes, &common)
	lvn.GinErr(c, 500, err, "Error while parsing webhook data")

	switch common.Type {
	case "OpportunityStageUpdate":
		err = runway.HandleOpportunityStageUpdate(bodyBytes)
		lvn.GinErr(c, 500, err, "Error while handling webhook")
	case "AppointmentCreate":
		err = runway.HandleAppointmentCreate(bodyBytes)
		lvn.GinErr(c, 500, err, "Error while handling webhook")
	case "AppointmentUpdate":
		err = automator.GhlTriggerAppointmentUpdated(context.Background(), bodyBytes)
		lvn.GinErr(c, 500, err, "Error while handling webhook")
	}

	c.Data(lvn.Res(200, nil, "Success"))
}

func appointmentWebhook(c *gin.Context) {

	bodyBytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyBytes))
	err = runway.HandleAppointmentCreate(bodyBytes)
	lvn.GinErr(c, 500, err, "Error while handling appointment webhook")

	c.Data(lvn.Res(200, nil, "Success"))

}
