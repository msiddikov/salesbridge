package models

import (
	"time"

	"gorm.io/gorm"
)

type OpenAIAssistant struct {
	AssistantID string `gorm:"not null;uniqueIndex"`
	ProfileID   uint   `gorm:"index"`
	Name        string `gorm:"not null"`
	GptModel    string `gorm:"not null"`
	gorm.Model
}

type OpenAIModelPricing struct {
	Model            string  `gorm:"primaryKey"`
	InputCentsPer1K  float64 `gorm:"not null"`
	OutputCentsPer1K float64 `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OpenAIUsageDaily struct {
	ID           uint      `gorm:"primaryKey"`
	ProfileID    uint      `gorm:"index:idx_openai_usage_profile_model_date,unique"`
	Model        string    `gorm:"index:idx_openai_usage_profile_model_date,unique"`
	UsageDate    time.Time `gorm:"index:idx_openai_usage_profile_model_date,unique"`
	InputTokens  int       `gorm:"not null;default:0"`
	OutputTokens int       `gorm:"not null;default:0"`
	Credits      float64   `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
