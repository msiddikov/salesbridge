package runway

import (
	cmn "client-runaway-zenoti/internal/common"
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"time"
)

type (
	set struct {
		Name          string
		LocationId    string
		LocationOldId string
		NewLeadsId    string
		Testing       bool
	}
)

var (
	sets = []set{
		{
			Name:          "Test",
			LocationId:    "ui9MKRggXJMrcUVOiRa8",
			LocationOldId: "Test",
			NewLeadsId:    "0b5b2ed9-2ff9-449a-b9dc-4948d16ebd19",
		},
		{
			Name:          "Solon",
			LocationId:    "TknDfZc7WZF0cq5YE0du",
			LocationOldId: "Solon",
			NewLeadsId:    "b3b75180-5c10-4bcd-972e-3252bf2daa80",
		},
		{
			Name:          "Avon",
			LocationId:    "jEbtFz8wbaNI1WWDDFIb",
			LocationOldId: "Avon",
			NewLeadsId:    "d31f880d-3baf-4474-bfee-81e0af59dcb7",
		},
		{
			Name:          "Tamaya",
			LocationId:    "70CWvxTnNfvgSQBIMlEA",
			LocationOldId: "Tamaya",
			NewLeadsId:    "9cf3bf75-18d8-4261-971c-214c24fa196f",
		},
		{
			Name:          "UpperArlington",
			LocationId:    "RF8DA0atVXW89ZHuaO5l",
			LocationOldId: "UpperArlington",
			NewLeadsId:    "32bfe321-9b4f-474e-8f5c-7a56fb8359fb",
		},
		{
			Name:          "Toledo",
			LocationId:    "u9fPb0fBqTehSErlKLZ6",
			LocationOldId: "Toledo",
			NewLeadsId:    "6900cb3e-bf85-46d5-be20-52a2484444ba",
		},
	}
)

func getWrongBookings(l config.Location, start, end time.Time) ([]runwayv2.Opportunity, error) {
	result := []runwayv2.Opportunity{}

	client, err := svc.NewClient(l.Id, "", "")
	if err != nil {
		return []runwayv2.Opportunity{}, err
	}

	loc := models.Location{}
	loc.Get(l.Id)
	opps, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  start,
		EndDate:    end,
		PipelineId: loc.PipelineId,
		StageId:    loc.BookId,
	})
	if err != nil {
		return []runwayv2.Opportunity{}, err
	}

	sales, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		StartDate:  start,
		EndDate:    end,
		PipelineId: loc.PipelineId,
		StageId:    loc.SalesId,
	})
	if err != nil {
		return []runwayv2.Opportunity{}, err
	}

	opps = append(opps, sales...)
	for _, opp := range opps {
		guest, err := zenoti.SearchGuest(opp.Contact.Phone, opp.Contact.Email, l)
		if err != nil {
			result = append(result, opp)
			continue
		}

		bookings, err := zenoti.GetGuestAppointments(guest.Id, l.Zenoti.Api)

		if err != nil || len(bookings) == 0 {
			result = append(result, opp)
		}
	}

	return result, nil
}

func CheckBookings() {
	fmt.Println("Checking started")
	l := config.GetLocationById("KEs4U6H0jTYhxk47aAKh")
	start := time.Now().Add(0 - 15*24*time.Hour)
	end := time.Now()
	opps, err := getWrongBookings(l, start, end)

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, opp := range opps {
		fmt.Printf("%s %s\n", opp.Name, opp.Contact.Email)
	}
	fmt.Println("Checking done")
}

func contactHasBookingsInZenoti(phone, email string, s set) bool {

	loc := config.GetLocationById(s.LocationOldId)
	if loc.Zenoti.Api == "" {
		return false
	}

	guest, err := zenoti.SearchGuest(zenoti.TrimPhoneNumber(phone), email, loc)
	if err != nil {
		return false // no zenoti contact
	}

	bookings, err := zenoti.GetGuestAppointments(guest.Id, loc.Zenoti.Api)

	if err != nil || len(bookings) == 0 {
		return false // no appointments
	}

	return true
}

func ForceCheckStage(opportunityId string, s set) {

	client, err := svc.NewClient(s.LocationId, "", "")
	if err != nil {
		return
	}

	opp, err := client.OpportunitiesGet(opportunityId)
	if err != nil {
		return
	}

	if contactHasBookingsInZenoti(opp.Contact.Phone, opp.Contact.Email, s) {
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

	if s.Testing {
		msg := fmt.Sprintf("Contact %s in %s should be moved back to new leads\n", opp.Name, s.Name)
		cmn.NotifySlack("", msg)
		return
	}

	update := runwayv2.OpportunityUpdateReq{
		PipelineId:      opp.PipelineId,
		PipelineStageId: s.NewLeadsId,
		Name:            opp.Name,
		Status:          opp.Status,
		MonetaryValue:   opp.MonetaryValue,
		AssignedTo:      opp.AssignedTo,
	}

	_, err = client.OpportunitiesUpdate(update, opp.Id)
	if err != nil {
		msg := fmt.Sprintf("Contact %s from %s error: %s\n", opp.Name, s.Name, err.Error())
		cmn.NotifySlack("", msg)
	}
	msg := fmt.Sprintf("Contact %s from %s has been moved back to new leads\n", opp.Name, s.Name)

	cmn.NotifySlack("", msg)
}

func ForceCheckLocation(s set) {
	fmt.Println("Forcecheck is started for " + s.Name)
	start := time.Now().Add(0 - 30*24*time.Hour)
	end := time.Now()

	loc := models.Location{}
	loc.Get(s.LocationId)

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
		ForceCheckStage(opp.Id, s)
	}

	fmt.Println("Forcecheck is finished")
}

func ForceCheckAll() {
	cmn.NotifySlack("", "Starting forceCheck job...")
	for _, s := range sets {
		ForceCheckLocation(s)
	}
	cmn.NotifySlack("", "Finished forceCheck job")
}
