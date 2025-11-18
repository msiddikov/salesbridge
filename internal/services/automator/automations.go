package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"regexp"
	"strings"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAutomations(c *gin.Context) {

	locationId := c.Param("locationId")
	automations := []models.Automation{}

	err := db.DB.Where("location_id = ?", locationId).Find(&automations).Error
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

	automation.Graph = payload.Graph
	err = db.DB.Model(&automation).Updates(payload).Error
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
