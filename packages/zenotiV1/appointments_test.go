package zenotiv1

import (
	"fmt"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	client := getClient()

	start, _ := time.Parse(time.RFC3339, "2023-07-14T00:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2023-07-15T00:00:00Z")

	filter := AppointmentFilter{
		StartDate: start,
		EndDate:   end,
		// TherapistId: "4736cbd2-b556-4822-bb60-b281f9499572", // Krista Elder
		// TherapistId:         "c1c6cf49-903c-442d-a796-eb7b78050b03", // Kaylee O'Donnell
		IncludeNoShowCancel: true,
	}

	res, err := client.AppointmentsListAppointments(filter)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestAppointmentsGetDetails(t *testing.T) {
	client := getFairlawnClient()

	res, err := client.AppointmentsGetDetails("4e46af27-b2dc-49f6-8c3b-aa615370a22a")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestGetAppointment(t *testing.T) {
	clientName := "Natasha"
	client := getFairlawnClient()

	filter := AppointmentFilter{
		StartDate:           time.Now().Add(0 * 24 * time.Hour),
		EndDate:             time.Now().Add(2 * 24 * time.Hour),
		IncludeNoShowCancel: true,
	}

	appts, err := client.AppointmentsListAllAppointments(filter)
	if err != nil {
		t.Error(err)
	}

	for _, a := range appts {
		if a.Guest.First_name == clientName {
			fmt.Println(a)
		}
	}
}
