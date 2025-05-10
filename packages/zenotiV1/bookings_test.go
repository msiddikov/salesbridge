package zenotiv1

import (
	"testing"
	"time"
)

func TestBooking(t *testing.T) {
	cli := getTrainingClient()
	serviceId := "3996df82-da85-4a50-8d9d-95a56348d93b"
	//guestId := "fd751580-3592-4015-9e56-32447f31e359" // test test
	guestId := "cee36dfc-b430-4fa9-8ba1-a0d28a3c1825" // Tailor swift
	therapistId := "e29100f6-819c-4325-8631-432372fbd726"
	datetime, _ := time.Parse(time.RFC3339, "2025-01-20T12:30:00Z")
	date := ZenotiDate{
		Time: datetime,
	}

	bookingReq := BookingReq{
		Date: date,
		Guests: []BookingReqGuest{
			{
				Id: guestId,
				Items: []BookingReqGuestsItems{
					{
						Item: BookingReqItem{
							Id: serviceId,
						},
						Therapist: BookingReqTherapist{
							Id: therapistId,
						},
					},
				},
			},
		},
	}

	book, err := cli.BookWithConfirm(bookingReq)
	if err != nil {
		t.Error(err)
	}
	t.Log(book)

}
