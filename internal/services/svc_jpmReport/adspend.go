package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/googleads"
	"client-runaway-zenoti/packages/meta"
	"context"
	"strconv"
	"time"
)

const (
	googleAdsId = "3082004096"
	metaAdsId   = "act_552470297607717"
)

func getAdSpends(loc models.Location, startDate, endDate time.Time) (adwords float64, meta float64, err error) {
	googleadsSpend, err := googleads.GetAdSpendGoPkg(context.Background(), googleAdsId, startDate, endDate)

	metaSpend, err := getMetaAdSpend(startDate, endDate)
	if err != nil {
		return 0, 0, err
	}

	return googleadsSpend, metaSpend, nil
}

func getMetaAdSpend(startDate, endDate time.Time) (float64, error) {
	token := models.Setting{}
	err := db.DB.First(&token, "key = ?", "meta_token").Error
	if err != nil {
		return 0, err
	}

	spend, err := meta.GetMetaAdSpend(token.Value, metaAdsId, "2025-09-01", "2025-09-30")
	if err != nil {
		return 0, err
	}

	spendFloat, err := strconv.ParseFloat(spend, 64)
	if err != nil {
		return 0, err
	}

	return spendFloat, nil
}
