package tests

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/types"
	zenoti_legacy "client-runaway-zenoti/internal/zenotiLegacy/zenoti"
	"fmt"
	"testing"
)

func Test_GetAppointments(t *testing.T) {
	for _, l := range config.GetLocations() {
		if l.Name != "Fairlawn" {
			continue
		}
		res, err := zenoti_legacy.GetAllAppointments(l)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(res)
	}
}

func Test_GetMembersNo(t *testing.T) {
	for _, l := range config.GetLocations() {
		if l.Name != "Fairlawn" {
			continue
		}
		res, err := zenoti_legacy.GetMembersNo(l)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(res)
	}
}

func Test_GetPastAppointments(t *testing.T) {
	noShows := []types.Appointment{}
	for _, l := range config.GetLocations() {
		if l.Name != "Fairlawn" {
			continue
		}
		res, err := zenoti_legacy.GetPastAppointments(l)
		if err != nil {
			t.Error(err)
		}
		// get only noShows
		for _, v := range res {
			if int(v.Status) < 0 {
				noShows = append(noShows, v)
			}
		}

	}
	fmt.Println(noShows)
}

func TestGetGuestAppointments(t *testing.T) {
	guestEmail := "hayden@test.com"

	loc := config.GetLocationById("TpwQvq1uDohQXHFebMQj") // Fairlawn
	guest, err := zenoti_legacy.SearchGuest("", guestEmail, loc)
	if err != nil {
		t.Error(err)
	}
	res, err := zenoti_legacy.GetGuestAppointments(guest.Id, loc.Zenoti.Api)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
