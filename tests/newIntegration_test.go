package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"testing"
)

func TestSyncStagesForLocation(t *testing.T) {

}

func TestTurnOfAllSyncs(t *testing.T) {
	locs := []models.Location{}
	db.DB.Where("sync_contacts = ? or sync_calendar=? or force_check=?", true, true, true).Find(&locs)
	for _, l := range locs {
		l.SyncCalendars = false
		l.SyncContacts = false
		db.DB.Save(&l)
	}
}

func TestUpdateTokens(t *testing.T) {
	runway.UpdateAllTokens()
}
