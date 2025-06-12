package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/integrations"
	"client-runaway-zenoti/internal/runway"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"testing"
)

func TestChangeZenotiAPIKey(t *testing.T) {
	oldApi := "57f1dc53e17044b984ea02cbc15b0588d6001d3c22994b2aba92d36169fcfd64"
	newApi := "893ccafc65d44472adca82df5f447217e513715d6de345ed9a46c2bd41e4168d"

	// select all locations with the old api key
	locations := []models.Location{}
	db.DB.Where("zenoti_api = ?", oldApi).Find(&locations)

	// update the api key for each location
	for _, location := range locations {
		location.ZenotiApi = newApi
		db.DB.Save(&location)
	}
}

func TestPushGGuestsToGHL(t *testing.T) {
	locationName := "Young Medical Spa"
	loc := models.Location{}
	db.DB.Where("name = ?", locationName).First(&loc)
	if loc.Id == "" {
		t.Errorf("Location %s not found", locationName)
		return
	}

	svc := runway.GetSvc()
	cli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		t.Errorf("Failed to create client from location ID %s: %v", loc.Id, err)
		return
	}

	zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		t.Errorf("Failed to create Zenoti client: %v", err)
		return
	}

	errs, err := zcli.GuestsIterateAll(56, integrations.PushGuestToGHL, cli)
	if err != nil {
		t.Errorf("Failed to iterate guests: %v", err)
		return
	}
	if len(errs) > 0 {
		t.Errorf("Errors occurred while pushing guests to GHL: %v", errs)
		return
	}
}

func TestPushGGuestsToGHL2(t *testing.T) {
	locationName := "Young Medical Spa - Lansdale"
	loc := models.Location{}
	db.DB.Where("name = ?", locationName).First(&loc)
	if loc.Id == "" {
		t.Errorf("Location %s not found", locationName)
		return
	}

	svc := runway.GetSvc()
	cli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		t.Errorf("Failed to create client from location ID %s: %v", loc.Id, err)
		return
	}

	zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		t.Errorf("Failed to create Zenoti client: %v", err)
		return
	}

	errs, err := zcli.GuestsIterateAll(1, integrations.PushGuestToGHL, cli)
	if err != nil {
		t.Errorf("Failed to iterate guests: %v", err)
		return
	}
	if len(errs) > 0 {
		t.Errorf("Errors occurred while pushing guests to GHL: %v", errs)
		return
	}
}

func TestPushGGuestsToGHL3(t *testing.T) {
	locationName := "Young Medical Spa - Wilkes-Barre/Scranton"
	loc := models.Location{}
	db.DB.Where("name = ?", locationName).First(&loc)
	if loc.Id == "" {
		t.Errorf("Location %s not found", locationName)
		return
	}

	svc := runway.GetSvc()
	cli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		t.Errorf("Failed to create client from location ID %s: %v", loc.Id, err)
		return
	}

	zcli, err := zenotiv1.NewClient(loc.Id, loc.ZenotiCenterId, loc.ZenotiApi)
	if err != nil {
		t.Errorf("Failed to create Zenoti client: %v", err)
		return
	}

	errs, err := zcli.GuestsIterateAll(23, integrations.PushGuestToGHL, cli)
	if err != nil {
		t.Errorf("Failed to iterate guests: %v", err)
		return
	}
	if len(errs) > 0 {
		t.Errorf("Errors occurred while pushing guests to GHL: %v", errs)
		return
	}
}
