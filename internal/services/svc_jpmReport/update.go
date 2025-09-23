package svc_jpmreport

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/tgbot"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"gorm.io/gorm/clause"
)

func UpdateAllLocationsReportDataForLargePeriod(start, end time.Time) error {
	// get all locations
	for _, lId := range locations {
		loc := models.Location{}
		err := models.DB.Where("name = ?", lId).First(&loc).Error
		if err != nil {
			return err
		}

		err = UpdateReportDataForLargePeriod(start, end, loc)
		if err != nil {
			tgbot.Notify("jpm_report_errors", fmt.Sprintf("Error updating report data for location %s: %v", loc.Name, err), true)
			continue
		}
	}
	return nil
}

func UpdateReportDataForLargePeriod(start, end time.Time, loc models.Location) error {
	// split the period into smaller chunks of 7 days
	curStart := start
	curEnd := start.Add(6 * 24 * time.Hour)
	curEnd = lvn.Ternary(curEnd.After(end), end, curEnd)

	for curStart.Before(end) {
		err := UpdateReportDataForPeriod(curStart, curEnd, loc)
		if err != nil {
			tgbot.Notify("jpm_report_errors", fmt.Sprintf("Error updating report data from %s to %s for location %s: %v", curStart.Format("2006-01-02"), curEnd.Format("2006-01-02"), loc.Name, err), true)
			continue
		}
		curStart = curEnd
		curEnd = curStart.Add(6 * 24 * time.Hour)
		curEnd = lvn.Ternary(curEnd.After(end), end, curEnd)
		tgbot.Notify("jpm_report", fmt.Sprintf("Finished updating report data from %s to %s for location %s", curStart.Format("2006-01-02"), curEnd.Format("2006-01-02"), loc.Name), false)
	}

	return nil

}

func UpdateReportDataForPeriod(start, end time.Time, l models.Location) error {
	// retrieve new leads from runway
	err := updateNewLeads(start, end, l)
	if err != nil {
		return err
	}

	// update guestIds
	err = updateGuestIds(l)
	if err != nil {
		return err
	}

	// update new sales from zenoti
	err = updateNewSales(start, end, l)
	if err != nil {
		return err
	}

	return nil
}

func updateNewLeads(start, end time.Time, loc models.Location) error {
	// Getting leads
	svc := runway.GetSvc()
	rcli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		return err
	}

	leads, err := rcli.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  start,
		EndDate:    end,
		PipelineId: loc.PipelineId,
	})
	if err != nil {
		return err
	}

	for _, lead := range leads {
		reportLead := models.JpmReportNewLead{
			Date:          lead.CreatedAt,
			LocationId:    loc.Id,
			Name:          lead.Name,
			ContactId:     lead.Contact.Id,
			OpportunityId: lead.Id,
			Phone:         lead.Contact.Phone,
			Email:         lead.Contact.Email,
			Source:        lead.Source,
		}

		// save to db, if opportunityId conflicts update all besides GuestId
		err = models.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "opportunity_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"date", "location_id", "name", "contact_id", "phone", "email", "source", "updated_at"}),
		}).Create(&reportLead).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func updateGuestIds(loc models.Location) error {
	// get all leads without guestId
	var leads []models.JpmReportNewLead
	err := models.DB.Where("location_id = ? AND (guest_id IS NULL OR guest_id = '')", loc.Id).Find(&leads).Error
	if err != nil {
		return err
	}

	svc := runway.GetSvc()
	rcli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		return err
	}

	for k, lead := range leads {
		// update lead info first
		opp, err := rcli.OpportunitiesGet(lead.OpportunityId)
		if err != nil {
			return err
		}

		lead.Name = opp.Name
		lead.Phone = opp.Contact.Phone
		lead.Email = opp.Contact.Email
		lead.Source = opp.Source

		// get guestId from zenoti
		zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)

		guest, err := zcli.GuestsGetByPhoneEmail(opp.Contact.Phone, opp.Contact.Email)
		if err == nil && len(guest) > 0 {
			lead.GuestId = strings.ToLower(guest[0].Id)
		}
		leads[k] = lead
	}

	// save all leads
	if len(leads) == 0 {
		return nil
	}

	err = models.DB.Save(&leads).Error
	if err != nil {
		return err
	}

	return nil
}

func updateNewSales(start, end time.Time, loc models.Location) error {
	// get all sales from zenoti for period
	zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		return err
	}

	invoices, err := zcli.ReportsAllCollections(start, end)
	if err != nil {
		return err
	}

	// get all guestIds
	var guestIdList []string
	err = models.DB.Model(&models.JpmReportNewLead{}).Where("location_id = ? AND guest_id IS NOT NULL AND guest_id != ''", loc.Id).Distinct().Pluck("guest_id", &guestIdList).Error
	if err != nil {
		return err
	}
	guestIds := make(map[string]bool, len(guestIdList))
	for _, id := range guestIdList {
		guestIds[id] = true
	}

	// if invoice guestId matches contact guestId, save the sale
	for _, invoice := range invoices {
		if invoice.Guest_id != "" && guestIds[invoice.Guest_id] {
			// save the sale
			sale := models.JpmReportInvoice{
				Date:       invoice.Created_Date.Time,
				InvoiceId:  invoice.Invoice_id,
				LocationId: loc.Id,
				GuestId:    strings.ToLower(invoice.Guest_id),
				Total:      invoice.Total_collection,
				IsConsult:  invoice.Total_collection == 0,
			}

			err = models.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "invoice_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"date", "location_id", "guest_id", "total", "is_consult", "updated_at"}),
			}).Create(&sale).Error

			if err != nil {
				return err
			}
		}
	}

	return nil
}
