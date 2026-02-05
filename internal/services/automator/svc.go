package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_cerbo"
	"client-runaway-zenoti/internal/services/svc_googleads"
	"client-runaway-zenoti/internal/services/svc_openai"
	"client-runaway-zenoti/packages/grafana"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		nodesCountFrom  *int
		nodesCountTo    *int
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

	// Get node counts for all runs in a single query
	if len(runs) > 0 {
		runIDs := make([]string, len(runs))
		for i, run := range runs {
			runIDs[i] = run.ID
		}

		type CountResult struct {
			RunID string
			Count int64
		}
		var counts []CountResult
		err := db.DB.Model(&models.AutomationRunNode{}).
			Select("run_id, COUNT(*) as count").
			Where("run_id IN ?", runIDs).
			Group("run_id").
			Scan(&counts).Error

		if err == nil {
			countMap := make(map[string]int)
			for _, c := range counts {
				countMap[c.RunID] = int(c.Count)
			}

			for i := range runs {
				runs[i].RunNodesQty = countMap[runs[i].ID]
			}
		}
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

	var nodesCountFrom, nodesCountTo *int
	if v := c.Query("nodesCountFrom"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			nodesCountFrom = &parsed
		}
	}
	if v := c.Query("nodesCountTo"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			nodesCountTo = &parsed
		}
	}

	return automationRunFilter{
		automationId:   automationId,
		status:         statusFilter,
		startedAfter:   startedAfter,
		startedBefore:  startedBefore,
		searchQuery:    searchQuery,
		batchRunID:     batchRunID,
		limit:          limit,
		offset:         offset,
		nodesCountFrom: nodesCountFrom,
		nodesCountTo:   nodesCountTo,
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
	lvn.GinErr(c, 400, err, "Could not retrieve automation run")

	// Check if this is a batch run (collection-based) or trigger-based run
	if run.BatchRunID != nil {
		// Collection-based run - need to start from collection node
		restartCollectionRun(c, run)
		return
	}

	// Trigger-based run - use the original logic
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

func restartCollectionRun(c *gin.Context, originalRun models.AutomationRun) {
	// Find the batch run to get the node ID
	var batchRun models.AutomationBatchRun
	err := db.DB.First(&batchRun, "id = ?", *originalRun.BatchRunID).Error
	if err != nil {
		lvn.GinErr(c, 404, err, "Batch run not found")
		return
	}

	// Load the node with automation
	var dbNode models.Node
	err = db.DB.
		Preload("Automation").
		Preload("Automation.Location").
		Preload("Automation.Location.ZenotiApiObj").
		Where("id = ?", batchRun.NodeID).
		First(&dbNode).Error
	if err != nil {
		lvn.GinErr(c, 404, err, "Collection node not found")
		return
	}

	// Find the API node from the graph
	var node models.APINode
	for _, n := range dbNode.Automation.Graph.Nodes {
		if n.ID == dbNode.ID {
			node = n
			break
		}
	}

	catalogNode, ok := getCatalogNode(node.Type)
	if !ok {
		lvn.GinErr(c, 500, errors.New("catalog node not found"), "Catalog node not found for type: "+node.Type)
		return
	}

	// Create new automation run
	ctx := context.Background()
	automation := dbNode.Automation

	payload := originalRun.TriggerPayload
	payloads := make(map[string]map[string]interface{})
	portName := originalRun.TriggerPort
	if portName == "" && len(catalogNode.Ports) > 0 {
		portName = catalogNode.Ports[0].Name
	}
	payloads[portName] = payload

	runtime := newAutomationRuntime(automation)
	runtime.runStatus = &models.AutomationRun{
		ID:             uuid.New().String(),
		AutomationID:   automation.ID,
		BatchRunID:     originalRun.BatchRunID,
		LocationID:     automation.LocationId,
		Status:         models.RunRunning,
		TriggerType:    node.Type,
		TriggerPayload: payload,
		TriggerPort:    portName,
		StartedAt:      time.Now(),
		RunNodes:       []models.AutomationRunNode{},
	}

	db.DB.Save(&runtime.runStatus)

	locName := dbNode.Automation.Location.Name
	locId := dbNode.Automation.LocationId

	// Run in background
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				errMsg := fmt.Sprintf("automation run panic: %v\n%s", r, stack)
				log.Printf("PANIC in automation run %s: %v\n%s", runtime.runStatus.ID, r, stack)

				grafana.Notify(locName, locId, "automation-run-error", errMsg)

				finishedAt := time.Now()
				runtime.runStatus.CompletedAt = &finishedAt
				runtime.runStatus.Status = models.RunFailed
				runtime.runStatus.ErrorMessage = fmt.Sprintf("panic: %v", r)
				db.DB.Save(&runtime.runStatus)
			}
		}()

		runResultErr := runtime.startFromEntry(ctx, node, payloads)
		finishedAt := time.Now()

		runtime.runStatus.CompletedAt = &finishedAt
		if runResultErr != nil {
			runtime.runStatus.Status = models.RunFailed
			runtime.runStatus.ErrorMessage = runResultErr.Error()
		} else {
			runtime.runStatus.Status = models.RunSuccess
		}
		db.DB.Save(&runtime.runStatus)
	}()

	c.Data(lvn.Res(200, runtime.runStatus, "Automation run restarted"))
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
	if filter.nodesCountFrom != nil {
		tx = tx.Where("(SELECT COUNT(*) FROM automation_run_nodes WHERE automation_run_nodes.run_id = automation_runs.id) >= ?", *filter.nodesCountFrom)
	}
	if filter.nodesCountTo != nil {
		tx = tx.Where("(SELECT COUNT(*) FROM automation_run_nodes WHERE automation_run_nodes.run_id = automation_runs.id) <= ?", *filter.nodesCountTo)
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

func GetLists(c *gin.Context) {

	listName := c.Param("listName")
	locationId := c.Param("locationId")

	var location models.Location
	err := db.DB.Preload("CerboApiObj").First(&location, "id = ?", locationId).Error
	lvn.GinErr(c, 400, err, "Could not find location")

	switch listName {
	case "cerboUsers":
		users, err := svc_cerbo.ListUsers(location)
		lvn.GinErr(c, 500, err, "failed to list cerbo users")

		c.Data(lvn.Res(200, users, "OK"))
		return
	case "googleAdsActions":
		list, err := svc_googleads.GetLocationConversionActionsList(location)
		lvn.GinErr(c, 500, err, "failed to list google ads conversion actions")

		c.Data(lvn.Res(200, list, "OK"))
		return
	case "aiAssistants":
		list, err := svc_openai.GetAssistantsList(location)
		lvn.GinErr(c, 500, err, "failed to list ai assistants")

		c.Data(lvn.Res(200, list, "OK"))
		return
	case "cerboEncounterTypes":
		types, err := svc_cerbo.GetEncounterTypesList(location)
		lvn.GinErr(c, 500, err, "failed to list cerbo encounter types")

		c.Data(lvn.Res(200, types, "OK"))
		return
	case "cerboFreeTextTypes":
		types, err := svc_cerbo.ListFreeTextNoteTypes(location)
		lvn.GinErr(c, 500, err, "failed to list cerbo free text types")

		c.Data(lvn.Res(200, types, "OK"))
		return
	}

	lvn.GinErr(c, 400, fmt.Errorf("unknown list name %q", listName), "unknown list name")
}

// GetCatalogData returns the automation node catalog for internal MCP tools
func GetCatalogData() Catalog {
	cat := designCatalog(catalogFull)
	return getCatalogWithImplementedNodes(cat)
}

// GetCatalogListData returns list data for dynamic node fields
func GetCatalogListData(listName string, location models.Location) (interface{}, error) {
	switch listName {
	case "cerboUsers":
		return svc_cerbo.ListUsers(location)
	case "googleAdsActions":
		return svc_googleads.GetLocationConversionActionsList(location)
	case "aiAssistants":
		return svc_openai.GetAssistantsList(location)
	case "cerboEncounterTypes":
		return svc_cerbo.GetEncounterTypesList(location)
	case "cerboFreeTextTypes":
		return svc_cerbo.ListFreeTextNoteTypes(location)
	default:
		return nil, fmt.Errorf("unknown list name: %s", listName)
	}
}

// AutomationInfo represents automation without sensitive data
type AutomationInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LocationID  string `json:"locationId"`
	State       string `json:"state"`
	NodeCount   int    `json:"nodeCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// AutomationRunInfo represents a run with node counts
type AutomationRunInfo struct {
	ID              string  `json:"id"`
	AutomationID    string  `json:"automationId"`
	BatchRunID      *string `json:"batchRunId,omitempty"`
	Status          string  `json:"status"`
	TriggerType     string  `json:"triggerType"`
	ErrorMessage    string  `json:"errorMessage,omitempty"`
	NodesExecuted   int     `json:"nodesExecuted"`
	NodesWithErrors int     `json:"nodesWithErrors"`
	StartedAt       string  `json:"startedAt"`
	CompletedAt     string  `json:"completedAt,omitempty"`
}

// BatchRunInfo represents a batch run
type BatchRunInfo struct {
	ID             string  `json:"id"`
	AutomationID   string  `json:"automationId"`
	NodeID         string  `json:"nodeId"`
	Status         string  `json:"status"`
	ItemsProcessed int     `json:"itemsProcessed"`
	TotalItems     *int    `json:"totalItems,omitempty"`
	ProgressPct    float64 `json:"progressPct"`
	ErrorMessage   string  `json:"errorMessage,omitempty"`
	StartedAt      string  `json:"startedAt,omitempty"`
	CompletedAt    string  `json:"completedAt,omitempty"`
}

// GetAutomationsForProfile returns automations for a profile, optionally filtered by location and state
func GetAutomationsForProfile(profileID uint, locationID, state string) ([]AutomationInfo, error) {
	query := db.DB.Model(&models.Automation{}).
		Joins("JOIN locations ON locations.id = automations.location_id").
		Where("locations.profile_id = ?", profileID)

	if locationID != "" {
		query = query.Where("automations.location_id = ?", locationID)
	}
	if state != "" {
		query = query.Where("automations.state = ?", state)
	}

	var automations []models.Automation
	err := query.Order("automations.created_at DESC").Find(&automations).Error
	if err != nil {
		return nil, err
	}

	result := make([]AutomationInfo, len(automations))
	for i, auto := range automations {
		result[i] = AutomationInfo{
			ID:          auto.ID,
			Name:        auto.Name,
			Description: auto.Description,
			LocationID:  auto.LocationId,
			State:       string(auto.State),
			NodeCount:   len(auto.Graph.Nodes),
			CreatedAt:   auto.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   auto.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return result, nil
}

// GetAutomationForProfile returns a single automation for a profile
func GetAutomationForProfile(profileID uint, automationID string) (*models.Automation, error) {
	var automation models.Automation
	err := db.DB.
		Joins("JOIN locations ON locations.id = automations.location_id").
		Where("automations.id = ? AND locations.profile_id = ?", automationID, profileID).
		First(&automation).Error
	if err != nil {
		return nil, err
	}
	return &automation, nil
}

// GetAutomationRunsForProfile returns automation runs for a profile
func GetAutomationRunsForProfile(profileID uint, automationID, status string, limit int) ([]AutomationRunInfo, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := db.DB.Model(&models.AutomationRun{}).
		Joins("JOIN automations ON automations.id = automation_runs.automation_id").
		Joins("JOIN locations ON locations.id = automations.location_id").
		Where("locations.profile_id = ?", profileID)

	if automationID != "" {
		query = query.Where("automation_runs.automation_id = ?", automationID)
	}
	if status != "" {
		query = query.Where("automation_runs.status = ?", status)
	}

	var runs []models.AutomationRun
	err := query.
		Order("automation_runs.started_at DESC").
		Limit(limit).
		Preload("RunNodes").
		Find(&runs).Error
	if err != nil {
		return nil, err
	}

	result := make([]AutomationRunInfo, len(runs))
	for i, run := range runs {
		info := AutomationRunInfo{
			ID:              run.ID,
			AutomationID:    run.AutomationID,
			BatchRunID:      run.BatchRunID,
			Status:          string(run.Status),
			TriggerType:     run.TriggerType,
			ErrorMessage:    run.ErrorMessage,
			NodesExecuted:   len(run.RunNodes),
			NodesWithErrors: run.RunNodesWithErrors,
			StartedAt:       run.StartedAt.Format("2006-01-02T15:04:05Z"),
		}
		if run.CompletedAt != nil {
			info.CompletedAt = run.CompletedAt.Format("2006-01-02T15:04:05Z")
		}
		result[i] = info
	}
	return result, nil
}

// GetAutomationRunForProfile returns a single automation run for a profile
func GetAutomationRunForProfile(profileID uint, runID string) (*models.AutomationRun, error) {
	var run models.AutomationRun
	err := db.DB.
		Joins("JOIN automations ON automations.id = automation_runs.automation_id").
		Joins("JOIN locations ON locations.id = automations.location_id").
		Where("automation_runs.id = ? AND locations.profile_id = ?", runID, profileID).
		Preload("RunNodes").
		First(&run).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

// GetBatchRunsForLocation returns batch runs for a location (after verifying profile ownership)
func GetBatchRunsForLocation(profileID uint, locationID, status string) ([]BatchRunInfo, error) {
	// Verify location belongs to profile
	var location models.Location
	err := db.DB.Where("id = ? AND profile_id = ?", locationID, profileID).First(&location).Error
	if err != nil {
		return nil, fmt.Errorf("location not found or access denied")
	}

	query := db.DB.Model(&models.AutomationBatchRun{}).Where("location_id = ?", locationID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var batchRuns []models.AutomationBatchRun
	err = query.Order("created_at DESC").Limit(50).Find(&batchRuns).Error
	if err != nil {
		return nil, err
	}

	result := make([]BatchRunInfo, len(batchRuns))
	for i, br := range batchRuns {
		info := BatchRunInfo{
			ID:             br.ID,
			AutomationID:   br.AutomationID,
			NodeID:         br.NodeID,
			Status:         string(br.Status),
			ItemsProcessed: br.ItemsProcessed,
			TotalItems:     br.TotalItems,
			ProgressPct:    br.ProgressPct,
			ErrorMessage:   br.ErrorMessage,
		}
		if br.StartedAt != nil {
			info.StartedAt = br.StartedAt.Format("2006-01-02T15:04:05Z")
		}
		if br.CompletedAt != nil {
			info.CompletedAt = br.CompletedAt.Format("2006-01-02T15:04:05Z")
		}
		result[i] = info
	}
	return result, nil
}
