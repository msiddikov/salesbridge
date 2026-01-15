package automator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	"github.com/google/uuid"
)

const (
	defaultPortOut     = "out"
	defaultPortSuccess = "success"
	defaultPortFail    = "fail"

	maxNodeExecutions = 512
)

type (
	// TriggerInput describes an automation trigger firing event.
	TriggerInput struct {
		LocationID  string
		TriggerType string
		Port        string
		Payload     map[string]interface{}
	}

	queuedNode struct {
		node           models.APINode
		incomingEdgeID string
		incomingPort   string
		parent         *queuedNode
	}

	edgeRef struct {
		id       string
		toNodeID string
	}

	automationRuntime struct {
		automation models.Automation
		nodes      map[string]models.APINode
		edges      map[string]map[string]edgeRef
		runStatus  *models.AutomationRun
	}

	collectionResult struct {
		items   []collectionItem
		total   int
		hasMore bool
	}

	collectionItem struct {
		payload   map[string]interface{}
		countsFor int
	}
)

// StartAutomationsForTrigger finds active automations for the given location whose
// entry trigger matches the provided trigger type and starts their execution.
func StartAutomationsForTrigger(ctx context.Context, input TriggerInput) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if input.LocationID == "" {
		return errors.New("automator: locationID is required")
	}
	if input.TriggerType == "" {
		return errors.New("automator: triggerType is required")
	}
	if input.Payload == nil {
		input.Payload = map[string]interface{}{}
	}

	query := db.DB.WithContext(ctx).
		Preload("Nodes").
		Preload("Edges").
		Preload("Location").
		Preload("Location.ZenotiApiObj").
		Where("location_id = ? AND state = ?", input.LocationID, models.StateActive)

	var automations []models.Automation
	if err := query.Find(&automations).Error; err != nil {
		return fmt.Errorf("automator: find automations: %w", err)
	}

	go func() {

		for _, automation := range automations {
			err := StartAutomationForOneTrigger(ctx, automation, input)
			if err != nil {
				log.Printf("automator: start automation for trigger: %s", err.Error())
			}
		}
	}()

	return nil
}

func StartAutomationForOneTrigger(ctx context.Context, automation models.Automation, input TriggerInput) error {
	if len(automation.Graph.Entry) == 0 || len(automation.Graph.Nodes) == 0 {
		return nil
	}

	runtime := newAutomationRuntime(automation)
	runtime.runStatus = &models.AutomationRun{
		ID:             uuid.New().String(),
		AutomationID:   automation.ID,
		LocationID:     automation.LocationId,
		Status:         models.RunRunning,
		TriggerType:    input.TriggerType,
		TriggerPort:    input.Port,
		TriggerPayload: input.Payload,
		StartedAt:      time.Now(),
		RunNodes:       []models.AutomationRunNode{},
	}

	entryNodes := runtime.entryNodesForType(input.TriggerType)
	if len(entryNodes) == 0 {
		return nil
	}

	db.DB.Save(&runtime.runStatus) // save runtime status only if the automation has entry nodes

	for _, entry := range entryNodes {
		payloads := make(map[string]map[string]interface{})
		payloads[input.Port] = clonePayload(input.Payload)
		if err := runtime.startFromEntry(ctx, entry, payloads); err != nil {
			runtime.runStatus.Status = models.RunWithErrors
			if runtime.runStatus.ErrorMessage == "" {
				runtime.runStatus.ErrorMessage = err.Error()
			}
		}
	}

	finishedAt := time.Now()
	runtime.runStatus.CompletedAt = &finishedAt
	if runtime.runStatus.RunNodesWithErrors == 0 && runtime.runStatus.Status != models.RunWithErrors {
		runtime.runStatus.Status = models.RunSuccess
	}
	db.DB.Save(&runtime.runStatus)
	return nil
}

