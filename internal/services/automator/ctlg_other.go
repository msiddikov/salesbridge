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
			othersActionTransformLogic,
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
			{Key: "right", Type: "string", Required: false},
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

var othersActionTransformLogic = Node{
	Id:          "others.transform.logic",
	Title:       "Logic Transform",
	Description: "Performs logical operations on boolean values.",
	ExecFunc:    othersTransformLogic,
	Type:        NodeTypeAction,
	Icon:        "ri:git-branch-line",
	Color:       ColorDefault,
	Ports: []NodePort{
		successPort(othersLogicTransformNodeFields),
		errorPort,
	},
	Fields: []NodeField{
		{Key: "left", Label: "Left Value", Type: "string", Required: true},
		{Key: "operator", Label: "Operator", Type: "string", Required: true, SelectOptions: []string{"and", "or", "not", "xor"}},
		{Key: "right", Label: "Right Value (not used for NOT)", Type: "string", Required: false},
	},
}

var othersLogicTransformNodeFields = []NodeField{
	{Key: "left", Label: "Left Value", Type: "string"},
	{Key: "operator", Label: "Operator", Type: "string"},
	{Key: "right", Label: "Right Value", Type: "string"},
	{Key: "result", Label: "Result", Type: "string"},
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

func othersTransformLogic(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	operator, _ := fields["operator"].(string)

	leftVal, ok := toBool(fields["left"])
	if !ok {
		return errorPayload(fmt.Errorf("unable to parse left value as boolean"), "error")
	}

	var result bool

	switch operator {
	case "not":
		result = !leftVal
	case "and":
		rightVal, ok := toBool(fields["right"])
		if !ok {
			return errorPayload(fmt.Errorf("unable to parse right value as boolean"), "error")
		}
		result = leftVal && rightVal
	case "or":
		rightVal, ok := toBool(fields["right"])
		if !ok {
			return errorPayload(fmt.Errorf("unable to parse right value as boolean"), "error")
		}
		result = leftVal || rightVal
	case "xor":
		rightVal, ok := toBool(fields["right"])
		if !ok {
			return errorPayload(fmt.Errorf("unable to parse right value as boolean"), "error")
		}
		result = leftVal != rightVal
	default:
		return errorPayload(fmt.Errorf("unknown operator: %s", operator), "error")
	}

	payloadData := map[string]interface{}{
		"left":     fields["left"],
		"right":    fields["right"],
		"operator": operator,
		"result":   fmt.Sprintf("%t", result),
	}

	return successPayload(payloadData)
}

func toBool(value interface{}) (bool, bool) {
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		lower := strings.ToLower(strings.TrimSpace(v))
		if lower == "true" || lower == "1" || lower == "yes" {
			return true, true
		}
		if lower == "false" || lower == "0" || lower == "no" || lower == "" {
			return false, true
		}
		return false, false
	case int:
		return v != 0, true
	case int64:
		return v != 0, true
	case float64:
		return v != 0, true
	default:
		return false, false
	}
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
