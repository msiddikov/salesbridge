package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (

	// Others category
	othersCategory = Category{
		Id:    "others",
		Name:  "Others",
		Icon:  "ri:pencil-ruler-2-line",
		Color: "#FFB020",
		Nodes: []Node{
			othersActionDelay,
			othersActionCondition,
			othersActionTransformNumber,
		},
	}

	// Actions
	othersActionDelay = Node{
		Id:          "others.delay",
		Title:       "Delay",
		Description: "Delays the workflow for a specified amount of time.",
		ExecFunc:    othersDelay,
		Type:        NodeTypeAction,
		Icon:        "ri:timer-flash-line",
		Color:       ColorDefault,
		Ports:       []NodePort{customPort("done", []NodeField{})},
		Fields: []NodeField{
			{Key: "duration", Type: "string", Required: true},
			{Key: "unit", Type: "string", Required: true, SelectOptions: []string{"seconds", "minutes"}},
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
			customPort("true", othersTransformationNodeFields),
			customPort("false", othersTransformationNodeFields),
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
	othersTransformationNodeFields = []NodeField{
		{Key: "left", Label: "Left Value", Type: "string"},
		{Key: "operator", Label: "Operator", Type: "string"},
		{Key: "right", Label: "Right Value", Type: "string"},
		{Key: "result", Label: "Result", Type: "string"},
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

func toString(value interface{}) string {
	return fmt.Sprint(value)
}

var othersActionTransformNumber = Node{
	Id:          "others.transform.number",
	Title:       "Transform Number",
	Description: "Transforms a value to a number.",
	ExecFunc:    othersTransformNumber,
	Type:        NodeTypeAction,
	Icon:        "ri:divide-line",
	Color:       ColorDefault,
	Ports: []NodePort{
		successPort(othersTransformationNodeFields),
		errorPort,
	},
	Fields: []NodeField{
		{Key: "left", Type: "string", Required: true},
		{Key: "operator", Type: "string", Required: true, SelectOptions: []string{"add", "subtract", "multiply", "divide"}},
		{Key: "right", Type: "string", Required: true},
	},
}

func othersTransformNumber(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	lv, ok := toFloat(fields["left"])
	if !ok {
		return errorPayload(fmt.Errorf("unable to parse left value"), "error")
	}

	rv, ok := toFloat(fields["right"])
	if !ok {
		return errorPayload(fmt.Errorf("unable to parse right value"), "error")
	}

	operator, _ := fields["operator"].(string)

	var result float64

	switch operator {
	case "add":
		result = lv + rv
	case "subtract":
		result = lv - rv
	case "multiply":
		result = lv * rv
	case "divide":
		if rv != 0 {
			result = lv / rv
		} else {
			return errorPayload(fmt.Errorf("division by zero"), "error")
		}
	}

	payloadData := map[string]interface{}{
		"left":     fields["left"],
		"right":    fields["right"],
		"operator": fields["operator"],
		"result":   fmt.Sprintf("%v", result),
	}

	return successPayload(payloadData)
}

func othersDelay(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	durationStr, ok := fields["duration"].(string)
	if !ok {
		return errorPayload(fmt.Errorf("invalid duration"), "error")
	}
	duration, ok := toFloat(durationStr)
	if !ok {
		return errorPayload(fmt.Errorf("invalid duration"), "error")
	}

	unit, ok := fields["unit"].(string)
	if !ok {
		return errorPayload(fmt.Errorf("invalid unit"), "error")
	}
	var timeDuration time.Duration
	switch unit {
	case "seconds":
		timeDuration = time.Duration(duration) * time.Second
	case "minutes":
		timeDuration = time.Duration(duration) * time.Minute
	default:
		return errorPayload(fmt.Errorf("invalid unit"), "error")
	}

	time.Sleep(timeDuration)

	return customPayload("done", map[string]interface{}{})
}
