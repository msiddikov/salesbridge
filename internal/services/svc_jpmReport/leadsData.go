package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"time"
)

func getLeadsWithAppointmentsDb(loc models.Location, startDate, endDate time.Time) (leadsData []LeadRawData, err error) {

	newLeads := []models.JpmReportNewLead{}
	err = models.DB.
		Where("location_id = ? AND date >= ? AND date <= ?",
			loc.Id,
			startDate,
			endDate).
		Preload("Invoices").
		Find(&newLeads).Error
	if err != nil {
		return nil, err
	}

	for _, nl := range newLeads {
		lead := LeadRawData{
			Name:        nl.Name,
			Source:      nl.Source,
			CreatedDate: nl.Date,
		}

		for _, inv := range nl.Invoices {
			lead.Appointments = append(lead.Appointments, LeadAppointment{
				Date:  inv.Date,
				Total: inv.Total,
				Link:  loc.GetZenotiAppointmentLink(inv.InvoiceId),
			})
			lead.Revenue += inv.Total
			lead.HadConsultation = true
			if inv.Total > 0 {
				lead.HadSale = true
			}
		}

		leadsData = append(leadsData, lead)
	}

	return leadsData, nil
}

func getLeadsWithAppointments(loc models.Location, startDate, endDate time.Time) (leadsData []LeadRawData, err error) {

	// Getting leads
	svc := runway.GetSvc()
	rcli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		return nil, err
	}

	leads, err := rcli.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  startDate,
		EndDate:    endDate,
		PipelineId: loc.PipelineId,
	})
	if err != nil {
		return nil, err
	}

	for _, opp := range leads {
		lead := fillAppointmentData(loc, opp)
		leadsData = append(leadsData, lead)
	}

	return leadsData, nil
}

func fillAppointmentData(l models.Location, opp runwayv2.Opportunity) LeadRawData {
	lead := LeadRawData{
		Name:        opp.Name,
		Source:      opp.Source,
		CreatedDate: opp.CreatedAt,
	}

	// Getting appointments
	zcli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)

	guest, err := zcli.GuestsGetByPhoneEmail(opp.Contact.Phone, opp.Contact.Email)
	if err != nil || len(guest) == 0 {
		return lead
	}

	appts, _, err := zcli.GuestsListAppointments(zenotiv1.GuestAppointmentsFilter{
		GuestId:   guest[0].Id,
		StartDate: lead.CreatedDate,
		EndDate:   time.Now(),
	})

	if err != nil || len(appts) == 0 {
		return lead
	}

	for _, a := range appts {
		if a.Invoice_status == zenotiv1.InvoiceClosed || a.Invoice_status == zenotiv1.InvoiceCampaignApplied || a.Invoice_status == zenotiv1.InvoiceCouponApplied {
			lead.Revenue += float64(a.Price.Sales)
			lead.HadConsultation = true

		}

		lead.Appointments = append(lead.Appointments, LeadAppointment{
			Date:  a.GetDate(),
			Total: float64(a.Price.Sales),
			Link:  l.GetZenotiAppointmentLink(a.Invoice_id),
		})
	}

	return lead
}

func getTotalBookedAndSales(leadsData []LeadRawData) (booked int, sales int, revenue float64) {
	for _, l := range leadsData {
		booked += len(l.Appointments)
		for _, a := range l.Appointments {
			if a.Total > 0 {
				sales++
				revenue += a.Total
			}
		}
	}
	return
}
