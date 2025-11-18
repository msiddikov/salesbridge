package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAutomationRuns(c *gin.Context) {
	automationId := c.Param("automationId")

	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit := parseQueryInt(c, "limit", defaultLimit)
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page := parseQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	statusFilter := c.Query("status")
	searchQuery := c.Query("query")
	batchRunID := c.Query("batchRunId")
	startedAfter, err := parseQueryTime(c, "startedAfter")
	if err != nil {
		lvn.GinErr(c, 400, err, "Invalid startedAfter parameter")
		return
	}
	startedBefore, err := parseQueryTime(c, "startedBefore")
	if err != nil {
		lvn.GinErr(c, 400, err, "Invalid startedBefore parameter")
		return
	}

	runs := []models.AutomationRun{}

	query := db.DB.
		Where("automation_id = ?", automationId)
	query = applyRunFilters(query, statusFilter, startedAfter, startedBefore, searchQuery, batchRunID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset)

	err = query.Find(&runs).Error
	lvn.GinErr(c, 500, err, "Error getting automation runs")

	var total int64
	countQuery := db.DB.Model(&models.AutomationRun{}).
		Where("automation_id = ?", automationId)
	countQuery = applyRunFilters(countQuery, statusFilter, startedAfter, startedBefore, searchQuery, batchRunID)
	err = countQuery.Count(&total).Error
	lvn.GinErr(c, 500, err, "Error counting automation runs")

	response := gin.H{
		"runs": runs,
		"pagination": gin.H{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasMore": int64(offset+len(runs)) < total,
		},
	}

	c.Data(lvn.Res(200, response, "OK"))
}

func GetAutomationRunDetails(c *gin.Context) {
	runId := c.Param("runId")

	var run models.AutomationRun
	err := db.DB.First(&run, "id = ?", runId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lvn.GinErr(c, 404, err, "Automation run not found")
		return
	}
	lvn.GinErr(c, 500, err, "Error getting automation run")

	const (
		defaultLimit = 25
		maxLimit     = 200
	)

	limit := parseQueryInt(c, "limit", defaultLimit)
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page := parseQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	var nodes []models.AutomationRunNode
	err = db.DB.
		Where("run_id = ?", run.ID).
		Order("sequence ASC").
		Limit(limit).
		Offset(offset).
		Find(&nodes).Error
	lvn.GinErr(c, 500, err, "Error getting automation run nodes")

	var total int64
	err = db.DB.Model(&models.AutomationRunNode{}).
		Where("run_id = ?", run.ID).
		Count(&total).Error
	lvn.GinErr(c, 500, err, "Error counting automation run nodes")

	response := gin.H{
		"run":   run,
		"nodes": nodes,
		"pagination": gin.H{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasMore": int64(offset+len(nodes)) < total,
		},
	}

	c.Data(lvn.Res(200, response, "OK"))
}

func StartFromAutomationRun(c *gin.Context) {
	runId := c.Param("runId")

	var run models.AutomationRun
	err := db.DB.First(&run, "id = ?", runId).Error
	lvn.GinErr(c, 400, err, "Could not retrieve automation")

	TriggerInput := TriggerInput{
		LocationID:  run.LocationID,
		TriggerType: run.TriggerType,
		Port:        "out",
		Payload:     run.TriggerPayload,
	}

	err = StartAutomationsForTrigger(context.Background(), TriggerInput)
	lvn.GinErr(c, 500, err, "Error starting automation from run")

	c.Data(lvn.Res(200, "Automation started", ""))
}

func parseQueryInt(c *gin.Context, key string, fallback int) int {
	val := c.Query(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseQueryTime(c *gin.Context, key string) (*time.Time, error) {
	val := c.Query(key)
	if val == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, fmt.Errorf("%s must be RFC3339", key)
	}
	return &parsed, nil
}

func applyRunFilters(tx *gorm.DB, status string, startedAfter, startedBefore *time.Time, searchQuery, batchRunID string) *gorm.DB {
	if batchRunID != "" {
		tx = tx.Where("batch_run_id = ?", batchRunID)
	} else {
		tx = tx.Where("batch_run_id IS NULL")
	}
	if status != "" && status != "all" {
		tx = tx.Where("status = ?", status)
	}
	if startedAfter != nil {
		tx = tx.Where("started_at >= ?", *startedAfter)
	}
	if startedBefore != nil {
		tx = tx.Where("started_at <= ?", *startedBefore)
	}
	if searchQuery != "" {
		tx = tx.Where("trigger_payload::text ILIKE ?", "%"+searchQuery+"%")
	}
	return tx
}
