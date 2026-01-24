package svc_openai

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type pricingPayload struct {
	InputCentsPer1K  float64 `json:"inputCentsPer1k"`
	OutputCentsPer1K float64 `json:"outputCentsPer1k"`
}

type pricingBatchItem struct {
	Model            string  `json:"model"`
	InputCentsPer1K  float64 `json:"inputCentsPer1k"`
	OutputCentsPer1K float64 `json:"outputCentsPer1k"`
}

type pricingResponse struct {
	Model            string  `json:"model"`
	InputCentsPer1K  float64 `json:"inputCentsPer1k"`
	OutputCentsPer1K float64 `json:"outputCentsPer1k"`
	UpdatedAt        int64   `json:"updatedAt"`
}

func ListPricing(c *gin.Context) {
	var pricing []models.OpenAIModelPricing

	err := db.DB.Order("model asc").Find(&pricing).Error
	lvn.GinErr(c, 400, err, "unable to list pricing")
	if err != nil {
		return
	}

	resp := make([]pricingResponse, 0, len(pricing))
	for _, p := range pricing {
		resp = append(resp, pricingResponse{
			Model:            p.Model,
			InputCentsPer1K:  p.InputCentsPer1K,
			OutputCentsPer1K: p.OutputCentsPer1K,
			UpdatedAt:        p.UpdatedAt.Unix(),
		})
	}

	c.Data(lvn.Res(200, resp, "OK"))
}

func UpsertPricing(c *gin.Context) {
	model := strings.TrimSpace(c.Param("model"))
	if model == "" {
		lvn.GinErr(c, 400, nil, "model is required")
		return
	}

	payload := pricingPayload{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")
	if err != nil {
		return
	}

	if payload.InputCentsPer1K < 0 || payload.OutputCentsPer1K < 0 {
		lvn.GinErr(c, 400, nil, "pricing must be non-negative")
		return
	}

	record := models.OpenAIModelPricing{
		Model:            model,
		InputCentsPer1K:  payload.InputCentsPer1K,
		OutputCentsPer1K: payload.OutputCentsPer1K,
	}

	err = db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "model"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"input_cents_per1k":  payload.InputCentsPer1K,
			"output_cents_per1k": payload.OutputCentsPer1K,
			"updated_at":         time.Now().UTC(),
		}),
	}).Create(&record).Error
	lvn.GinErr(c, 400, err, "unable to save pricing")
	if err != nil {
		return
	}

	resp := pricingResponse{
		Model:            record.Model,
		InputCentsPer1K:  record.InputCentsPer1K,
		OutputCentsPer1K: record.OutputCentsPer1K,
		UpdatedAt:        time.Now().UTC().Unix(),
	}

	c.Data(lvn.Res(200, resp, "OK"))
}

func UpsertPricingBatch(c *gin.Context) {
	payload := []pricingBatchItem{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")
	if err != nil {
		return
	}

	if len(payload) == 0 {
		lvn.GinErr(c, 400, nil, "pricing list is required")
		return
	}

	orderedModels := make([]string, 0, len(payload))
	itemsByModel := make(map[string]pricingBatchItem, len(payload))
	for _, item := range payload {
		model := strings.TrimSpace(item.Model)
		if model == "" {
			lvn.GinErr(c, 400, nil, "model is required")
			return
		}
		if item.InputCentsPer1K < 0 || item.OutputCentsPer1K < 0 {
			lvn.GinErr(c, 400, nil, "pricing must be non-negative")
			return
		}

		if _, exists := itemsByModel[model]; !exists {
			orderedModels = append(orderedModels, model)
		}
		item.Model = model
		itemsByModel[model] = item
	}

	now := time.Now().UTC()
	records := make([]models.OpenAIModelPricing, 0, len(orderedModels))
	for _, model := range orderedModels {
		item := itemsByModel[model]
		records = append(records, models.OpenAIModelPricing{
			Model:            model,
			InputCentsPer1K:  item.InputCentsPer1K,
			OutputCentsPer1K: item.OutputCentsPer1K,
			UpdatedAt:        now,
		})
	}

	err = db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "model"}},
		DoUpdates: clause.AssignmentColumns([]string{"input_cents_per1_k", "output_cents_per1_k", "updated_at"}),
	}).Create(&records).Error
	lvn.GinErr(c, 400, err, "unable to save pricing")
	if err != nil {
		return
	}

	resp := make([]pricingResponse, 0, len(orderedModels))
	for _, model := range orderedModels {
		item := itemsByModel[model]
		resp = append(resp, pricingResponse{
			Model:            model,
			InputCentsPer1K:  item.InputCentsPer1K,
			OutputCentsPer1K: item.OutputCentsPer1K,
			UpdatedAt:        now.Unix(),
		})
	}

	c.Data(lvn.Res(200, resp, "OK"))
}
