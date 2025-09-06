package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/packages/googleads"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"context"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func getLeadsApptsCollections(loc models.Location, startDate, endDate time.Time) (leads []runwayv2.Opportunity, appointments []zenotiv1.Appointment, collections []zenotiv1.Collection, err error) {

	// Getting leads
	svc := runway.GetSvc()
	rcli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	leads, err = rcli.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  startDate,
		EndDate:    endDate,
		PipelineId: loc.PipelineId,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Getting appointments
	zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		return nil, nil, nil, err
	}

	appointments, err = zcli.AppointmentsListAllAppointments(zenotiv1.AppointmentFilter{
		StartDate: startDate,
		EndDate:   endDate,

		IncludeNoShowCancel: true,
	})

	if err != nil {
		return nil, nil, nil, err
	}

	// Getting collections
	collections, err = zcli.ReportsAllCollections(startDate, endDate)
	if err != nil {
		return nil, nil, nil, err
	}

	return leads, appointments, collections, nil
}

func getAdSpends(loc models.Location, startDate, endDate time.Time) (adwords float64, meta float64, err error) {
	googleadsSpend, err := googleads.GetAdSpendGoPkg(context.Background(), googleAdsId, startDate, endDate)

	return googleadsSpend, 0, err
}

func calculateTotals(data []ReportData) []ReportData {
	total := ReportData{}
	total.SourceBreakdown = make(map[string]SourceBreakdown)
	total.Label = "Total"

	for _, d := range data {

		// Aggregate core metrics
		total.CoreMetrics.Leads += d.CoreMetrics.Leads
		total.CoreMetrics.Booked += d.CoreMetrics.Booked
		total.CoreMetrics.Sales += d.CoreMetrics.Sales
		total.CoreMetrics.Revenue += d.CoreMetrics.Revenue

		// Aggregate graph data
		for _, ds := range d.GraphData.Values {
			for i, tds := range total.GraphData.Values {
				if ds.Label != tds.Label {
					continue
				}
				for j, v := range ds.Data {
					total.GraphData.Values[i].Data[j] += v
				}
			}
		}

		// Aggregate source breakdown
		for k, v := range d.SourceBreakdown {

			sb, ok := total.SourceBreakdown[k]
			if !ok {
				sb = SourceBreakdown{}
			}

			sb.Leads += v.Leads
			sb.Consultations += v.Consultations
			sb.Revenue += v.Revenue
			total.SourceBreakdown[k] = sb
		}
	}

	data = append([]ReportData{total}, data...)

	return data
}

func distributeAdSpends(data []ReportData) []ReportData {

	metaCostPerLead := float64(0)
	adwordsCostPerLead := float64(0)

	for i, d := range data {
		if d.Label == "Total" {

			sb, ok := d.SourceBreakdown["Paid Search"]
			if !ok {
				sb = SourceBreakdown{}
			}
			if sb.Leads > 0 {
				adwordsCostPerLead = adwordsAdSpend / float64(sb.Leads)
			}
			sb = fillBreakdown(sb, adwordsAdSpend)
			data[i].SourceBreakdown["Paid Search"] = sb

			sb, ok = d.SourceBreakdown["Facebook paid"]
			if !ok {
				sb = SourceBreakdown{}
			}
			sb = fillBreakdown(sb, metaAdSpend)
			if sb.Leads > 0 {
				metaCostPerLead = metaAdSpend / float64(sb.Leads)
			}
			data[i].SourceBreakdown["Facebook paid"] = sb
		}
	}

	for i, d := range data {
		if d.Label == "Total" {
			continue
		}

		sb, ok := d.SourceBreakdown["Paid Search"]
		if !ok {
			sb = SourceBreakdown{}
		}
		sb = fillBreakdown(sb, adwordsCostPerLead*float64(sb.Leads))
		data[i].SourceBreakdown["Paid Search"] = sb

		sb, ok = d.SourceBreakdown["Facebook paid"]
		if !ok {
			sb = SourceBreakdown{}
		}
		sb = fillBreakdown(sb, metaCostPerLead*float64(sb.Leads))
		data[i].SourceBreakdown["Facebook paid"] = sb

	}

	return data
}

func fillBreakdown(data SourceBreakdown, adspend float64) SourceBreakdown {
	data.AdSpend = adspend
	data.CostPerLead = lvn.Ternary(data.Leads == 0, float64(0), data.AdSpend/float64(data.Leads))
	data.CostPerConsultation = lvn.Ternary(data.Consultations == 0, float64(0), data.AdSpend/float64(data.Consultations))
	data.Roas = lvn.Ternary(data.AdSpend == 0, float64(0), 100*data.Revenue/data.AdSpend)
	return data
}
