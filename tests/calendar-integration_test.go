package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	integrations_zenoti "client-runaway-zenoti/internal/integrations/zenoti"
	"client-runaway-zenoti/internal/runway"
	"strings"
	"testing"
)

func TestCalendarSync(t *testing.T) {
	locationId := "LPanRUHnQMyHqq2O1u8c"
	loc := models.Location{}
	db.DB.Where("id = ?", locationId).First(&loc)

	integrations_zenoti.SyncCalendarsForLocation(loc)
}

func TestDeleteAllEvents(t *testing.T) {
	slots := []models.BlockSlot{}
	db.DB.Where("location_id", testLocId).Find(&slots)

	for _, s := range slots {
		db.DB.Delete(&s)
	}
}

func TestDeleteAllCalendarsRegistered(t *testing.T) {
	locationId := "LPanRUHnQMyHqq2O1u8c"
	loc := models.Location{}
	err := db.DB.Where("id = ?", locationId).First(&loc).Error
	if err != nil {
		t.Error(err)
	}
	if loc.Id == "" {
		t.Error("Location not found")
	}

	calendars := []models.Calendar{}
	db.DB.Where("location_id = ?", loc.Id).Find(&calendars)

	for _, c := range calendars {
		db.DB.Delete(&c)
	}
}

func TestDeleteAllCalendars(t *testing.T) {
	locationId := "LPanRUHnQMyHqq2O1u8c"
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(locationId)

	calendars, err := client.CalendarsGet()

	if err != nil {
		t.Error(err)
	}

	for _, c := range calendars {
		if !strings.Contains(c.Name, "Zenoti (") {
			continue
		}
		err := client.CalendarsDelete(c.Id)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestTurnOffAllIntegrations(t *testing.T) {
	locs := []models.Location{}
	db.DB.Where("sync_calendars=? or sync_contacts=?", true, true).Find(&locs)
	for _, l := range locs {
		l.SyncCalendars = false
		l.SyncContacts = false
		db.DB.Save(&l)
	}
}
