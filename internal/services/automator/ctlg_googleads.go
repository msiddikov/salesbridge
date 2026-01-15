package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_googleads"
	"client-runaway-zenoti/packages/googleads"
	"context"
)

var (
	// Google Ads Category
	gaCategory = Category{
		Id:    "ga",
		Name:  "Google Ads",
		Icon:  "ri:share-forward-2",
		Color: "#4C6FFF",
		Nodes: []Node{
			gaActionUploadConversionData,
		},
	}
)

var gaActionUploadConversionData = Node{
	Id:          "ga.conversionUpload",
	Title:       "Upload Conversion Data",
	Description: "Uploads offline conversion data to Google Ads.",
	ExecFunc:    gaUploadConversionData,
	Type:        NodeTypeAction,
	Icon:        "ri:search-2-line",
	Kind:        "Conversions",
	Color:       ColorAction,
	Ports: []NodePort{
		successPort([]NodeField{}),
		errorPort,
	},
	Fields: []NodeField{
		{Key: "googleAdsActionId", Type: "string"},
		{Key: "gclid", Type: "string"},
		{Key: "eventTime", Type: "string"},
		{Key: "value", Type: "string"},
		{Key: "orderId", Type: "string"},
	},
}

func gaUploadConversionData(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	client, err := svc_googleads.CliForLocation(l.Id, l.ProfileID)
	if err != nil {
		return errorPayload(err, "unable to create google ads client")
	}
	defer client.Close()

	eventTime, err := parseTime(fields["eventTime"].(string))
	if err != nil {
		return errorPayload(err, "invalid event time format")
	}
	value, ok := toFloat(fields["value"])
	if !ok {
		return errorPayload(err, "invalid value format")
	}

	actionId, ok := fields["googleAdsActionId"].(string)
	if !ok || actionId == "" {
		return errorPayload(err, "googleAdsActionId is required")
	}

	req := googleads.ConversionRequest{
		ConversionActionID: actionId,
		Gclid:              fields["gclid"].(string),
		EventTime:          eventTime,
		ConversionValue:    value,
		OrderId:            fields["orderId"].(string),
	}
	err = client.SendConversion(req)
	if err != nil {
		return errorPayload(err, "unable to upload conversion data")
	}

	return successPayload(map[string]interface{}{})
}
