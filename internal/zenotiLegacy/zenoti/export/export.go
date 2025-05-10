package export

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti"
)

func UpdateGuest(email, phone string, l config.Location) error {
	return zenoti.UpdateGuest(email, phone, l)
}
