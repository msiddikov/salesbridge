package tests

import (
	"client-runaway-zenoti/internal/zenoti"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"testing"
)

func TestGetZenotiLink(t *testing.T) {
	locId := "TpwQvq1uDohQXHFebMQj"                   // Fairlawn
	guestId := "f97e6fef-0140-4229-9783-ce7c993ba5c0" // Hayden2 test

	link := zenoti.GuestsGetLinkByIdLocationId(guestId, locId)
	fmt.Println(link)
}

func TestGetCenters(t *testing.T) {
	api := "db42ef1415b24d6792076776c11bdb3dd62ce56639b84d50aa104c2bed6796ee"
	centers, err := zenotiv1.CentersListAll(api)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Found %d centers\n", len(centers))
}
