package tests

import (
	"client-runaway-zenoti/internal/db/models"
	svc_jpmreport "client-runaway-zenoti/internal/services/svc_jpmReport"
	"testing"
	"time"
)

func TestUpdateData(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2025-09-04")
	end, _ := time.Parse("2006-01-02", "2025-09-05")
	locName := "Young Medical Spa - Lansdale"
	loc := models.Location{}
	models.DB.Where("name = ?", locName).First(&loc)

	err := svc_jpmreport.UpdateReportDataForLargePeriod(start, end, loc)
	if err != nil {
		t.Errorf("Error updating report data: %v", err)
	}
}
