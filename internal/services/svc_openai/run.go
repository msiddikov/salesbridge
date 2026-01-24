package svc_openai

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	openaiv1 "client-runaway-zenoti/packages/openaiV1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RunThroughAssistant(id, text string) (string, error) {
	assistantID, err := parseAssistantID(id)
	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("text is required")
	}

	var assistant models.OpenAIAssistant
	if err := db.DB.First(&assistant, "id = ?", assistantID).Error; err != nil {
		return "", fmt.Errorf("assistant not found")
	}

	client, err := openAIClient()
	if err != nil {
		return "", err
	}

	thread, err := client.ThreadsCreate(openaiv1.ThreadRequest{})
	if err != nil {
		return "", err
	}

	_, err = client.ThreadsMessagesCreate(thread.ID, openaiv1.ThreadMessageRequest{
		Role:    "user",
		Content: text,
	})
	if err != nil {
		return "", err
	}

	run, err := client.RunsCreate(thread.ID, openaiv1.RunRequest{
		AssistantID: assistant.AssistantID,
	})
	if err != nil {
		return "", err
	}

	run, err = waitForRunCompletion(client, thread.ID, run.ID)
	if err != nil {
		return "", err
	}

	modelName := run.Model
	if modelName == "" {
		modelName = assistant.GptModel
	}
	if modelName != "" {
		err = recordAIUsage(assistant.ProfileID, modelName, run.Usage)
		if err != nil {
			return "", err
		}
	}

	messages, err := client.ThreadsMessagesList(thread.ID, 20)
	if err != nil {
		return "", err
	}

	responseText, err := firstAssistantMessage(messages)
	if err != nil {
		return "", err
	}

	return responseText, nil
}

func parseAssistantID(raw string) (uint, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, fmt.Errorf("assistant_id is required")
	}
	parsed, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("assistant_id is invalid")
	}
	return uint(parsed), nil
}

func waitForRunCompletion(client openaiv1.Client, threadID, runID string) (openaiv1.Run, error) {
	run := openaiv1.Run{}
	for i := 0; i < 30; i++ {
		res, err := client.RunsGet(threadID, runID)
		if err != nil {
			return run, err
		}
		run = res
		switch run.Status {
		case "completed":
			return run, nil
		case "failed", "cancelled", "expired", "requires_action":
			return run, fmt.Errorf("run %s", run.Status)
		}
		time.Sleep(2 * time.Second)
	}

	return run, fmt.Errorf("run timed out")
}

func firstAssistantMessage(messages openaiv1.ThreadMessagesListResponse) (string, error) {
	for _, msg := range messages.Data {
		if msg.Role != "assistant" {
			continue
		}
		for _, item := range msg.Content {
			if item.Type == "text" && item.Text != nil {
				return item.Text.Value, nil
			}
		}
	}
	return "", fmt.Errorf("assistant message not found")
}

func recordAIUsage(profileID uint, model string, usage openaiv1.ResponseUsage) error {
	if profileID == 0 || model == "" {
		return nil
	}

	usageDate := time.Now().UTC()
	usageDate = time.Date(usageDate.Year(), usageDate.Month(), usageDate.Day(), 0, 0, 0, 0, time.UTC)

	pricing := models.OpenAIModelPricing{}
	pricingErr := db.DB.First(&pricing, "model = ?", model).Error

	inputCost := 0.0
	outputCost := 0.0
	if pricingErr == nil {
		inputCost = (float64(usage.InputTokens) + float64(usage.PromptTokens)/1000.0) * pricing.InputCentsPer1K
		outputCost = (float64(usage.OutputTokens) + float64(usage.CompletionTokens)/1000.0) * pricing.OutputCentsPer1K
	}

	costCents := inputCost + outputCost

	record := models.OpenAIUsageDaily{
		ProfileID:    profileID,
		Model:        model,
		UsageDate:    usageDate,
		InputTokens:  usage.InputTokens + usage.PromptTokens,
		OutputTokens: usage.OutputTokens + usage.CompletionTokens,
		CostCents:    costCents,
		Points:       costCents,
	}

	tbl := "open_ai_usage_dailies"

	return db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "profile_id"}, {Name: "model"}, {Name: "usage_date"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"input_tokens":  gorm.Expr(tbl + ".input_tokens + EXCLUDED.input_tokens"),
			"output_tokens": gorm.Expr(tbl + ".output_tokens + EXCLUDED.output_tokens"),
			"cost_cents":    gorm.Expr(tbl + ".cost_cents + EXCLUDED.cost_cents"),
			"points":        gorm.Expr(tbl + ".points + EXCLUDED.points"),
			"updated_at":    time.Now().UTC(),
		}),
	}).Create(&record).Error
}
