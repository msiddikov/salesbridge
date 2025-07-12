package integrations

import (
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

// Pushes Guest into GHL if not found in it
// params:
// 1. runwayV2 client
func PushGuestToGHL(guest zenotiv1.Guest, params ...any) (err error) {
	if len(params) < 1 {
		return fmt.Errorf("not enough params")
	}

	runwayV2Client, ok := params[0].(runwayv2.Client)
	if !ok {
		return fmt.Errorf("first param must be RunwayV2Client")
	}

	// check if guest is already in GHL
	contacts, err := runwayV2Client.ContactsFindByEmailPhone(guest.Personal_info.Email, guest.Personal_info.Mobile_phone.Number)
	if err != nil {
		return err
	}

	if len(contacts) > 0 {
		// guest already exists in GHL, no need to push
		return nil
	}

	if guest.Personal_info.Mobile_phone.Number == "" {
		return fmt.Errorf("guest phone number is empty")
	}
	if guest.Personal_info.Email == "" {
		return fmt.Errorf("guest email is empty")
	}

	// push guest to GHL
	contact := runwayv2.Contact{
		FirstName:  guest.Personal_info.First_name,
		LastName:   guest.Personal_info.Last_name,
		Email:      guest.Personal_info.Email,
		Phone:      guest.Personal_info.Mobile_phone.Number,
		LocationId: runwayV2Client.GetLocationId(),

		Address1:   guest.Address_info.Address_1,
		City:       guest.Address_info.City,
		PostalCode: guest.Address_info.Zip_code,

		Tags:   []string{"sales-bridge"},
		Source: "Zenoti",
	}
	// dob := guest.Personal_info.DateOfBirth.Time.Format("2006-01-02")
	// contact.CustomFields = []runwayv2.CustomFieldValue{
	// 	{
	// 		Name:        "Date of Birth",
	// 		Field_value: dob,
	// 	},
	// }

	if guest.Personal_info.Gender == 1 {
		contact.Gender = "male"
	}
	if guest.Personal_info.Gender == 2 {
		contact.Gender = "female"
	}

	_, err = runwayV2Client.ContactsCreate(contact)

	return err
}

// Pushes Guest update into GHL if found in it
// params:
// 1. runwayV2 client
// Updates only DND settings for Email and SMS from the Guest
func PushGuestUpdateToGHL(guest zenotiv1.Guest, params ...any) (err error) {
	if len(params) < 1 {
		return fmt.Errorf("not enough params")
	}

	runwayV2Client, ok := params[0].(runwayv2.Client)
	if !ok {
		return fmt.Errorf("first param must be RunwayV2Client")
	}

	// check if guest is already in GHL
	contacts, err := runwayV2Client.ContactsFindByEmailPhone(guest.Personal_info.Email, guest.Personal_info.Mobile_phone.Number)
	if err != nil {
		return err
	}

	if len(contacts) == 0 {
		// guest not found in GHL, no need to push
		return nil
	}

	// push dnd settings to GHL
	contact := runwayv2.Contact{
		Id: contacts[0].Id,
		DndSettings: runwayv2.ContactDndSettings{
			Email: runwayv2.ContactDndSetting{
				Status: lvn.Ternary(guest.Preferences.Receive_Marketing_Email, runwayv2.ContactDndStatusActive, runwayv2.ContactDndStatusInactive),
			},
			Sms: runwayv2.ContactDndSetting{
				Status: lvn.Ternary(guest.Preferences.Receive_Marketing_SMS, runwayv2.ContactDndStatusActive, runwayv2.ContactDndStatusInactive),
			},
		},
	}

	_, err = runwayV2Client.ContactsUpdate(contact)

	return err
}
