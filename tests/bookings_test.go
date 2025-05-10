package tests

import (
	"client-runaway-zenoti/internal/runway"
	"testing"
)

func TestBookingWebhook(t *testing.T) {
	payload := runway.AppointmentCreateWebhookPayload{}

	err := runway.CreateBooking(payload)
	if err != nil {
		t.Error(err)
	}

}