func StartAutomationsForCollection(ctx context.Context, dbNode models.Node, batchRun models.AutomationBatchRun) {

	if ctx == nil {
		ctx = context.Background()
	}

	node := models.APINode{}

	for _, n := range dbNode.Automation.Graph.Nodes {
		if n.ID == dbNode.ID {
			node = n
			break
		}
	}

	var runErr error

	errorRunTime := newAutomationRuntime(dbNode.Automation)
	errorRunTime.runStatus = &models.AutomationRun{
		ID:           uuid.New().String(),
		AutomationID: dbNode.Automation.ID,
		BatchRunID:   &batchRun.ID,
		LocationID:   dbNode.Automation.LocationId,
		Status:       models.RunFailed,
		TriggerType:  node.Type,
		StartedAt:    time.Now(),
		RunNodes:     []models.AutomationRunNode{},
	}

	catalogNode, ok := getCatalogNode(node.Type)
	if !ok || catalogNode.CollectorFunc == nil {
		errorRunTime.runStatus.ErrorMessage = "collector function not found for node type"
		err := db.DB.Save(&errorRunTime.runStatus).Error
		if err != nil {
			log.Printf("automator: failed to save error runtime: %s", err.Error())
		}
		return
	}

	automation := dbNode.Automation

	nodeConfig := node.Config.EdgeConfig("")

	// for date-based pagination, ensure pages are set
	if _, ok := nodeConfig["pageFrom"]; !ok {
		nodeConfig["pageFrom"] = float64(1)
	}
	if _, ok := nodeConfig["pageTo"]; !ok {
		nodeConfig["pageTo"] = float64(0)
	}

	nodeConfig["page"] = nodeConfig["pageFrom"]
	limit, ok := nodeConfig["limit"].(float64)
	if !ok || limit < 1 {
		nodeConfig["limit"] = float64(20)
		err := db.DB.Save(&dbNode).Error
		if err != nil {
			errorRunTime.runStatus.ErrorMessage = fmt.Sprintf("automator: save node config: %s", err.Error())
			db.DB.Save(&errorRunTime.runStatus)
			return
		}
	}

	res, err := catalogNode.CollectorFunc(ctx, nodeConfig, automation.Location)
	if err != nil {
		errorRunTime.runStatus.ErrorMessage = fmt.Sprintf("automator: collect data for collection node: %s", err.Error())
		db.DB.Save(&errorRunTime.runStatus)
		return
	}

	itemsToProcess := getTotalItemsToProcess(res.total, int(nodeConfig["limit"].(float64)), int(nodeConfig["pageFrom"].(float64)), int(nodeConfig["pageTo"].(float64)))
	batchRun.TotalItems = &itemsToProcess
	pages := getTotalPagesToProcess(res.total, int(nodeConfig["limit"].(float64)), int(nodeConfig["pageFrom"].(float64)), int(nodeConfig["pageTo"].(float64)))
	batchRun.TotalPages = &pages
	batchRun.CurrentPage = int(nodeConfig["pageFrom"].(float64))
	batchRun.ItemsProcessed = 0
	batchRun.PageFrom = int(nodeConfig["pageFrom"].(float64))
	pageTo := int(nodeConfig["pageTo"].(float64))
	batchRun.PageTo = &pageTo
	db.DB.Save(&batchRun)

	for (len(res.items) == 0 && res.hasMore) || (len(res.items) > 0) || (nodeConfig["page"].(float64) < nodeConfig["pageTo"].(float64) && nodeConfig["pageTo"].(float64) != 0) {
		if ctx.Err() != nil {
			now := time.Now()
			batchRun.Status = models.BatchRunCanceled
			batchRun.ErrorMessage = "automator: batch run canceled"
			batchRun.CompletedAt = &now
			db.DB.Save(&batchRun)
			return
		}

		for _, item := range res.items {
			if ctx.Err() != nil {
				now := time.Now()
				batchRun.Status = models.BatchRunCanceled
				batchRun.ErrorMessage = "automator: batch run canceled"
				batchRun.CompletedAt = &now
				db.DB.Save(&batchRun)
				return
			}

			payload := clonePayload(item.payload)
			payloads := make(map[string]map[string]interface{})
			payloads[catalogNode.Ports[0].Name] = payload

			runtime := newAutomationRuntime(automation)
			runtime.runStatus = &models.AutomationRun{
				ID:           uuid.New().String(),
				AutomationID: automation.ID,
				BatchRunID:   &batchRun.ID,
				LocationID:   automation.LocationId,
				Status:       models.RunRunning,

				TriggerType:    node.Type,
				TriggerPayload: payload,
				TriggerPort:    catalogNode.Ports[0].Name,
				StartedAt:      time.Now(),
				RunNodes:       []models.AutomationRunNode{},
			}

			db.DB.Save(&runtime.runStatus)

			runResultErr := runtime.startFromEntry(ctx, node, payloads)
			finishedAt := time.Now()

			runtime.runStatus.CompletedAt = &finishedAt
			if runResultErr != nil {
				if errors.Is(runResultErr, context.Canceled) || ctx.Err() != nil {
					runtime.runStatus.Status = models.RunCanceled
					runtime.runStatus.ErrorMessage = "automation run canceled"
					db.DB.Save(&runtime.runStatus)

					batchRun.Status = models.BatchRunCanceled
					batchRun.ErrorMessage = "automator: batch run canceled"
					batchRun.CompletedAt = &finishedAt
					db.DB.Save(&batchRun)
					return
				}
				runtime.runStatus.Status = models.RunFailed
				runtime.runStatus.ErrorMessage = runResultErr.Error()
				runErr = errors.Join(runErr, fmt.Errorf("automation %s: %w", automation.ID, runResultErr))

				batchRun.Status = models.BatchRunWithErrors
				batchRun.RunsWithErrors++
				batchRun.ErrorMessage = fmt.Sprintf("%v runs finished with error", batchRun.RunsWithErrors)
			} else {
				runtime.runStatus.Status = models.RunSuccess
			}
			db.DB.Save(&runtime.runStatus)

			batchRun.ItemsProcessed += item.countsFor
			if batchRun.TotalItems != nil && *batchRun.TotalItems > 0 {
				batchRun.ProgressPct = float64(batchRun.ItemsProcessed) / float64(*batchRun.TotalItems) * 100.0
			}
			db.DB.Save(&batchRun)
		}

		nodeConfig["page"] = nodeConfig["page"].(float64) + 1
		if nodeConfig["page"].(float64) > nodeConfig["pageTo"].(float64) && nodeConfig["pageTo"].(float64) != 0 {
			break
		}
		res, err = catalogNode.CollectorFunc(ctx, nodeConfig, automation.Location)
		if err != nil {
			errorRunTime.runStatus.ErrorMessage = fmt.Sprintf("automator: collect data for collection node: %s", err.Error())
			db.DB.Save(&errorRunTime.runStatus)
			return
		}
		batchRun.CurrentPage = int(nodeConfig["page"].(float64))
		db.DB.Save(&batchRun)
	}

	finishedAt := time.Now()

	batchRun.CompletedAt = &finishedAt
	if runErr != nil {
		batchRun.Status = models.BatchRunFailed
	} else {
		batchRun.Status = models.BatchRunSuccess
	}
	db.DB.Save(&batchRun)
}

