package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"testing"
	"time"
)

func TestCreateBlockSlotDB(t *testing.T) {
	slot := models.BlockSlot{
		LocationId: testLocId,
		CalendarId: testCalendarId,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Hour),
		Title:      "Test Block Slot",
		Notes:      "Test Block Slot Notes",
	}

	err := db.DB.Create(&slot).Error

	if err != nil {
		t.Error(err)
	}
}

func TestEditBlockSlotDB(t *testing.T) {
	slot := models.BlockSlot{Id: "dTaV0jqFBj81W4rkeXdg"}
	db.DB.Where(&slot).Find(&slot)

	slot.Title = "Test Block Slot 2"
	slot.Notes = "Test Block Slot Notes 2"

	err := db.DB.Save(&slot).Error
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteBlocSlotDB(t *testing.T) {
	slot := models.BlockSlot{Id: "dTaV0jqFBj81W4rkeXdg"}
	db.DB.Where(&slot).Find(&slot)
	err := db.DB.Delete(&slot).Error

	if err != nil {
		t.Error(err)
	}

}
