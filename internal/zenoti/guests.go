package zenoti

import (
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
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