func newAutomationRuntime(auto models.Automation) *automationRuntime {
	rt := &automationRuntime{
		automation: auto,
		nodes:      make(map[string]models.APINode, len(auto.Graph.Nodes)),
		edges:      make(map[string]map[string]edgeRef),
	}

	for _, node := range auto.Graph.Nodes {
		rt.nodes[node.ID] = node
	}

	for _, edge := range auto.Graph.Edges {
		fromID := edge.FromNodeId
		if fromID == "" {
			continue
		}
		fromPort := edge.FromPort
		if fromPort == "" {
			if node, ok := rt.nodes[fromID]; ok {
				fromPort = defaultSuccessForKind(node.Kind)
			}
		}
		toID := edge.ToNodeId
		if toID == "" {
			continue
		}

		if _, ok := rt.edges[fromID]; !ok {
			rt.edges[fromID] = make(map[string]edgeRef)
		}
		rt.edges[fromID][fromPort] = edgeRef{
			id:       edge.ID,
			toNodeID: toID,
		}
	}

	return rt
}

func (rt *automationRuntime) entryNodesForType(triggerType string) []models.APINode {
	var entries []models.APINode
	for _, entryID := range rt.automation.Graph.Entry {
		node, ok := rt.nodes[entryID]
		if !ok {
			log.Printf("automator: automation %s references unknown entry node %s", rt.automation.ID, entryID)
			continue
		}
		if node.Type != triggerType {
			continue
		}
		entries = append(entries, node)
	}
	return entries
}

