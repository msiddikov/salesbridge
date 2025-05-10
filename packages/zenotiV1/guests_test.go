package zenotiv1

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdateGuest(t *testing.T) {
	client := getClient()

	err := client.GuestsUpdate("hayden@test.com", "000000000000")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateContact(t *testing.T) {
	client := getClient()
	date, err := time.Parse("2006-01-02", "1990-01-01")
	if err != nil {
		t.Error(err)
	}

	guest := Guest{
		Personal_info: Personal_info{
			First_name:  "Hayden2",
			Last_name:   "Test",
			Email:       "hayden@test2.com",
			DateOfBirth: ZenotiTime{Time: date},
			Mobile_phone: Phone_info{
				Number: "000000000001",
			},
		},
	}

	guest, err = client.GuestsCreate(guest)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(guest)
}

func TestGetByPhoneEmail(t *testing.T) {
	client := getClient()

	guests, err := client.GuestsGetByPhoneEmail("+14708253372", "robynhood1954@yahoo.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(guests)
}
