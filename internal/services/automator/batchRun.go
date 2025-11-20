package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var batchRunCancels sync.Map

func StartBatchRun(c *gin.Context) {
	batchRun := models.AutomationBatchRun{}
	err := c.BindJSON(&batchRun)
	lvn.GinErr(c, 500, err, "Error binding JSON")

	// create batch run record
	now := time.Now()
	batchRun.StartedAt = &now
	batchRun.Status = models.BatchRunRunning
	batchRun.ID = uuid.New().String()
	batchRun.CreatorID = c.MustGet("user").(models.User).ID

	nodeId := batchRun.NodeID

	dbNode := models.Node{}

	err = db.DB.
		Preload("Automation").
		Preload("Automation.Location").
		Preload("Automation.Location.ZenotiApiObj").
		Where("id = ?", nodeId).
		First(&dbNode).Error
	lvn.GinErr(c, 404, err, "Node not found")

	batchRun.LocationID = dbNode.Automation.LocationId
	err = db.DB.WithContext(c.Request.Context()).Create(&batchRun).Error
	lvn.GinErr(c, 500, err, "Error creating batch run record")

	ctx, cancel := context.WithCancel(context.Background())
	registerBatchRunCancel(batchRun.ID, cancel)

	go func() {
		defer func() {
			cancel()
			unregisterBatchRunCancel(batchRun.ID)
		}()
		StartAutomationsForCollection(ctx, dbNode, batchRun)
	}()
	c.Data(lvn.Res(200, batchRun, "Batch run started"))

}

func GetBatchRuns(c *gin.Context) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	locationId := c.Param("locationId")
	if locationId == "" {
		lvn.GinErr(c, 400, errors.New("locationId required"), "locationId is required")
		return
	}

	limit := parseBatchQueryInt(c, "limit", defaultLimit)
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page := parseBatchQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	status := c.Query("status")
	automationId := c.Query("automationId")
	nodeId := c.Query("nodeId")

	startedAfter, err := parseBatchQueryTime(c, "startedAfter")
	if err != nil {
		lvn.GinErr(c, 400, err, "Invalid startedAfter parameter")
		return
	}
	startedBefore, err := parseBatchQueryTime(c, "startedBefore")
	if err != nil {
		lvn.GinErr(c, 400, err, "Invalid startedBefore parameter")
		return
	}

	query := db.DB.Model(&models.AutomationBatchRun{}).Where("location_id = ?", locationId)
	query = applyBatchRunFilters(query, status, automationId, nodeId, startedAfter, startedBefore)

	var batchRuns []models.AutomationBatchRun
	err = query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&batchRuns).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Data(lvn.Res(200, gin.H{
			"batchRuns": []models.AutomationBatchRun{},
			"pagination": gin.H{
				"page":    page,
				"limit":   limit,
				"total":   0,
				"hasMore": false,
			},
		}, "OK"))
		return
	}
	lvn.GinErr(c, 500, err, "Error getting batch runs")

	var total int64
	err = query.Count(&total).Error
	lvn.GinErr(c, 500, err, "Error counting batch runs")

	response := gin.H{
		"data": gin.H{

			"batchRuns": batchRuns,
			"pagination": gin.H{
				"page":    page,
				"limit":   limit,
				"total":   total,
				"hasMore": int64(offset+len(batchRuns)) < total,
			},
		},
		"message": "success",
		"id":      true,
	}

	c.JSON(200, response)
}