func (rt *automationRuntime) startFromEntry(ctx context.Context, entry models.APINode, payload map[string]map[string]interface{}) error {

	entryWrapper := &queuedNode{node: entry}
	queue := []*queuedNode{}
	queue = append(queue, rt.nextNodes(ctx, entryWrapper, payload)...)
	nodePayloads := make(map[string]map[string]map[string]interface{})
	nodePayloads[entry.ID] = payload

	executed := 0
	var runErr error

	for len(queue) > 0 {
		if ctx.Err() != nil {
			return errors.Join(runErr, ctx.Err())
		}
		if executed >= maxNodeExecutions {
			return errors.Join(runErr, fmt.Errorf("automation %s exceeded %d node executions", rt.automation.ID, maxNodeExecutions))
		}

		current := queue[0]
		queue = queue[1:]
		executed++

		currentNode := current.node
		effectiveConfig := currentNode.Config.EdgeConfig(current.incomingEdgeID)

		fieldValues := substNodeFields(current, effectiveConfig, nodePayloads)

		startTime := time.Now()

		results := executeNode(ctx, currentNode, fieldValues, rt.automation.Location)
		finishedTime := time.Now()

		if len(results) == 0 {
			results = errorPayload(nil, "Node didn't answered on any port")
		}

		runNode := models.AutomationRunNode{
			ID:             uuid.New().String(),
			RunID:          rt.runStatus.ID,
			NodeID:         currentNode.ID,
			NodeName:       currentNode.Name,
			NodeType:       currentNode.Type,
			Sequence:       executed,
			InputFields:    fieldValues,
			OutputPayloads: results,
			Status:         models.RunSuccess,
			StartedAt:      startTime,
			CompletedAt:    &finishedTime,
		}

		if errPayload, ok := results["error"]; ok {
			runNode.Status = models.RunFailed
			runNode.ErrorMessage = fmt.Sprintf("%v: %v", errPayload["message"], errPayload["error"])
			rt.runStatus.Status = models.RunWithErrors
			rt.runStatus.RunNodesWithErrors++
			runErr = fmt.Errorf("%v node(s) have errors", rt.runStatus.RunNodesWithErrors)
		}

		rt.runStatus.RunNodes = append(rt.runStatus.RunNodes, runNode)

		db.DB.Save(&runNode)

		if _, ok := nodePayloads[currentNode.ID]; !ok {
			nodePayloads[currentNode.ID] = make(map[string]map[string]interface{})
		}

		for port, resultPayload := range results {
			if resultPayload == nil {
				resultPayload = map[string]interface{}{}
			}
			nodePayloads[currentNode.ID][port] = resultPayload

			nextID, edgeID := rt.next(currentNode.ID, port)
			if nextID != "" {
				nextNode, ok := rt.nodes[nextID]
				if !ok {
					log.Printf("automator: automation %s references unknown node %s", rt.automation.ID, nextID)
					continue
				}
				child := &queuedNode{
					node:           nextNode,
					incomingEdgeID: edgeID,
					incomingPort:   port,
					parent:         current,
				}
				queue = append(queue, child)
			}
		}
	}

	return runErr
}

func (rt *automationRuntime) next(nodeID, port string) (string, string) {
	if nodeEdges, ok := rt.edges[nodeID]; ok {
		if ref, ok := nodeEdges[port]; ok {
			return ref.toNodeID, ref.id
		}
	}
	return "", ""
}

func executeNode(ctx context.Context, node models.APINode, payload map[string]interface{}, location models.Location) map[string]map[string]interface{} {

	defaultResult := map[string]map[string]interface{}{}

	catalogNode, ok := getCatalogNode(node.Type)
	if !ok || catalogNode.ExecFunc == nil {
		return defaultResult
	}

	result := catalogNode.ExecFunc(ctx, payload, location)
	return result
}

