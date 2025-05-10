package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
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
