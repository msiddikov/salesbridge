package googleads

import (
	"context"
	"testing"
)

func TestGetAdSpend(t *testing.T) {
	ctx := context.Background()
	customerID := "3082004096"
	startDate := "2025-06-01"
	endDate := "2025-06-30"

	spend, err := GetAdSpendGoPkg(ctx, customerID, startDate, endDate)
	if err != nil {
		t.Fatalf("GetAdSpendGoPkg failed: %v", err)
	}
	t.Logf("Ad Spend: %f", spend)
}
