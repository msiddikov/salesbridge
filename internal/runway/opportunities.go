package runway

import (
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"strings"
	"time"
)

func CheckAutoCreateOpportunity(opp runwayv2.Opportunity, l models.Location) (zenotiv1.Guest, error) {
	zguest := zenotiv1.Guest{}

	if !l.AutoCreateContacts {
		return zguest, nil
	}

	// get the zenoti client
	c, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		fmt.Println(err)
		return zguest, err
	}

	// check if the guest already exists
	guests, err := c.GuestsGetByPhoneEmail(opp.Contact.Phone, opp.Contact.Email)
	if err != nil && err.Error() != "guest not found" {
		fmt.Println(err)
		return zguest, err
	}

	if len(guests) > 0 {
		return guests[0], nil
	}
	// prepare guest info
	guest, err := PrepareGuestInfoFromTheOpportunity(opp)
	if err != nil {
		return zguest, err
	}

	// create the guest if it doesn't exist

	return c.GuestsCreate(guest)

}

func PrepareGuestInfoFromTheOpportunity(opp runwayv2.Opportunity) (zenotiv1.Guest, error) {
	guest := zenotiv1.Guest{
		Personal_info: zenotiv1.Personal_info{
			First_name: opp.Contact.FirstName,
			Last_name:  opp.Contact.LastName,
			Gender:     1,
			Email:      opp.Contact.Email,
			Mobile_phone: zenotiv1.Phone_info{
				Number: opp.Contact.Phone,
			},
			DateOfBirth: zenotiv1.ZenotiTime{Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		Preferences: zenotiv1.GuestPreferences{
			Receive_Transactional_SMS:   true,
			Receive_Transactional_Email: true,
		},
		Tags: []string{"JFM"},
	}

	if guest.Personal_info.First_name == "" && opp.Name != "" {
		// if the first name is empty, parse the opportunity name by space and use the first word as the first name
		parsedName := strings.Fields(opp.Name)
		if len(parsedName) > 0 {
			guest.Personal_info.First_name = parsedName[0]
		}
	}

	if guest.Personal_info.Last_name == "" && opp.Name != "" {
		// if the last name is empty, parse the opportunity name by space and use the rest of the words as the last name
		parsedName := strings.Fields(opp.Name)
		if len(parsedName) > 1 {
			guest.Personal_info.Last_name = parsedName[1]
		}
	}

	if guest.Personal_info.First_name == "" || guest.Personal_info.Last_name == "" {
		return guest, fmt.Errorf("first name or last name is empty")
	}

	if guest.Personal_info.Email == "" || guest.Personal_info.Mobile_phone.Number == "" {
		return guest, fmt.Errorf("email or phone is empty")
	}

	return guest, nil
}
