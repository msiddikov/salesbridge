package svc_jpmreport

import (
	"time"
)

func buildGraphData(startDate, endDate time.Time, leadsData []LeadRawData) GraphData {
	dates := []string{}
	totalLeadsData := []float64{}

	// Create a map to hold counts for each date
	leadsCountMap := make(map[string]int)

	// Populate the maps with counts
	for _, lead := range leadsData {
		dateStr := lead.CreatedDate.Format("2006-01-02")
		leadsCountMap[dateStr]++
	}

	// Iterate over the date range and build the data arrays
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dates = append(dates, dateStr)

		totalLeadsData = append(totalLeadsData, float64(leadsCountMap[dateStr]))
	}

	return GraphData{
		Dates: dates,
		Values: []GraphDataset{
			{
				Label: "Total leads",
				Data:  totalLeadsData,
			},
		},
	}
}

func buildGraphDataOld(startDate, endDate time.Time, leadsData []LeadRawData) GraphData {
	dates := []string{}
	leadToConsultationData := []float64{}
	consultToProcedureData := []float64{}

	// Create a map to hold counts for each date
	leadsCountMap := make(map[string]int)
	consultationsCountMap := make(map[string]int)
	proceduresCountMap := make(map[string]int)

	// Populate the maps with counts
	for _, lead := range leadsData {
		dateStr := lead.CreatedDate.Format("2006-01-02")
		leadsCountMap[dateStr]++
		for _, appt := range lead.Appointments {
			if appt.Total > 0 {
				proceduresCountMap[dateStr]++
			} else {
				consultationsCountMap[dateStr]++
			}
		}
	}

	// Iterate over the date range and build the data arrays
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dates = append(dates, dateStr)
		leadToConsult := float64(0)
		consultToProcedure := float64(0)

		leads, leadsOk := leadsCountMap[dateStr]
		consultations, consultationsOk := consultationsCountMap[dateStr]
		procedures, proceduresOk := proceduresCountMap[dateStr]

		if leadsOk && consultationsOk && leads != 0 {
			leadToConsult = 100 * (float64(consultations) / float64(leads))
		}

		if consultationsOk && proceduresOk && consultations != 0 {
			consultToProcedure = 100 * (float64(procedures) / float64(consultations))
		}

		leadToConsultationData = append(leadToConsultationData, leadToConsult)
		consultToProcedureData = append(consultToProcedureData, consultToProcedure)
	}

	return GraphData{
		Dates: dates,
		Values: []GraphDataset{
			{
				Label: "Lead to Consultation",
				Data:  leadToConsultationData,
			},
			{
				Label: "Consultation to Procedure",
				Data:  consultToProcedureData,
			},
		},
	}
}
