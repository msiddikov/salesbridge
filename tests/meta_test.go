package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/meta"
	"context"
	"testing"
)

func TestRenewToken(t *testing.T) {
	short := "EAAK2EsTlzmEBP4ttxW7VOHcTBw7nkQi07oQaZADLBAITFbauRCj0hpdEpLGGUlmQrXOpJBsGZB3iKeTvGKh6arZCXPwyEZB7cZCoDF3bx8Uo1g6g7kfC2KE9tC5mzR2A4Qw8bDOrvsUNngwnqy8BScSgOdtXkDx6qCpAJY1nrLyksV9u9OZCWa2jML2ZCYFt6PuagZDZD"
	long, err := meta.ExchangeForLongLivedToken(context.Background(), "763141682482785", "4a2a73aa201509af22b12621b3d8741f", short)
	if err != nil {
		t.Fatalf("failed to exchange token: %v", err)
	}
	settings := models.Setting{}
	err = db.DB.Where("key = ?", "meta_token").First(&settings).Error
	if err != nil {
		t.Fatalf("failed to get setting: %v", err)
	}

	settings.Value = long.AccessToken
	err = db.DB.Save(&settings).Error
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
