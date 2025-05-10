package runway

import (
	"client-runaway-zenoti/internal/tgbot"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"encoding/json"
)

type (
	AppointmentCreateWebhookPayload struct {
		Contact_Id string
		First_name string
		Last_name  string
		Email      string
		Phone      string
		Location   struct {
			Id string
		}
		Calendar struct {
			Id               string
			SelectedTimeZone string
			StartTime        zenotiv1.ZenotiTime
		}
	}
)

func HandleOpportunityStageUpdate(payloadBytes []byte) (err error) {
	payload := struct {
		LocationId      string
		Id              string
		ContactId       string
		PipelineId      string
		PipelineStageId string
		Source          string
	}{}

	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return
	}

	//ForceCheckStage(payload.Id, payload.LocationId)
	return
}

func TestOpportunityStageUpdate() {

}

func HandleAppointmentCreate(payloadBytes []byte) (err error) {

	payload := AppointmentCreateWebhookPayload{}

	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return
	}

	tgbot.Notify("New Appointment", "New appointment for "+payload.First_name+" "+payload.Last_name+"\n"+payload.Email+"\n"+payload.Phone+"\n"+payload.Location.Id, false)

	err = CreateBooking(payload)
	if err != nil {
		tgbot.Notify("New Appointment", "Error creating booking: "+err.Error(), false)
	}
	return err
}
