package svc_openai

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	openaiv1 "client-runaway-zenoti/packages/openaiV1"
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type assistantPayload struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	Instructions string `json:"instructions,omitempty"`
}

type assistantResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Model     string `json:"model"`
	CreatedAt int64  `json:"createdAt"`
}

func ListAssistants(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	var assistants []models.OpenAIAssistant

	err := db.DB.Where("profile_id = ?", user.ProfileID).Order("created_at desc").Find(&assistants).Error
	lvn.GinErr(c, 400, err, "unable to list assistants")
	if err != nil {
		return
	}

	resp := make([]assistantResponse, 0, len(assistants))
	for _, a := range assistants {
		resp = append(resp, assistantResponse{
			ID:        a.ID,
			Name:      a.Name,
			Model:     a.GptModel,
			CreatedAt: a.CreatedAt.Unix(),
		})
	}

	c.Data(lvn.Res(200, resp, "OK"))
}

func CreateAssistant(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	payload := assistantPayload{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")
	if err != nil {
		return
	}

	if payload.Name == "" || payload.Model == "" {
		lvn.GinErr(c, 400, nil, "name and model are required")
		return
	}

	client, err := openAIClient()
	lvn.GinErr(c, 400, err, "unable to init ai client")
	if err != nil {
		return
	}

	assistant, err := client.AssistantsCreate(openaiv1.AssistantRequest{
		Model:        payload.Model,
		Name:         payload.Name,
		Instructions: payload.Instructions,
	})
	lvn.GinErr(c, 400, err, "unable to create assistant")
	if err != nil {
		return
	}

	record := models.OpenAIAssistant{
		AssistantID: assistant.ID,
		ProfileID:   user.ProfileID,
		Name:        payload.Name,
		GptModel:    payload.Model,
	}

	err = db.DB.Create(&record).Error
	lvn.GinErr(c, 400, err, "unable to persist assistant")
	if err != nil {
		return
	}

	c.Data(lvn.Res(200, assistantResponse{
		ID:        record.ID,
		Name:      record.Name,
		Model:     record.GptModel,
		CreatedAt: record.CreatedAt.Unix(),
	}, "OK"))
}

func UpdateAssistant(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	id, err := strconv.ParseUint(c.Param("assistantId"), 10, 64)
	lvn.GinErr(c, 400, err, "invalid assistant id")
	if err != nil {
		return
	}

	payload := assistantPayload{}
	err = c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")
	if err != nil {
		return
	}

	if payload.Name == "" && payload.Model == "" && payload.Instructions == "" {
		lvn.GinErr(c, 400, nil, "no fields provided")
		return
	}

	var record models.OpenAIAssistant
	err = db.DB.Where("profile_id = ? AND id = ?", user.ProfileID, id).First(&record).Error
	lvn.GinErr(c, 404, err, "assistant not found")
	if err != nil {
		return
	}

	client, err := openAIClient()
	lvn.GinErr(c, 400, err, "unable to init ai client")
	if err != nil {
		return
	}

	updateReq := openaiv1.AssistantUpdateRequest{}
	if payload.Name != "" {
		updateReq.Name = &payload.Name
	}
	if payload.Model != "" {
		updateReq.Model = &payload.Model
	}
	if payload.Instructions != "" {
		updateReq.Instructions = &payload.Instructions
	}

	_, err = client.AssistantsUpdate(record.AssistantID, updateReq)
	lvn.GinErr(c, 400, err, "unable to update assistant")
	if err != nil {
		return
	}

	updates := map[string]interface{}{}
	if payload.Name != "" {
		updates["name"] = payload.Name
		record.Name = payload.Name
	}
	if payload.Model != "" {
		updates["gpt_model"] = payload.Model
		record.GptModel = payload.Model
	}
	if len(updates) > 0 {
		if err := db.DB.Model(&record).Updates(updates).Error; err != nil {
			lvn.GinErr(c, 400, err, "unable to update assistant record")
			return
		}
	}

	c.Data(lvn.Res(200, assistantResponse{
		ID:        record.ID,
		Name:      record.Name,
		Model:     record.GptModel,
		CreatedAt: record.CreatedAt.Unix(),
	}, "OK"))
}

func DeleteAssistant(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	id, err := strconv.ParseUint(c.Param("assistantId"), 10, 64)
	lvn.GinErr(c, 400, err, "invalid assistant id")
	if err != nil {
		return
	}

	var record models.OpenAIAssistant
	err = db.DB.Where("profile_id = ? AND id = ?", user.ProfileID, id).First(&record).Error
	lvn.GinErr(c, 404, err, "assistant not found")
	if err != nil {
		return
	}

	client, err := openAIClient()
	lvn.GinErr(c, 400, err, "unable to init ai client")
	if err != nil {
		return
	}

	_, err = client.AssistantsDelete(record.AssistantID)
	lvn.GinErr(c, 400, err, "unable to delete assistant")
	if err != nil {
		return
	}

	if err := db.DB.Delete(&record).Error; err != nil {
		lvn.GinErr(c, 400, err, "unable to delete assistant record")
		return
	}

	c.Data(lvn.Res(200, gin.H{"deleted": record.ID}, "OK"))
}

func ListModels(c *gin.Context) {
	models := []gin.H{
		{"id": "gpt-5-mini", "name": "GPT-5 mini"},
		{"id": "gpt-4.1-nano", "name": "GPT-4.1 nano"},
		{"id": "gpt-4.1-mini", "name": "GPT-4.1 mini"},
	}

	c.Data(lvn.Res(200, models, "OK"))
}
