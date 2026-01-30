package svc_googleads

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/googleads"
	"fmt"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// GetLocationConversionActions lists conversion actions for the account configured on a location.
func GetLocationConversionActions(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locID := c.Param("locationId")
	if locID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("location id required"), "invalid location id")
		return
	}

	cli, err := CliForLocation(locID, user.ProfileID)
	lvn.GinErr(c, 400, err, "unable to create google ads client")
	defer cli.Close()

	actions, err := cli.ListConversionActions()
	lvn.GinErr(c, 400, err, "unable to list conversion actions")

	c.Data(lvn.Res(200, actions, "success"))
}

// GetLocationConversionActions lists conversion actions for the account configured on a location.
func GetLocationConversionActionsList(l models.Location) (map[string]string, error) {

	cli, err := CliForLocation(l.Id, l.ProfileID)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	actions, err := cli.ListConversionActions()
	if err != nil {
		return nil, err
	}

	actionMap := make(map[string]string)
	for _, action := range actions {
		// Get everything after last '/' in ResourceName
		parts := strings.Split(action.ResourceName, "/")
		id := parts[len(parts)-1]
		actionMap[action.Name] = id
	}

	return actionMap, nil
}

// UploadConversionData uploads offline conversion data for the location's configured account.
func UploadConversionData(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locID := c.Param("locationId")
	if locID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("location id required"), "invalid location id")
		return
	}

	var payload struct {
		ConversionActionID string `json:"conversionActionId" binding:"required"`
		Gclid              string `json:"gclid"`
		Gbraid             string `json:"gbraid"`
		Wbraid             string `json:"wbraid"`

		ConversionValue float64   `json:"conversionValue"`
		CurrencyCode    string    `json:"currencyCode"`
		EventTime       time.Time `json:"eventTime"`
		ValidateOnly    bool      `json:"validateOnly"`

		UTMCustomVariableResourceNames map[string]string `json:"utmCustomVariableResourceNames"`
	}
	err := c.ShouldBindJSON(&payload)
	lvn.GinErr(c, 400, err, "invalid payload")
	if err != nil {
		return
	}
	if payload.Gclid == "" && payload.Gbraid == "" && payload.Wbraid == "" {
		lvn.GinErr(c, 400, fmt.Errorf("gclid, gbraid, or wbraid required"), "invalid payload")
		return
	}

	cli, err := CliForLocation(locID, user.ProfileID)
	lvn.GinErr(c, 400, err, "unable to create google ads client")
	if err != nil {
		return
	}
	defer cli.Close()

	err = cli.SendConversion(googleads.ConversionRequest{
		ConversionActionID: payload.ConversionActionID,
		Gclid:              payload.Gclid,
		Gbraid:             payload.Gbraid,
		Wbraid:             payload.Wbraid,
		ConversionValue:    payload.ConversionValue,
		CurrencyCode:       payload.CurrencyCode,
		EventTime:          payload.EventTime,
		ValidateOnly:       payload.ValidateOnly,

		UTMCustomVariableResourceNames: payload.UTMCustomVariableResourceNames,
	})
	lvn.GinErr(c, 400, err, "unable to upload conversion")
	if err != nil {
		return
	}
}
