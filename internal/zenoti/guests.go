package zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/integrations"
	"client-runaway-zenoti/internal/runway"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"encoding/json"
	"errors"
)

func GuestsGetByPhoneNumberLocationId(phone, email, locId string) (guest zenotiv1.Guest, err error) {
	loc := models.Location{}
	loc.Get(locId)
	client := MustGetClientFromLocation(loc)

	guests, err := client.GuestsGetByPhoneEmail(phone, email)
	if err != nil {
		return
	}
	return guests[0], nil
}

func GuestsGetLinkByIdLocationId(id, locId string) string {
	loc := models.Location{}
	loc.Get(locId)

	return loc.ZenotiUrl + "/Guests/GuestProfile.aspx?UserId=" + id
}

func GuestCreatedWebhookHandler(WebhookBodyBytes []byte) error {
	type WebhookBody struct {
		Data zenotiv1.Guest `json:"data"`
	}

	var webhookBody WebhookBody
	if err := json.Unmarshal(WebhookBodyBytes, &webhookBody); err != nil {
		return err
	}

	// Process the guest data
	guest := webhookBody.Data
	// TODO: Implement your logic here

	loc := models.Location{}
	db.DB.Where("zenoti_center_id = ?", guest.Center_id).First(&loc)
	if loc.Id == "" {
		return errors.New("location not found for center ID: " + guest.Center_id)

	}

	svc := runway.GetSvc()
	cli, err := svc.NewClientFromId(loc.Id)
	if err != nil {
		return errors.New("failed to create client from location ID: " + loc.Id + ", error: " + err.Error())
	}

	err = integrations.PushGuestToGHL(guest, cli)
	if err != nil {
		return errors.New("failed to push guest to GHL: " + err.Error())
	}

	return nil
}
