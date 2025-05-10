package runway

import (
	cmn "client-runaway-zenoti/internal/common"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"time"
)

func contactHasBookingsInZenotiV2(phone, email string, loc models.Location) bool {

	client, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		return false // no zenoti integration
	}

	guests, err := client.GuestsGetByPhoneEmail(phone, email)
	if err != nil || len(guests) == 0 {
		return false // no zenoti contact
	}

	bookings, _, err := client.GuestsListAppointments(zenotiv1.GuestAppointmentsFilter{
		GuestId: guests[0].Id,
		Page:    1,
		Size:    1,
	})

	if err != nil || len(bookings) == 0 {
		return false // no appointments
	}

	return true
}

func ForceCheckStageV2(opportunityId string, loc models.Location) {

	client, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	opp, err := client.OpportunitiesGet(opportunityId)
	if err != nil {
		return
	}

	if contactHasBookingsInZenotiV2(opp.Contact.Phone, opp.Contact.Email, loc) {
		return
	}

	contact := models.Contact{
		ContactId:  opp.Contact.Id,
		LocationId: opp.Contact.LocationId,
	}

	hasSales, err := contact.HasSalesOrAppointments(db.DB)
	if hasSales || err != nil {
		return
	}

	update := runwayv2.OpportunityUpdateReq{
		PipelineId:      opp.PipelineId,
		PipelineStageId: loc.NewId,
		Name:            opp.Name,
		Status:          opp.Status,
		MonetaryValue:   opp.MonetaryValue,
		AssignedTo:      opp.AssignedTo,
	}

	_, err = client.OpportunitiesUpdate(update, opp.Id)
	if err != nil {
		msg := fmt.Sprintf("Contact %s from %s error: %s\n", opp.Name, loc.Name, err.Error())
		cmn.NotifySlack("", msg)
	}
	msg := fmt.Sprintf("Contact %s from %s has been moved back to new leads\n", opp.Name, loc.Name)

	cmn.NotifySlack("", msg)
}

func ForceCheckLocationV2(loc models.Location) {
	fmt.Println("Forsecheck is started for " + loc.Name)
	start := time.Now().Add(0 - 30*24*time.Hour)
	end := time.Now()

	client, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	opps, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  start,
		EndDate:    end,
		PipelineId: loc.PipelineId,
		StageId:    loc.BookId,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	sales, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  start,
		EndDate:    end,
		PipelineId: loc.PipelineId,
		StageId:    loc.SalesId,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	opps = append(opps, sales...)

	for _, opp := range opps {
		ForceCheckStageV2(opp.Id, loc)
	}

	fmt.Println("Forcecheck is finished")
}

func ForceCheckAllV2() {
	cmn.NotifySlack("", "Starting forceCheck job...")
	locs := []models.Location{}
	err := db.DB.Where("force_check = ?", true).Find(&locs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, l := range locs {
		ForceCheckLocationV2(l)
	}
	cmn.NotifySlack("", "Finished forceCheck job")
}
