package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db/models"
	"math"
	"sync"
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
		LeadRawData     []LeadRawData              `json:"lead_raw_data,omitempty"`
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

	LeadRawData struct {
		Name            string
		Source          string
		CreatedDate     time.Time
		Status          string
		Stage           string
		HadConsultation bool
		HadSale         bool
		ConsultDate     time.Time
		Revenue         float64
		Appointments    []LeadAppointment
	}

	LeadAppointment struct {
		Date   time.Time
		Status string
		Total  float64
		Link   string
	}
)

var (
	metaAdSpend    float64
	adwordsAdSpend float64
	locations      = []string{"Young Medical Spa", "Young Medical Spa - Lansdale", "Young Medical Spa - Wilkes-Barre/Scranton"}
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
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// build data
	res := []ReportData{}
	wg := sync.WaitGroup{}
	for _, locName := range locations {
		loc := models.Location{}
		models.DB.Where("name = ?", locName).First(&loc)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err = getReportForLocation(loc, startDate, endDate, &res)
		}()
	}
	wg.Wait()
	if err != nil {
		lvn.GinErr(c, 500, err, "Unable to get report for locations")
	}

	res = calculateTotals(res)

	adwordsAdSpend, metaAdSpend, err = getAdSpends(models.Location{}, startDate, endDate)
	lvn.GinErr(c, 500, err, "Unable to get ad spends")
	res = distributeAdSpends(res)

	c.JSON(200, res)
}

func getReportForLocation(loc models.Location, startDate, endDate time.Time, results *[]ReportData) (ReportData, error) {
	res := ReportData{
		Label: loc.Name,
	}
	// get leads, appointments and collections for period
	leadsData, err := getLeadsWithAppointmentsDb(loc, startDate, endDate)
	if err != nil {
		return ReportData{}, err
	}

	booked, sales, revenue := getTotalBookedAndSales(leadsData)
	res.CoreMetrics = CoreMetrics{
		Leads:   len(leadsData),
		Booked:  booked,
		Sales:   sales,
		Revenue: revenue,
	}

	if res.CoreMetrics.Leads > 0 {
		res.CoreMetrics.LeadsToBooked = math.Round((float64(res.CoreMetrics.Booked) / float64(res.CoreMetrics.Leads)) * 100)
	}

	if res.CoreMetrics.Booked > 0 {
		res.CoreMetrics.BookedToSales = math.Round((float64(res.CoreMetrics.Sales) / float64(res.CoreMetrics.Booked)) * 100)
	}

	res.GraphData = buildGraphData(startDate, endDate, leadsData)

	res.SourceBreakdown = breakBySource(leadsData)
	res.LeadRawData = leadsData

	*results = append(*results, res)
	return res, nil
}
