package tests

import (
	"client-runaway-zenoti/packages/googleads"
	"context"
	"testing"
	"time"
)

func TestGetAdSpend(t *testing.T) {
	ctx := context.Background()
	customerID := "3082004096"
	startDate, _ := time.Parse("2006-01-02", "2025-06-01")
	endDate, _ := time.Parse("2006-01-02", "2025-06-30")

	spend, err := googleads.GetAdSpendGoPkg(ctx, customerID, startDate, endDate)
	if err != nil {
		t.Fatalf("GetAdSpendGoPkg failed: %v", err)
	}
	t.Logf("Ad Spend: %f", spend)
}
