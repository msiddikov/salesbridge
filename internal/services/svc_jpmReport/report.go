package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db/models"
	"math"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type (
	ReportData struct {
		Label           string                     `json:"label"`
		CoreMetrics     CoreMetrics                `json:"core_metrics"`
		GraphData       GraphData                  `json:"graph_data"`
		SourceBreakdown map[string]SourceBreakdown `json:"source_breakdown"`
	}

	CoreMetrics struct {
		Leads         int     `json:"total_leads"`
		LeadsToBooked float64 `json:"leads_to_booked"`
		Booked        int     `json:"booked"`
		BookedToSales float64 `json:"booked_to_sales"`
		Sales         int     `json:"sales"`
		Revenue       float64 `json:"revenue"`
	}

	GraphData struct {
		Dates  []string       `json:"dates"`
		Values []GraphDataset `json:"values"`
	}

	GraphDataset struct {
		Label string    `json:"label"`
		Data  []float64 `json:"data"`
	}

	SourceBreakdown struct {
		Source              string  `json:"source"`
		AdSpend             float64 `json:"ad_spend"`
		Leads               int     `json:"leads"`
		CostPerLead         float64 `json:"cost_per_lead"`
		Consultations       int     `json:"consultations"`
		CostPerConsultation float64 `json:"cost_per_consultation"`
		Revenue             float64 `json:"revenue"`
		Roas                float64 `json:"roas"`
	}
)

var (
	metaAdSpend    float64
	adwordsAdSpend float64
	locations      = []string{"Young Medical Spa", "Young Medical Spa - Lansdale", "Young Medical Spa - Wilkes-Barre/Scranton"}
)

const (
	googleAdsId = "3082004096"
)

func GetReport(c *gin.Context) {
	// collect params
	startDateString := c.Query("start")
	endDateString := c.Query("end")

	if startDateString == "" && endDateString == "" {
		c.Data(lvn.Res(400, "Both start and end are required", "Bad request"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateString)
	lvn.GinErr(c, 400, err, "Invalid start")
	endDate, err := time.Parse("2006-01-02", endDateString)
	lvn.GinErr(c, 400, err, "Invalid end")

	// build data
	res := []ReportData{}

	for _, locName := range locations {
		loc := models.Location{}
		models.DB.Where("name = ?", locName).First(&loc)

		report, err := getReportForLocation(loc, startDate, endDate)
		lvn.GinErr(c, 500, err, "Unable to get report for location "+loc.Name)

		res = append(res, report)
	}

	res = calculateTotals(res)

	adwordsAdSpend, metaAdSpend, err = getAdSpends(models.Location{}, startDate, endDate)
	lvn.GinErr(c, 500, err, "Unable to get ad spends")
	res = distributeAdSpends(res)

	c.JSON(200, res)
}

func getReportForLocation(loc models.Location, startDate, endDate time.Time) (ReportData, error) {
	res := ReportData{
		Label: loc.Name,
	}
	// get leads, appointments and collections for period
	leads, appointments, collections, err := getLeadsApptsCollections(loc, startDate, endDate)
	if err != nil {
		return ReportData{}, err
	}

	res.CoreMetrics = CoreMetrics{
		Leads:   len(leads),
		Booked:  len(appointments),
		Sales:   len(collections),
		Revenue: 0,
	}

	if len(leads) > 0 {
		res.CoreMetrics.LeadsToBooked = math.Round((float64(len(appointments)) / float64(len(leads))) * 100)
	}

	if len(appointments) > 0 {
		res.CoreMetrics.BookedToSales = math.Round((float64(len(collections)) / float64(len(appointments))) * 100)
	}

	for _, c := range collections {
		res.CoreMetrics.Revenue += c.Total_collection
	}

	res.GraphData = buildGraphData(startDate, endDate, leads, collections)

	res.SourceBreakdown = breakBySource(leads, appointments, collections)

	return res, nil
}
