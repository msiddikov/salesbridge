package svc_openai

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"fmt"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type usageResponse struct {
	Date         string  `json:"date"`
	Model        string  `json:"model"`
	InputTokens  int     `json:"inputTokens"`
	OutputTokens int     `json:"outputTokens"`
	CostCents    float64 `json:"costCents"`
	Points       float64 `json:"points"`
}

func ListUsage(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	start, end, err := usageRange(c)
	lvn.GinErr(c, 400, err, "invalid date range")
	if err != nil {
		return
	}

	query := db.DB.Where("profile_id = ? AND usage_date >= ? AND usage_date < ?", user.ProfileID, start, end)
	if model := c.Query("model"); model != "" {
		query = query.Where("model = ?", model)
	}

	var usage []models.OpenAIUsageDaily
	err = query.Order("usage_date asc").Find(&usage).Error
	lvn.GinErr(c, 400, err, "unable to list usage")
	if err != nil {
		return
	}

	resp := make([]usageResponse, 0, len(usage))
	for _, u := range usage {
		resp = append(resp, usageResponse{
			Date:         u.UsageDate.Format("2006-01-02"),
			Model:        u.Model,
			InputTokens:  u.InputTokens,
			OutputTokens: u.OutputTokens,
			CostCents:    u.CostCents,
			Points:       u.Points,
		})
	}

	c.Data(lvn.Res(200, resp, "OK"))
}

func usageRange(c *gin.Context) (time.Time, time.Time, error) {
	if month := c.Query("month"); month != "" {
		start, err := time.Parse("2006-01", month)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start, end, nil
	}

	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" && endStr == "" {
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start, end, nil
	}
	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("start and end are required")
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	return start, end, nil
}
