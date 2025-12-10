package googleads

import (
	"fmt"
	"time"

	"github.com/shenzhencenter/google-ads-pb/enums"
	"github.com/shenzhencenter/google-ads-pb/services"
)

// ConversionRequest carries conversion details plus optional UTM custom variable mapping.
type ConversionRequest struct {
	ConversionActionID string
	Gclid              string
	Gbraid             string
	Wbraid             string

	ConversionValue float64
	CurrencyCode    string
	EventTime       time.Time
	ValidateOnly    bool

	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMTerm     string
	UTMContent  string

	UTMCustomVariableResourceNames map[string]string
}

// SendConversion uploads a click conversion with optional UTM custom variables.
// ctx must already include auth/developer headers (Service.WithHeaders).
func (c *Client) SendConversion(req ConversionRequest) error {

	if req.ConversionActionID == "" {
		return fmt.Errorf("googleads: conversion action id required")
	}
	if req.Gclid == "" && req.Gbraid == "" && req.Wbraid == "" {
		return fmt.Errorf("googleads: one of gclid, gbraid, or wbraid is required")
	}

	eventTime := req.EventTime
	if eventTime.IsZero() {
		eventTime = time.Now()
	}
	currency := req.CurrencyCode
	if currency == "" {
		currency = "USD"
	}

	var customVars []*services.CustomVariable

	if len(req.UTMCustomVariableResourceNames) > 0 {
		addVar := func(key, value string) {
			resName := req.UTMCustomVariableResourceNames[key]
			if resName == "" || value == "" {
				return
			}
			customVars = append(customVars, &services.CustomVariable{
				ConversionCustomVariable: resName,
				Value:                    value,
			})
		}
		addVar("source", req.UTMSource)
		addVar("medium", req.UTMMedium)
		addVar("campaign", req.UTMCampaign)
		addVar("term", req.UTMTerm)
		addVar("content", req.UTMContent)
	}

	convDateString := eventTime.Format("2006-01-02 15:04:05-07:00")
	actionStr := fmt.Sprintf("customers/%s/conversionActions/%s", c.CustomerInfo.CustomerID, req.ConversionActionID)
	click := &services.ClickConversion{
		ConversionAction:   &actionStr,
		ConversionDateTime: &convDateString,
		CurrencyCode:       &currency,
		ConversionValue:    &req.ConversionValue,
		CustomVariables:    customVars,
		Gclid:              &req.Gclid,
		Gbraid:             req.Gbraid,
		Wbraid:             req.Wbraid,
	}

	svc := services.NewConversionUploadServiceClient(c.grpcConn)
	_, err := svc.UploadClickConversions(c.ctx, &services.UploadClickConversionsRequest{
		CustomerId:     c.CustomerInfo.CustomerID,
		ValidateOnly:   req.ValidateOnly,
		Conversions:    []*services.ClickConversion{click},
		PartialFailure: true,
	})
	return err
}

// ConversionActionSummary is a lightweight view of a conversion action.
type ConversionActionSummary struct {
	ResourceName string
	Name         string
	Type         string
	Status       string
}

// ListConversionActions fetches conversion actions for the specified customer account.
// If managerID is provided, it is used as login-customer-id for the request.
func (c *Client) ListConversionActions() ([]ConversionActionSummary, error) {

	query := `
SELECT
  conversion_action.resource_name,
  conversion_action.name,
  conversion_action.type,
  conversion_action.status
FROM conversion_action
`

	req := &services.SearchGoogleAdsRequest{
		CustomerId: c.CustomerInfo.CustomerID,
		Query:      query,
	}

	svc := services.NewGoogleAdsServiceClient(c.grpcConn)
	resp, err := svc.Search(c.ctx, req)
	if err != nil {
		return nil, err
	}

	out := make([]ConversionActionSummary, 0, len(resp.Results))
	for _, row := range resp.Results {
		ca := row.GetConversionAction()
		if ca == nil {
			continue
		}
		out = append(out, ConversionActionSummary{
			ResourceName: ca.GetResourceName(),
			Name:         ca.GetName(),
			Type:         ca.GetType().String(),
			Status:       ca.GetStatus().String(),
		})
	}
	return out, nil
}

// ConversionAdjustmentRequest carries data to restate or retract a conversion.
type ConversionAdjustmentRequest struct {
	CustomerID       string
	ClientCustomerID string
	ManagerID        string

	ConversionActionID  string
	AdjustmentType      enums.ConversionAdjustmentTypeEnum_ConversionAdjustmentType
	Gclid               string
	OriginalEventTime   time.Time
	OrderID             string
	RestatementValue    float64
	RestatementCurrency string
	AdjustmentTime      time.Time
	ValidateOnly        bool
}

// SendConversionAdjustment uploads a conversion adjustment (restatement/retraction/enhancement).
// ctx must already include auth/developer headers (Service.WithHeaders).
func (c *Client) SendConversionAdjustment(req ConversionAdjustmentRequest) error {

	if req.ConversionActionID == "" {
		return fmt.Errorf("googleads: conversion action id required")
	}
	if req.Gclid == "" && req.OrderID == "" {
		return fmt.Errorf("googleads: gclid or order_id required")
	}

	adjTime := req.AdjustmentTime
	if adjTime.IsZero() {
		adjTime = time.Now()
	}
	origTime := req.OriginalEventTime
	if origTime.IsZero() {
		return fmt.Errorf("googleads: original conversion time required for adjustment")
	}

	adjType := req.AdjustmentType
	if adjType == enums.ConversionAdjustmentTypeEnum_UNSPECIFIED {
		adjType = enums.ConversionAdjustmentTypeEnum_RESTATEMENT
	}

	adjTimeStr := adjTime.Format("2006-01-02 15:04:05-07:00")
	origTimeStr := origTime.Format("2006-01-02 15:04:05-07:00")
	actionRes := fmt.Sprintf("customers/%s/conversionActions/%s", c.CustomerInfo.CustomerID, req.ConversionActionID)

	convAdj := &services.ConversionAdjustment{
		ConversionAction:   &actionRes,
		AdjustmentType:     adjType,
		AdjustmentDateTime: &adjTimeStr,
	}

	if req.Gclid != "" {
		convAdj.GclidDateTimePair = &services.GclidDateTimePair{
			Gclid:              &req.Gclid,
			ConversionDateTime: &origTimeStr,
		}
	}
	if req.OrderID != "" {
		convAdj.OrderId = &req.OrderID
	}
	if adjType == enums.ConversionAdjustmentTypeEnum_RESTATEMENT || adjType == enums.ConversionAdjustmentTypeEnum_ENHANCEMENT {
		currency := req.RestatementCurrency
		if currency == "" {
			currency = "USD"
		}
		convAdj.RestatementValue = &services.RestatementValue{
			AdjustedValue: &req.RestatementValue,
			CurrencyCode:  &currency,
		}
	}

	svc := services.NewConversionAdjustmentUploadServiceClient(c.grpcConn)
	_, err := svc.UploadConversionAdjustments(c.ctx, &services.UploadConversionAdjustmentsRequest{
		CustomerId:            c.CustomerInfo.CustomerID,
		ValidateOnly:          req.ValidateOnly,
		ConversionAdjustments: []*services.ConversionAdjustment{convAdj},
		PartialFailure:        true,
	})
	return err
}