func GetBatchRunDetails(c *gin.Context) {
	batchRunID := c.Param("batchRunId")
	if batchRunID == "" {
		lvn.GinErr(c, 400, errors.New("batchRunId required"), "batchRunId is required")
		return
	}

	var batchRun models.AutomationBatchRun
	err := db.DB.First(&batchRun, "id = ?", batchRunID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lvn.GinErr(c, 404, err, "Batch run not found")
		return
	}
	lvn.GinErr(c, 500, err, "Error getting batch run")

	const (
		defaultLimit = 25
		maxLimit     = 200
	)

	limit := parseBatchQueryInt(c, "limit", defaultLimit)
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page := parseBatchQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	statusFilter := c.Query("status")
	queryFilter := c.Query("query")

	var runs []models.AutomationRun
	runQuery := db.DB.
		Where("batch_run_id = ?", batchRun.ID)
	runQuery = applyBatchRunDetailsFilters(runQuery, statusFilter, queryFilter)

	err = runQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&runs).Error
	lvn.GinErr(c, 500, err, "Error getting batch run items")

	var total int64
	countQuery := db.DB.Model(&models.AutomationRun{}).
		Where("batch_run_id = ?", batchRun.ID)
	countQuery = applyBatchRunDetailsFilters(countQuery, statusFilter, queryFilter)
	err = countQuery.Count(&total).Error
	lvn.GinErr(c, 500, err, "Error counting batch run items")

	response := gin.H{
		"batchRun": batchRun,
		"runs":     runs,
		"pagination": gin.H{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasMore": int64(offset+len(runs)) < total,
		},
	}

	c.Data(lvn.Res(200, response, "OK"))
}

func CancelBatchRun(c *gin.Context) {
	batchRunID := c.Param("batchRunId")
	if batchRunID == "" {
		lvn.GinErr(c, 400, errors.New("batchRunId required"), "batchRunId is required")
		return
	}

	var batchRun models.AutomationBatchRun
	err := db.DB.First(&batchRun, "id = ?", batchRunID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lvn.GinErr(c, 404, err, "Batch run not found")
		return
	}
	lvn.GinErr(c, 500, err, "Error getting batch run")

	terminal := batchRun.Status == models.BatchRunSuccess ||
		batchRun.Status == models.BatchRunFailed ||
		batchRun.Status == models.BatchRunCanceled
	cancelTriggered := false
	if !terminal {
		cancelTriggered = requestBatchRunCancel(batchRunID)
	}

	now := time.Now()
	update := map[string]interface{}{
		"status":        models.BatchRunCanceled,
		"completed_at":  &now,
		"error_message": "Canceled by user",
	}
	err = db.DB.Model(&batchRun).Updates(update).Error
	lvn.GinErr(c, 500, err, "Error updating batch run")

	batchRun.Status = models.BatchRunCanceled
	batchRun.CompletedAt = &now
	batchRun.ErrorMessage = "Canceled by user"

	message := "Batch run canceled"
	if terminal {
		message = "Batch run already completed"
	} else if !cancelTriggered {
		message = "Batch run cancel requested"
	}

	c.Data(lvn.Res(200, batchRun, message))
}

func parseBatchQueryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return val
}

func parseBatchQueryTime(c *gin.Context, key string) (*time.Time, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be RFC3339", key)
	}
	return &parsed, nil
}

func applyBatchRunFilters(tx *gorm.DB, status, automationId, nodeId string, startedAfter, startedBefore *time.Time) *gorm.DB {
	if status != "" && status != "all" {
		tx = tx.Where("status = ?", status)
	}
	if automationId != "" {
		tx = tx.Where("automation_id = ?", automationId)
	}
	if nodeId != "" {
		tx = tx.Where("node_id = ?", nodeId)
	}
	if startedAfter != nil {
		tx = tx.Where("(started_at IS NOT NULL AND started_at >= ?)", *startedAfter)
	}
	if startedBefore != nil {
		tx = tx.Where("(started_at IS NOT NULL AND started_at <= ?)", *startedBefore)
	}
	return tx
}

func applyBatchRunDetailsFilters(tx *gorm.DB, status, query string) *gorm.DB {
	if status != "" && status != "all" {
		tx = tx.Where("status = ?", status)
	}
	if query != "" {
		tx = tx.Where("trigger_payload::text ILIKE ?", "%"+query+"%")
	}
	return tx
}

func registerBatchRunCancel(batchRunID string, cancel context.CancelFunc) {
	if batchRunID == "" || cancel == nil {
		return
	}
	batchRunCancels.Store(batchRunID, cancel)
}

func requestBatchRunCancel(batchRunID string) bool {
	if batchRunID == "" {
		return false
	}
	if val, ok := batchRunCancels.Load(batchRunID); ok {
		if cancel, ok := val.(context.CancelFunc); ok {
			cancel()
			return true
		}
	}
	return false
}

func unregisterBatchRunCancel(batchRunID string) {
	if batchRunID == "" {
		return
	}
	batchRunCancels.Delete(batchRunID)
}
