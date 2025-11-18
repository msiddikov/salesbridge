package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var (

	// Others category
	othersCategory = Category{
		Id:   "others",
		Name: "Others",
		Nodes: []Node{
			othersActionDelay,
			othersActionCondition,
		},
	}

	// Actions
	othersActionDelay = Node{
		Id:          "others.delay",
		Title:       "Delay",
		Description: "Delays the workflow for a specified amount of time.",
		Type:        NodeTypeAction,
		Icon:        "ri:timer-flash-line",
		Color:       ColorDefault,
		Ports: []NodePort{
			{Name: "done"},
		},
		Fields: []NodeField{
			{Key: "duration", Type: "number", Required: true},
			{Key: "unit", Type: "string", Required: true},
		},
	}

	othersActionCondition = Node{
		Id:          "others.condition",
		Title:       "Condition",
		Description: "Branches the workflow based on a condition.",
		ExecFunc:    othersCondition,
		Type:        NodeTypeAction,
		Icon:        "ri:divide-line",
		Color:       ColorDefault,
		Ports: []NodePort{
			customPort("true", othersConditionNodeFields),
			customPort("false", othersConditionNodeFields),
		},
		Fields: []NodeField{
			{Key: "left", Type: "string", Required: true},
			{Key: "operator", Type: "string", Required: true, SelectOptions: []string{"equals", "not_equals", "greater_than", "less_than", "contains"}},
			{Key: "right", Type: "string", Required: true},
		},
	}

	//////////////////////////////////////////////////
	//                  Node Fields
	///////////////////////////////////////////////////
	othersConditionNodeFields = []NodeField{
		{Key: "left", Label: "Left Value", Type: "string"},
		{Key: "operator", Label: "Operator", Type: "string"},
		{Key: "right", Label: "Right Value", Type: "string"},
		{Key: "result", Label: "Result", Type: "bool"},
	}
)

//////////////////////////////////////////////////
//
//                  Functions
//
///////////////////////////////////////////////////

func othersCondition(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	left := fields["left"]
	right := fields["right"]
	operator, _ := fields["operator"].(string)

	result := evaluateCondition(left, right, operator)
	port := "false"
	if result {
		port = "true"
	}

	payloadData := map[string]interface{}{
		"left":     left,
		"right":    right,
		"operator": operator,
		"result":   result,
	}

	return customPayload(port, payloadData)
}

func evaluateCondition(left, right interface{}, operator string) bool {
	switch operator {
	case "equals":
		return fmt.Sprint(left) == fmt.Sprint(right)
	case "not_equals":
		return fmt.Sprint(left) != fmt.Sprint(right)
	case "greater_than":
		lv, lok := toFloat(left)
		rv, rok := toFloat(right)
		return lok && rok && lv > rv
	case "less_than":
		lv, lok := toFloat(left)
		rv, rok := toFloat(right)
		return lok && rok && lv < rv
	case "contains":
		return strings.Contains(strings.ToLower(fmt.Sprint(left)), strings.ToLower(fmt.Sprint(right)))
	default:
		return false
	}
}

func toFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return f, err == nil
	default:
		return 0, false
	}
}
