package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	integrations_zenoti "client-runaway-zenoti/internal/integrations/zenoti"
	"client-runaway-zenoti/internal/runway"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func TestUpdateStages(t *testing.T) {
	l := models.Location{}
	locationName := "Hand & Stone Massage and Facial Spa - Chesapeake"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	integrations_zenoti.UpdateStagesForLocation(l)
}

func TestUpdateSales(t *testing.T) {
	l := models.Location{}
	locationName := "Hand and Stone Massage and Facial Spa - Olathe"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	l.SalesSyncDate = time.Now().AddDate(0, 0, -5)

	err := integrations_zenoti.UpdateSalesChunks(l)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateCollectionAllLocations(t *testing.T) {
	locations := []models.Location{}
	// find locations starting with "Hand and Stone Massage and Facial Spa"
	db.DB.Where("name LIKE ?", "Hand%").Find(&locations)
	fmt.Println(locations)

	updatedLocations := []string{}

	for _, l := range locations {
		// if location is in the updatedLocations list
		if slices.Contains(updatedLocations, l.Name) {
			return
		}

		l.SalesSyncDate = time.Now().AddDate(0, 0, -5)

		err := integrations_zenoti.UpdateSalesChunks(l)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateNotesNew(t *testing.T) {
	l := models.Location{}
	locationName := "VIO Med Spa Carmel"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	integrations_zenoti.UpdateNotesV2Location(l, true)
}

func TestForceCheckV2(t *testing.T) {
	l := models.Location{}
	locationName := "VIO Med Spa Carmel"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	runway.ForceCheckLocationV2(l)
}

func TestOppsWithNotes(t *testing.T) {
	l := models.Location{}
	locationName := "VIO Med Spa Carmel"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	_, ress, err := integrations_zenoti.OpportunitiesWithNotes(l, "01/03/2024: https://viomedspa.zenoti.com/Appointment/DlgAppointment1.aspx?invoiceid=86e77970-dec0-4a6c-aed3-fe1901dc51d7")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ress)
}

func TestScheduledJobs(t *testing.T) {
	integrations_zenoti.StartScheduledJobs()
}

func TestUpdateSingleAppt(t *testing.T) {
	apptId := "e513e6f2-881d-4b65-8426-c7c7156a69d4"

	l := models.Location{}
	locationName := "VIO Med Spa Fairlawn (DWY)"
	db.DB.Where("name = ?", locationName).First(&l)
	if l.Id == "" {
		t.Error("Location not found")
	}

	client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		t.Error(err)
	}

	rService := runway.GetSvc()
	runwayCli, err := rService.NewClientFromId(l.Id)

	appt, err := client.AppointmentsGetDetails(apptId)
	if err != nil {
		t.Error(err)
	}

	err = runway.UpdateAppt(appt[0], l, runwayCli)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateCollection(t *testing.T) {
	invoiceId := "30854a40-d431-4892-8cbb-1e5c7185ab15"
	date, _ := time.Parse("1/02/2006", "1/26/2025")

	// getting location
	l := models.Location{}
	locationName := "Hand and Stone Massage and Facial Spa - Doral"
	db.DB.Where("name = ?", locationName).First(&l)
	if l.Id == "" {
		t.Error("Location not found")
	}

	// getting zenoti client
	client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		t.Error(err)
	}

	// getting runway client
	rService := runway.GetSvc()
	runwayCli, err := rService.NewClientFromId(l.Id)

	// getting collections
	neededCollection := zenotiv1.Collection{}
	collections, _ := client.ReportsCollections(date.Add(-24*time.Hour), date.Add(24*time.Hour))
	for _, c := range collections {
		if c.Invoice_id == invoiceId {
			neededCollection = c
			break
		}
	}

	// filling in guest info
	guest, err := client.GuestsGetById(neededCollection.Guest_id)
	if err != nil {
		t.Error(err)
	}
	neededCollection.Guest = guest

	_, err = runway.UpdateCollection(neededCollection, l, runwayCli)
	if err != nil {
		t.Error(err)
	}
}
