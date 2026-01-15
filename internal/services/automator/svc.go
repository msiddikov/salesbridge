package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	automationRunFilter struct {
		automationId    string
		status          string
		startedAfter    *time.Time
		startedBefore   *time.Time
		searchQuery     string
		batchRunID      string
		limit           int
		offset          int
		preloadNodeRuns bool
	}
)

func GetAutomationRuns(c *gin.Context) {
	filter := getAutomationRunFilterFromContext(c)

	if filter.limit > 100 {
		filter.limit = 100
	}

	runs, total, err := getAutomationRuns(filter)
	lvn.GinErr(c, 500, err, "Error getting automation runs")

	page := (filter.offset / filter.limit) + 1
	limit := filter.limit
	offset := filter.offset

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

func getAutomationRuns(filter automationRunFilter) ([]models.AutomationRun, int64, error) {
	runs := []models.AutomationRun{}

	query := db.DB.Model(&models.AutomationRun{})
	if filter.automationId != "" {
		query = query.Where("automation_id = ?", filter.automationId)
	}

	query = applyRunFilters(query, filter).
		Order("created_at desc").
		Limit(filter.limit).
		Offset(filter.offset)

	err := query.Find(&runs).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := db.DB.Model(&models.AutomationRun{})
	if filter.automationId != "" {
		countQuery = countQuery.Where("automation_id = ?", filter.automationId)
	}
	countQuery = applyRunFilters(countQuery, filter)
	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return runs, total, nil
}

func getAutomationRunFilterFromContext(c *gin.Context) automationRunFilter {
	automationId := c.Query("automationId")

	const (
		defaultLimit = 20
	)

	limit := parseQueryInt(c, "limit", defaultLimit)
	if limit <= 0 {
		limit = defaultLimit
	}

	page := parseQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	statusFilter := c.Query("status")
	searchQuery := c.Query("query")
	batchRunID := c.Query("batchRunId")
	startedAfter, _ := parseQueryTime(c, "startedAfter")
	startedBefore, _ := parseQueryTime(c, "startedBefore")

	return automationRunFilter{
		automationId:  automationId,
		status:        statusFilter,
		startedAfter:  startedAfter,
		startedBefore: startedBefore,
		searchQuery:   searchQuery,
		batchRunID:    batchRunID,
		limit:         limit,
		offset:        offset,
	}
}

func ExportAutomationRuns(c *gin.Context) {
	filter := getAutomationRunFilterFromContext(c)
	filter.preloadNodeRuns = true

	filter.limit = 1000000
	filter.offset = 0

	runs, _, err := getAutomationRuns(filter)
	lvn.GinErr(c, 500, err, "Error getting automation runs")

	batchStatuses, err := fetchBatchRunStatuses(runs)
	lvn.GinErr(c, 500, err, "Error getting batch run info")

	triggerKeys := collectTriggerPayloadKeys(runs)

	headers := []string{
		"startedDate",
		"endedDate",
		"duration",
		"status",
		"nodeRunsQty",
		"nodeRuns",
		"Status",
		"ExecutionsWithErrors",
		"errorMessage",
	}
	headers = append(headers, triggerKeys...)

	filename := "automation-runs.csv"

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Status(200)

	writer := csv.NewWriter(c.Writer)
	if err := writer.Write(headers); err != nil {
		writer.Flush()
		lvn.GinErr(c, 500, err, "Error writing CSV header")
		return
	}

	for _, run := range runs {
		started := run.StartedAt.Format(time.RFC3339)
		ended := ""
		duration := ""
		if run.CompletedAt != nil {
			ended = run.CompletedAt.Format(time.RFC3339)
			duration = run.CompletedAt.Sub(run.StartedAt).String()
		}

		sortedNodes := append([]models.AutomationRunNode(nil), run.RunNodes...)
		sort.Slice(sortedNodes, func(i, j int) bool {
			return sortedNodes[i].Sequence < sortedNodes[j].Sequence
		})

		nodeNames := make([]string, 0, len(sortedNodes))
		for _, nodeRun := range sortedNodes {
			nodeNames = append(nodeNames, nodeRunDisplayName(nodeRun))
		}

		batchStatus := ""
		if run.BatchRunID != nil {
			if val, ok := batchStatuses[*run.BatchRunID]; ok {
				batchStatus = val
			}
		}

		row := []string{
			started,
			ended,
			duration,
			string(run.Status),
			strconv.Itoa(len(sortedNodes)),
			strings.Join(nodeNames, ", "),
			batchStatus,
			strconv.Itoa(run.RunNodesWithErrors),
			run.ErrorMessage,
		}

		for _, key := range triggerKeys {
			val := ""
			if run.TriggerPayload != nil {
				if v, ok := run.TriggerPayload[key]; ok {
					val = fmt.Sprintf("%v", v)
				}
			}
			row = append(row, val)
		}

		if err := writer.Write(row); err != nil {
			writer.Flush()
			lvn.GinErr(c, 500, err, "Error writing CSV row")
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		lvn.GinErr(c, 500, err, "Error finalizing CSV export")
		return
	}
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

func StartTriggerForAutomation(c *gin.Context) {
	automationId := c.Param("automationId")

	var automation models.Automation
	err := db.DB.First(&automation, "id = ?", automationId).Preload("Location").Error
	lvn.GinErr(c, 400, err, "Could not retrieve automation")

	var input TriggerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		lvn.GinErr(c, 400, err, "Invalid input")
		return
	}

	input.LocationID = automation.LocationId

	err = StartAutomationForOneTrigger(context.Background(), automation, input)
	lvn.GinErr(c, 500, err, "Error starting automation from trigger")

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

func applyRunFilters(tx *gorm.DB, filter automationRunFilter) *gorm.DB {
	if filter.batchRunID != "" {
		tx = tx.Where("batch_run_id = ?", filter.batchRunID)
	} else {
		tx = tx.Where("batch_run_id IS NULL")
	}
	if filter.status != "" && filter.status != "all" {
		tx = tx.Where("status = ?", filter.status)
	}
	if filter.startedAfter != nil {
		tx = tx.Where("started_at >= ?", *filter.startedAfter)
	}
	if filter.startedBefore != nil {
		tx = tx.Where("started_at <= ?", *filter.startedBefore)
	}
	if filter.searchQuery != "" {
		tx = tx.Where("trigger_payload::text ILIKE ?", "%"+filter.searchQuery+"%")
	}
	if filter.automationId != "" {
		tx = tx.Where("automation_id = ?", filter.automationId)
	}
	if filter.preloadNodeRuns {
		tx = tx.Preload("RunNodes")
	}
	return tx
}

func fetchBatchRunStatuses(runs []models.AutomationRun) (map[string]string, error) {
	idsSet := make(map[string]struct{})
	for _, run := range runs {
		if run.BatchRunID != nil && *run.BatchRunID != "" {
			idsSet[*run.BatchRunID] = struct{}{}
		}
	}
	if len(idsSet) == 0 {
		return map[string]string{}, nil
	}
	ids := make([]string, 0, len(idsSet))
	for id := range idsSet {
		ids = append(ids, id)
	}
	var batchRuns []models.AutomationBatchRun
	if err := db.DB.Where("id IN ?", ids).Find(&batchRuns).Error; err != nil {
		return nil, err
	}
	res := make(map[string]string, len(batchRuns))
	for _, batch := range batchRuns {
		res[batch.ID] = string(batch.Status)
	}
	return res, nil
}

func collectTriggerPayloadKeys(runs []models.AutomationRun) []string {
	keySet := make(map[string]struct{})
	for _, run := range runs {
		if run.TriggerPayload == nil {
			continue
		}
		for key := range run.TriggerPayload {
			if key == "" {
				continue
			}
			keySet[key] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func nodeRunDisplayName(node models.AutomationRunNode) string {
	name := strings.TrimSpace(node.NodeName)
	if name != "" {
		return name
	}
	if node.NodeType != "" {
		return strings.TrimSpace(node.NodeType)
	}
	if node.NodeID != "" {
		return node.NodeID
	}
	return "Unnamed node"
}
