package integrations_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"testing"
)

func TestSyncWithEmptyEmail(t *testing.T) {
	// finding location
	l := models.Location{}
	locationName := "VIO Med Spa Carmel"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		t.Error("Location not found")
	}

	// getting appointments
	apts, err := getAppointments(l)
	if err != nil {
		t.Error(err)
	}

	// filtering appointments
	newappts := []zenotiv1.Appointment{}
	for _, a := range apts {
		if a.Guest.Email == "" && a.Guest.Mobile.Number != "" {
			newappts = append(newappts, a)
		}
	}

	// updating bookings
	fmt.Println("Updating bookings for " + l.Name)
	err = runway.UpdateBookings(newappts, l)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Booking successfully updated for " + l.Name)
}
