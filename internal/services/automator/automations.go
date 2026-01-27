package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"errors"
	"fmt"
	"regexp"
	"strings"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func GetAutomations(c *gin.Context) {

	locationId := c.Param("locationId")
	automations := []models.Automation{}

	err := db.DB.Where("location_id = ?", locationId).Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		{Column: clause.Column{Name: "state"}, Desc: false},
		{Column: clause.Column{Name: "created_at"}, Desc: true},
	}}).Find(&automations).Error
	lvn.GinErr(c, 400, err, "error while getting automations")

	c.Data(lvn.Res(200, automations, ""))
}

func CreateAutomation(c *gin.Context) {

	locationId := c.Param("locationId")
	payload := models.Automation{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	payload.LocationId = locationId
	payload.CreatorId = c.MustGet("user").(models.User).ID
	payload.UpdaterId = c.MustGet("user").(models.User).ID

	if validationErrors := validateAutomationGraph(payload); len(validationErrors) > 0 {
		c.Data(lvn.Res(400, gin.H{
			"errors": validationErrorsToStrings(validationErrors),
		}, "automation validation failed"))
		return
	}

	err = db.DB.Create(&payload).Error
	lvn.GinErr(c, 400, err, "error while creating automation")

	c.Data(lvn.Res(200, payload, ""))
}

func UpdateAutomation(c *gin.Context) {

	automationId := c.Param("automationId")
	payload := models.Automation{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	var automation models.Automation
	err = db.DB.First(&automation, "id = ?", automationId).Error
	lvn.GinErr(c, 400, err, "error while getting automation")

	payload.UpdaterId = c.MustGet("user").(models.User).ID

	if validationErrors := validateAutomationGraph(payload); len(validationErrors) > 0 {
		c.Data(lvn.Res(400, gin.H{
			"errors": validationErrorsToStrings(validationErrors),
		}, "automation validation failed"))
		return
	}
	automation.Name = payload.Name
	automation.Description = payload.Description
	automation.LocationId = payload.LocationId
	automation.State = payload.State
	automation.UpdaterId = payload.UpdaterId
	automation.Graph = payload.Graph

	err = db.DB.Model(&automation).Updates(&automation).Error

	lvn.GinErr(c, 400, err, "error while updating automation")

	c.Data(lvn.Res(200, automation, ""))
}

func DeleteAutomation(c *gin.Context) {

	automationId := c.Param("automationId")

	err := db.DB.Delete(&models.Automation{}, "id = ?", automationId).Error
	lvn.GinErr(c, 400, err, "error while deleting automation")

	c.Data(lvn.Res(200, "", ""))
}

func DuplicateAutomation(c *gin.Context) {
	automationId := c.Param("automationId")

	payload := struct {
		LocationId string `json:"locationId"`
	}{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")
	if err != nil {
		return
	}

	if payload.LocationId == "" {
		c.Data(lvn.Res(400, "", "locationId is required"))
		return
	}

	var automation models.Automation
	err = db.DB.Preload("Nodes").Preload("Edges").First(&automation, "id = ?", automationId).Error
	lvn.GinErr(c, 400, err, "error while getting automation")
	if err != nil {
		return
	}

	nodeIDMap := make(map[string]string, len(automation.Graph.Nodes))
	for _, node := range automation.Graph.Nodes {
		nodeIDMap[node.ID] = uuid.New().String()
	}

	newNodes := make([]models.APINode, 0, len(automation.Graph.Nodes))
	for _, node := range automation.Graph.Nodes {
		newNode := node
		newNode.ID = nodeIDMap[node.ID]
		newNode.Config = cloneAndReplaceConfig(node.Config, nodeIDMap)
		newNodes = append(newNodes, newNode)
	}

	newEdges := make([]models.APIEdge, 0, len(automation.Graph.Edges))
	for _, edge := range automation.Graph.Edges {
		newEdge := edge
		newEdge.ID = uuid.New().String()
		if mapped, ok := nodeIDMap[edge.FromNodeId]; ok {
			newEdge.FromNodeId = mapped
		}
		if mapped, ok := nodeIDMap[edge.ToNodeId]; ok {
			newEdge.ToNodeId = mapped
		}
		newEdges = append(newEdges, newEdge)
	}

	newEntry := make([]string, 0, len(automation.Graph.Entry))
	for _, entry := range automation.Graph.Entry {
		if mapped, ok := nodeIDMap[entry]; ok {
			newEntry = append(newEntry, mapped)
			continue
		}
		newEntry = append(newEntry, entry)
	}

	user := c.MustGet("user").(models.User)
	newAutomation := models.Automation{
		ID:          uuid.New().String(),
		Name:        automation.Name,
		Description: automation.Description,
		LocationId:  payload.LocationId,
		State:       automation.State,
		CreatorId:   user.ID,
		UpdaterId:   user.ID,
		Graph: models.Graph{
			Nodes: newNodes,
			Edges: newEdges,
			Entry: newEntry,
		},
		Notes: automation.Notes,
	}

	err = db.DB.Create(&newAutomation).Error
	lvn.GinErr(c, 400, err, "error while duplicating automation")
	if err != nil {
		return
	}

	c.Data(lvn.Res(200, newAutomation, ""))

}

func validateAutomationGraph(auto models.Automation) []error {
	graph := auto.Graph
	var errs []error

	if len(graph.Nodes) == 0 {
		return []error{errors.New("automation graph must contain at least one node")}
	}

	nodeByID := make(map[string]models.APINode, len(graph.Nodes))
	catalogByNodeID := make(map[string]Node, len(graph.Nodes))

	for idx, node := range graph.Nodes {
		nodeLbl := nodeLabel(node)
		if node.ID == "" {
			errs = append(errs, fmt.Errorf("node at position %d must include an id", idx))
			continue
		}
		if _, exists := nodeByID[node.ID]; exists {
			if node.Name != "" {
				errs = append(errs, fmt.Errorf("duplicate node id %s (%s)", node.ID, node.Name))
			} else {
				errs = append(errs, fmt.Errorf("duplicate node id %s", node.ID))
			}
			continue
		}
		nodeByID[node.ID] = node

		if node.Type == "" {
			errs = append(errs, fmt.Errorf("node %s is missing node type", nodeLbl))
			continue
		}

		catalogNode, ok := getCatalogNode(node.Type)
		if !ok {
			errs = append(errs, fmt.Errorf("node %s references unknown catalog node %s", nodeLbl, node.Type))
			continue
		}
		catalogByNodeID[node.ID] = catalogNode

		if catalogNode.Type != "" && string(node.Kind) != string(catalogNode.Type) {
			errs = append(errs, fmt.Errorf("node %s kind %s does not match catalog kind %s", nodeLbl, node.Kind, catalogNode.Type))
		}
	}

	edgeByID := make(map[string]models.APIEdge, len(graph.Edges))
	edgesFrom := make(map[string][]models.APIEdge)
	for idx, edge := range graph.Edges {
		if edge.ID == "" {
			errs = append(errs, fmt.Errorf("edge at position %d must include an id", idx))
			continue
		}
		if _, exists := edgeByID[edge.ID]; exists {
			errs = append(errs, fmt.Errorf("duplicate edge id %s", edge.ID))
			continue
		}

		if edge.FromNodeId == "" {
			errs = append(errs, fmt.Errorf("edge %s missing fromNodeId", edge.ID))
		} else if _, ok := nodeByID[edge.FromNodeId]; !ok {
			errs = append(errs, fmt.Errorf("edge %s references unknown from node %s", edge.ID, edge.FromNodeId))
		}
		if edge.ToNodeId == "" {
			errs = append(errs, fmt.Errorf("edge %s missing toNodeId", edge.ID))
		} else if _, ok := nodeByID[edge.ToNodeId]; !ok {
			errs = append(errs, fmt.Errorf("edge %s references unknown to node %s", edge.ID, edge.ToNodeId))
		}

		edgeByID[edge.ID] = edge
		if edge.FromNodeId != "" {
			edgesFrom[edge.FromNodeId] = append(edgesFrom[edge.FromNodeId], edge)
		}
	}

	for idx, entryID := range graph.Entry {
		if entryID == "" {
			errs = append(errs, fmt.Errorf("entry reference at position %d is empty", idx))
			continue
		}
		node, ok := nodeByID[entryID]
		if !ok {
			errs = append(errs, fmt.Errorf("entry node %s does not exist in the node list", entryID))
			continue
		}
		if node.Kind != models.KindTrigger && node.Kind != models.KindCollection {
			errs = append(errs, fmt.Errorf("entry node %s must be a trigger or collection node", nodeLabel(node)))
		}
	}

	for _, node := range graph.Nodes {
		catalogNode, ok := catalogByNodeID[node.ID]
		if !ok {
			continue
		}
		errs = append(errs, validateNodeConfig(node, catalogNode, edgeByID, edgesFrom, nodeByID, catalogByNodeID)...)
	}

	return errs
}

func validateNodeConfig(node models.APINode, catalogNode Node, edgeByID map[string]models.APIEdge, edgesFrom map[string][]models.APIEdge, nodeByID map[string]models.APINode, catalogByNodeID map[string]Node) []error {
	var errs []error
	nodeLbl := nodeLabel(node)

	requiredFields := requiredFieldKeys(catalogNode.Fields)
	if len(node.Config) == 0 {
		if len(requiredFields) > 0 {
			errs = append(errs, fmt.Errorf("node %s is missing configuration for required fields: %s", nodeLbl, strings.Join(requiredFields, ", ")))
		}
		return errs
	}

	for cfgKey, cfg := range node.Config {
		if cfgKey == "default" {
			continue
		}
		cfgMap := cfg
		if cfgMap == nil {
			cfgMap = map[string]interface{}{}
		}

		if cfgKey != "" && cfgKey != models.DefaultNodeConfigEdge {
			if edge, ok := edgeByID[cfgKey]; !ok {
				errs = append(errs, fmt.Errorf("node %s configuration references unknown edge %s", nodeLbl, cfgKey))
			} else if edge.ToNodeId != node.ID {
				errs = append(errs, fmt.Errorf("node %s configuration references edge %s that does not point to the node", nodeLbl, cfgKey))
			}
		}

		if len(requiredFields) > 0 {
			missing := missingRequiredFields(cfgMap, requiredFields)
			if len(missing) > 0 {
				errs = append(errs, fmt.Errorf("node %s config %s missing required fields: %s", nodeLbl, configLabel(cfgKey), strings.Join(missing, ", ")))
			}
		}

		errs = append(errs, validateConfigValueReferences(node.ID, cfgMap, nodeByID, edgesFrom, catalogByNodeID)...)
	}

	return errs
}

func requiredFieldKeys(fields []NodeField) []string {
	res := []string{}
	for _, field := range fields {
		if field.Required {
			res = append(res, field.Key)
		}
	}
	return res
}

func missingRequiredFields(values map[string]interface{}, required []string) []string {
	if len(required) == 0 {
		return nil
	}
	missing := []string{}
	for _, key := range required {
		val, ok := values[key]
		if !ok || isEmptyConfigValue(val) {
			missing = append(missing, key)
		}
	}
	return missing
}

func isEmptyConfigValue(val interface{}) bool {
	if val == nil {
		return true
	}
	if str, ok := val.(string); ok {
		return strings.TrimSpace(str) == ""
	}
	return false
}

func configLabel(edgeID string) string {
	if edgeID == "" || edgeID == models.DefaultNodeConfigEdge {
		return "default"
	}
	return fmt.Sprintf("edge %s", edgeID)
}

func validateConfigValueReferences(nodeID string, value interface{}, nodeByID map[string]models.APINode, edgesFrom map[string][]models.APIEdge, catalogByNodeID map[string]Node) []error {
	var errs []error
	nodeLbl := nodeLabelByID(nodeID, nodeByID)

	switch v := value.(type) {
	case map[string]interface{}:
		for _, val := range v {
			errs = append(errs, validateConfigValueReferences(nodeID, val, nodeByID, edgesFrom, catalogByNodeID)...)
		}
	case []interface{}:
		for _, val := range v {
			errs = append(errs, validateConfigValueReferences(nodeID, val, nodeByID, edgesFrom, catalogByNodeID)...)
		}
	case string:
		matches := nodeReferenceRegex.FindAllStringSubmatch(v, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			refID := match[1]
			fieldExpr := ""
			if len(match) > 2 {
				fieldExpr = strings.TrimSpace(match[2])
			}
			fieldExpr = strings.TrimPrefix(fieldExpr, ".")
			if fieldExpr == "" {
				errs = append(errs, fmt.Errorf("node %s has invalid placeholder referencing %s in value %q", nodeLbl, refID, v))
				continue
			}
			fieldPath := strings.Split(fieldExpr, ".")
			errs = append(errs, validatePlaceholderReference(nodeID, refID, fieldPath, v, nodeByID, edgesFrom, catalogByNodeID)...)
		}
	}

	return errs
}

func validatePlaceholderReference(currentNodeID, refNodeID string, fieldPath []string, rawValue string, nodeByID map[string]models.APINode, edgesFrom map[string][]models.APIEdge, catalogByNodeID map[string]Node) []error {
	var errs []error
	targetLbl := nodeLabelByID(currentNodeID, nodeByID)
	refLbl := nodeLabelByID(refNodeID, nodeByID)

	if refNodeID == "" {
		errs = append(errs, fmt.Errorf("node %s references an empty node id in value %q", targetLbl, rawValue))
		return errs
	}

	if _, ok := nodeByID[refNodeID]; !ok {
		errs = append(errs, fmt.Errorf("node %s references unknown node %s in config value %q", targetLbl, refNodeID, rawValue))
		return errs
	}

	if currentNodeID == refNodeID {
		errs = append(errs, fmt.Errorf("node %s cannot reference itself in config value %q", targetLbl, rawValue))
		return errs
	}

	if len(fieldPath) == 0 || fieldPath[0] == "" {
		errs = append(errs, fmt.Errorf("node %s has invalid placeholder referencing %s in value %q", targetLbl, refLbl, rawValue))
		return errs
	}

	ports := findPortsForReference(refNodeID, currentNodeID, edgesFrom, nodeByID)
	if len(ports) == 0 {
		errs = append(errs, fmt.Errorf("node %s is not connected to %s, cannot reference it in value %q", targetLbl, refLbl, rawValue))
		return errs
	}

	catalogNode, ok := catalogByNodeID[refNodeID]
	if !ok {
		return errs
	}

	fieldKey := fieldPath[0]
	for _, port := range ports {
		if nodePortHasField(catalogNode, port, fieldKey) {
			return errs
		}
	}

	errs = append(errs, fmt.Errorf("node %s expects field %s from %s via port(s) %s, but it is not provided", targetLbl, fieldKey, refLbl, strings.Join(ports, ", ")))
	return errs
}

func findPortsForReference(fromNodeID, toNodeID string, edgesFrom map[string][]models.APIEdge, nodeByID map[string]models.APINode) []string {
	if fromNodeID == "" || toNodeID == "" || fromNodeID == toNodeID {
		return nil
	}

	type state struct {
		nodeID    string
		firstPort string
	}

	queue := []state{}
	visited := make(map[string]map[string]bool)
	ports := []string{}

	markVisited := func(nodeID, port string) bool {
		if nodeID == "" || port == "" {
			return false
		}
		if visited[nodeID] == nil {
			visited[nodeID] = make(map[string]bool)
		}
		if visited[nodeID][port] {
			return false
		}
		visited[nodeID][port] = true
		return true
	}

	appendPort := func(port string) {
		for _, existing := range ports {
			if existing == port {
				return
			}
		}
		ports = append(ports, port)
	}

	startEdges := edgesFrom[fromNodeID]
	for _, edge := range startEdges {
		portName := edge.FromPort
		if portName == "" {
			if fromNode, ok := nodeByID[fromNodeID]; ok {
				portName = defaultSuccessForKind(fromNode.Kind)
			}
		}
		if portName == "" {
			continue
		}
		next := state{
			nodeID:    edge.ToNodeId,
			firstPort: portName,
		}
		if next.nodeID == "" || !markVisited(next.nodeID, next.firstPort) {
			continue
		}
		if next.nodeID == toNodeID {
			appendPort(next.firstPort)
		} else {
			queue = append(queue, next)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.nodeID == toNodeID {
			appendPort(cur.firstPort)
			continue
		}

		for _, edge := range edgesFrom[cur.nodeID] {
			next := state{
				nodeID:    edge.ToNodeId,
				firstPort: cur.firstPort,
			}
			if next.nodeID == "" || !markVisited(next.nodeID, next.firstPort) {
				continue
			}
			if next.nodeID == toNodeID {
				appendPort(next.firstPort)
			} else {
				queue = append(queue, next)
			}
		}
	}

	return ports
}

func nodePortHasField(node Node, portName, fieldKey string) bool {
	if portName == "" || fieldKey == "" {
		return false
	}
	for _, port := range node.Ports {
		if port.Name != portName {
			continue
		}
		for _, field := range port.Payload {
			if field.Key == fieldKey {
				return true
			}
		}
	}
	return false
}

func nodeLabel(node models.APINode) string {
	switch {
	case node.Name != "" && node.ID != "":
		return fmt.Sprintf("%s (%s)", node.Name, node.ID)
	case node.Name != "":
		return node.Name
	default:
		return node.ID
	}
}

func nodeLabelByID(nodeID string, nodeByID map[string]models.APINode) string {
	if node, ok := nodeByID[nodeID]; ok {
		return nodeLabel(node)
	}
	return nodeID
}

func validationErrorsToStrings(errs []error) []string {
	res := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		res = append(res, err.Error())
	}
	return res
}

var nodeReferenceRegex = regexp.MustCompile(`\{\{\s*([0-9a-fA-F-]{36})([^}]*)\}\}`)

func cloneAndReplaceConfig(cfg models.NodeConfig, replacements map[string]string) models.NodeConfig {
	if cfg == nil {
		return nil
	}
	res := make(models.NodeConfig, len(cfg))
	for edgeID, edgeCfg := range cfg {
		if edgeCfg == nil {
			res[edgeID] = nil
			continue
		}
		cloned := make(map[string]interface{}, len(edgeCfg))
		for key, val := range edgeCfg {
			cloned[key] = replaceNodeReferences(val, replacements)
		}
		res[edgeID] = cloned
	}
	return res
}

func replaceNodeReferences(value interface{}, replacements map[string]string) interface{} {
	switch v := value.(type) {
	case string:
		return replaceNodeIDsInString(v, replacements)
	case map[string]interface{}:
		cloned := make(map[string]interface{}, len(v))
		for key, val := range v {
			cloned[key] = replaceNodeReferences(val, replacements)
		}
		return cloned
	case []interface{}:
		cloned := make([]interface{}, 0, len(v))
		for _, item := range v {
			cloned = append(cloned, replaceNodeReferences(item, replacements))
		}
		return cloned
	case map[string]string:
		cloned := make(map[string]string, len(v))
		for key, val := range v {
			cloned[key] = replaceNodeIDsInString(val, replacements)
		}
		return cloned
	case []string:
		cloned := make([]string, 0, len(v))
		for _, item := range v {
			cloned = append(cloned, replaceNodeIDsInString(item, replacements))
		}
		return cloned
	default:
		return value
	}
}

func replaceNodeIDsInString(value string, replacements map[string]string) string {
	if value == "" || len(replacements) == 0 {
		return value
	}
	replaced := nodeReferenceRegex.ReplaceAllStringFunc(value, func(match string) string {
		groups := nodeReferenceRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		if newID, ok := replacements[groups[1]]; ok {
			return strings.Replace(match, groups[1], newID, 1)
		}
		return match
	})
	for oldID, newID := range replacements {
		if replaced == oldID {
			return newID
		}
	}
	return replaced
}
