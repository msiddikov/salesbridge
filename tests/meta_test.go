package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/meta"
	"context"
	"testing"
)

func TestRenewToken(t *testing.T) {
	short := "EAAK2EsTlzmEBPU4GqUPSTPTIgMVAHrJqFOWfe3ZCw6KLvhpOrPilfpvZASYWLnrVFH40WS3YZBtlJyZA3qCMLh71Iff3btjBXzActeqbe7T6s27ZCnYsg7SnQZCH7P9pTJyZBXM8QI80y1iMYxzemNL8wGAFIBrjA7XZCeZBHbjOA4E7UyXf3XG9VEpNyZCLPM4AW7H76QdKSE4kxqXhAJP3qf6kPpR6aqtk2yO7A0QVayfXvDZCAZDZD"
	long, err := meta.ExchangeForLongLivedToken(context.Background(), "763141682482785", "4a2a73aa201509af22b12621b3d8741f", short)
	if err != nil {
		t.Fatalf("failed to exchange token: %v", err)
	}

	setting := models.Setting{
		Key:   "meta_token",
		Value: long.AccessToken,
	}
	err = db.DB.Create(&setting).Error
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}
}

func TestGetMetaAdSpend(t *testing.T) {
	token := models.Setting{}
	err := db.DB.First(&token, "key = ?", "meta_token").Error
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}

	spend, err := meta.GetMetaAdSpend(token.Value, "act_552470297607717", "2025-09-01", "2025-09-30")
	if err != nil {
		t.Fatalf("failed to get ad spend: %v", err)
	}

	t.Logf("Ad spend: %s", spend)
}