func defaultSuccessForKind(kind models.NodeKind) string {
	if kind == models.KindTrigger {
		return defaultPortOut
	}
	return defaultPortSuccess
}

func clonePayload(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	cp := make(map[string]interface{}, len(src))
	for key, value := range src {
		cp[key] = value
	}
	return cp
}

var placeholderPattern = regexp.MustCompile(`\{\{\s*([^{}]+?)\s*\}\}`)
var singlePlaceholderPattern = regexp.MustCompile(`^\{\{\s*([^{}]+?)\s*\}\}$`)

func substNodeFields(current *queuedNode, nodeConfig map[string]interface{}, nodePayloads map[string]map[string]map[string]interface{}) map[string]interface{} {
	if len(nodeConfig) == 0 {
		return map[string]interface{}{}
	}

	result := make(map[string]interface{}, len(nodeConfig))
	for key, value := range nodeConfig {
		result[key] = substConfigValue(current, nodePayloads, value)
	}

	return result
}

func substConfigValue(current *queuedNode, nodePayloads map[string]map[string]map[string]interface{}, value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return substNodeFields(current, v, nodePayloads)
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = substConfigValue(current, nodePayloads, item)
		}
		return out
	case string:
		return substConfigString(current, nodePayloads, v)
	default:
		return v
	}
}

func substConfigString(current *queuedNode, nodePayloads map[string]map[string]map[string]interface{}, input string) interface{} {
	trimmed := strings.TrimSpace(input)
	if match := singlePlaceholderPattern.FindStringSubmatch(trimmed); match != nil {
		if val, ok := resolvePlaceholder(current, nodePayloads, match[1]); ok {
			return val
		}
		return input
	}

	matches := placeholderPattern.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 {
		return input
	}

	var builder strings.Builder
	last := 0
	replaced := false

	for _, match := range matches {
		start, end := match[0], match[1]
		tokenStart, tokenEnd := match[2], match[3]

		builder.WriteString(input[last:start])
		token := strings.TrimSpace(input[tokenStart:tokenEnd])
		if val, ok := resolvePlaceholder(current, nodePayloads, token); ok {
			builder.WriteString(fmt.Sprint(val))
			replaced = true
		} else {
			builder.WriteString(input[start:end])
		}

		last = end
	}

	builder.WriteString(input[last:])
	if replaced {
		return builder.String()
	}

	return input
}

func resolvePlaceholder(current *queuedNode, nodePayloads map[string]map[string]map[string]interface{}, placeholder string) (interface{}, bool) {
	placeholder = strings.TrimSpace(placeholder)
	if placeholder == "" {
		return nil, false
	}

	parts := strings.SplitN(placeholder, ".", 2)
	if len(parts) != 2 {
		return nil, false
	}

	nodeID := parts[0]
	fieldPath := parts[1]
	if nodeID == "" || fieldPath == "" {
		return nil, false
	}

	portsVals, ok := nodePayloads[nodeID]
	if !ok || len(portsVals) == 0 {
		return nil, false
	}

	pathSegments := strings.Split(fieldPath, ".")
	portCandidate := portForChild(current, nodeID)
	if portCandidate == "" {
		return nil, false
	}

	value, ok := digPayloadValue(portsVals[portCandidate], pathSegments)
	if !ok {
		return nil, false
	}

	return value, true
}

func portForChild(current *queuedNode, ancestorID string) string {
	if current == nil || ancestorID == "" {
		return ""
	}

	node := current
	for node != nil {
		if node.parent != nil && node.parent.node.ID == ancestorID {
			return node.incomingPort
		}
		node = node.parent
	}

	return ""
}

func digPayloadValue(value interface{}, path []string) (interface{}, bool) {
	current := value
	for _, segment := range path {
		if segment == "" {
			return nil, false
		}

		switch typed := current.(type) {
		case map[string]interface{}:
			next, ok := typed[segment]
			if !ok {
				return nil, false
			}
			current = next
		case map[string]string:
			next, ok := typed[segment]
			if !ok {
				return nil, false
			}
			current = next
		case []interface{}:
			index, err := strconv.Atoi(segment)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}
			current = typed[index]
		default:
			return nil, false
		}
	}

	return current, true
}
