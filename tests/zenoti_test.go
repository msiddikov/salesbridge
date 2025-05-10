package tests

import (
	"client-runaway-zenoti/internal/zenoti"
	"fmt"
	"testing"
)

func TestGetZenotiLink(t *testing.T) {
	locId := "TpwQvq1uDohQXHFebMQj"                   // Fairlawn
	guestId := "f97e6fef-0140-4229-9783-ce7c993ba5c0" // Hayden2 test

	link := zenoti.GuestsGetLinkByIdLocationId(guestId, locId)
	fmt.Println(link)
}
