package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"context"
	"encoding/json"
	"strings"
)

var (
	// Zentoti Category
	zenotiCategory = Category{
		Id:   "zenoti",
		Name: "Zenoti",
		Nodes: []Node{
			zenotiTriggerGuestCreated,
			zenotiTriggerAppointmentCreated,
			zenotiTriggerInvoiceClosed,
			zenotiActionGetGuest,
			zenotiActionFindGuest,
			zenotiActionCreateGuest,
			zenotiActionUpdateGuest,
		},
	}

	// Triggers
	zenotiTriggerGuestCreated = Node{
		Id:          "zenoti.guest.created",
		Title:       "Guest Created",
		Description: "Triggers when a new guest is created in Zenoti.",
		Type:        NodeTypeTrigger,
		Icon:        "ri:form",
		Color:       ColorTrigger,
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: zenotiGuestNodeFields,
			},
		},
	}

	zenotiTriggerAppointmentCreated = Node{
		Id:          "zenoti.appointment.created",
		Title:       "Appointment Created",
		Description: "Triggers when a new appointment is created in Zenoti.",
		Type:        NodeTypeTrigger,
		Icon:        "ri:form",
		Color:       ColorTrigger,
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: zenotiAppointmentNodeFields,
			},
		},
	}

	zenotiTriggerInvoiceClosed = Node{
		Id:          "zenoti.invoice.closed",
		Title:       "Invoice Closed",
		Description: "Triggers when an invoice is closed in Zenoti.",
		Type:        NodeTypeTrigger,
		Icon:        "ri:form",
		Color:       ColorTrigger,
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: zenotiInvoiceNodeFields,
			},
		},
	}

	// Actions
	zenotiActionFindGuest = Node{
		Id:          "zenoti.guest.find",
		Title:       "Find Guest",
		Description: "Finds a guest in Zenoti. with email or phone.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: zenotiGuestNodeFields,
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "email", Type: "string"},
			{Key: "phone", Type: "string"},
		},
	}

	zenotiActionGetGuest = Node{
		Id:          "zenoti.guest.get",
		Title:       "Get Zenoti Guest",
		Description: "Gets a guest in Zenoti by guest ID.",
		ExecFunc:    zenotiActionGetGuestById,
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(zenotiGuestNodeFields),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "guestId", Label: "Guest ID", Type: "string"},
		},
	}

	zenotiActionCreateGuest = Node{
		Id:          "zenoti.guest.create",
		Title:       "Create Guest",
		Description: "Creates a new guest in Zenoti.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: zenotiGuestNodeFields,
			},
			errorPort,
		},
		Fields: zenotiGuestNodeFields,
	}

	zenotiActionUpdateGuest = Node{
		Id:          "zenoti.guest.update",
		Title:       "Update Guest",
		Description: "Updates an existing guest in Zenoti.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: zenotiGuestNodeFields,
			},
			errorPort,
		},
		Fields: zenotiGuestNodeFields,
	}

	//////////////////////////////////////////////////
	//                  Node Fields
	///////////////////////////////////////////////////

	zenotiInvoiceNodeFields = []NodeField{
		{Key: "invoiceId", Type: "string"},
		{Key: "guestId", Type: "string"},
		{Key: "amount", Type: "number"},
		{Key: "status", Type: "string"},
	}

	zenotiAppointmentNodeFields = []NodeField{
		{Key: "appointmentId", Label: "Appointment ID", Type: "string"},
		{Key: "guestId", Label: "Guest ID", Type: "string"},
		{Key: "guestFirstName", Label: "Guest First Name", Type: "string"},
		{Key: "guestLastName", Label: "Guest Last Name", Type: "string"},
		{Key: "guestEmail", Label: "Guest Email", Type: "string"},
		{Key: "centerId", Label: "Center ID", Type: "string"},
		{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
		{Key: "serviceIds", Label: "Service IDs", Type: "string"},
		{Key: "serviceNames", Label: "Service Names", Type: "string"},
		{Key: "date", Label: "Appointment Date", Type: "string"},
	}

	zenotiGuestNodeFields = []NodeField{
		{Key: "id", Label: "Guest ID", Type: "string"},
		{Key: "centerId", Label: "Center ID", Type: "string"},
		{Key: "firstName", Label: "First Name", Type: "string"},
		{Key: "lastName", Label: "Last Name", Type: "string"},
		{Key: "email", Label: "Email", Type: "string"},
		{Key: "phone", Label: "Phone", Type: "string"},
		{Key: "dob", Label: "Date of Birth", Type: "string"},

		{Key: "Address1", Label: "Address 1", Type: "string"},
		{Key: "Address2", Label: "Address 2", Type: "string"},
		{Key: "City", Label: "City", Type: "string"},

		{Key: "receiveMarketingEmails", Label: "Receive Marketing Emails", Type: "boolean"},
		{Key: "receiveTransactionalEmails", Label: "Receive Transactional Emails", Type: "boolean"},
		{Key: "receiveMarketingSms", Label: "Receive Marketing SMS", Type: "boolean"},
		{Key: "receiveTransactionalSms", Label: "Receive Transactional SMS", Type: "boolean"},

		{Key: "tags", Label: "Tags", Type: "string"},
	}
)

//////////////////////////////////////////////////
//
//                  Functions
//
///////////////////////////////////////////////////

func ZenotiTriggerAppointmentCreated(ctx context.Context, WebhookBodyBytes []byte) error {
	type WebhookBody struct {
		Data zenotiv1.AppointmentGroupWebhookData `json:"data"`
	}

	var webhookBody WebhookBody
	if err := json.Unmarshal(WebhookBodyBytes, &webhookBody); err != nil {
		return err
	}

	locs := []models.Location{}
	err := db.DB.Where("zenoti_center_id = ?", webhookBody.Data.Center_Id).Find(&locs).Error
	if err != nil {
		return err
	}

	for _, loc := range locs {

		res := mapZenotiAppointmentGroupToNodePayload(webhookBody.Data)
		triggerInput := TriggerInput{
			LocationID:  loc.Id,
			TriggerType: "zenoti.appointment.created",
			Port:        "out",
			Payload:     res,
		}
		err = StartAutomationsForTrigger(ctx, triggerInput)
		if err != nil {
			return err
		}
	}

	return nil
}

func zenotiActionGetGuestById(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	guestId, ok := fields["guestId"].(string)
	if !ok || guestId == "" {
		return errorPayload(nil, "guestId is required and must be a string")
	}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)

	if err != nil {
		return errorPayload(err, "failed to create zenoti client")
	}

	guest, err := zenotiCli.GuestsGetById(guestId)
	if err != nil {
		return errorPayload(err, "failed to get guest by ID")
	}

	return successPayload(mapZenotiGuestToNodePayload(guest))
}

//////////////////////////////////////////////////
//
//                  Helpers
//
///////////////////////////////////////////////////

func mapZenotiAppointmentGroupToNodePayload(apptGroup zenotiv1.AppointmentGroupWebhookData) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = apptGroup.Appointment_Group_Id
	res["guestId"] = apptGroup.Guest.Id
	res["guestFirstName"] = apptGroup.Guest.FirstName
	res["guestLastName"] = apptGroup.Guest.LastName
	res["guestEmail"] = apptGroup.Guest.Email
	res["centerId"] = apptGroup.Center_Id
	res["invoiceId"] = apptGroup.Invoice_id
	res["serviceIds"] = ""
	res["serviceNames"] = ""

	if len(apptGroup.Appointments) > 0 {
		res["date"] = apptGroup.Appointments[0].Start_time_in_center.Time.Format("2006-01-02")
	}

	for _, appt := range apptGroup.Appointments {
		res["serviceIds"] = res["serviceIds"].(string) + appt.Service_id + ","
		res["serviceNames"] = res["serviceNames"].(string) + appt.Service_Name + ","
	}

	return res
}

func mapZenotiGuestToNodePayload(guest zenotiv1.Guest) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = guest.Id
	res["centerId"] = guest.Center_id
	res["firstName"] = guest.Personal_info.First_name
	res["lastName"] = guest.Personal_info.Last_name
	res["email"] = guest.Personal_info.Email
	res["phone"] = guest.Personal_info.Mobile_phone.Number
	res["dob"] = guest.DateOfBirth.Time.Format("2006-01-02 15:04:05")

	res["Address1"] = guest.Address_info.Address_1
	res["Address2"] = guest.Address_info.Address_2
	res["City"] = guest.Address_info.City

	res["receiveMarketingEmails"] = guest.Preferences.Receive_Marketing_Email
	res["receiveTransactionalEmails"] = guest.Preferences.Receive_Transactional_Email
	res["receiveMarketingSms"] = guest.Preferences.Receive_Marketing_SMS
	res["receiveTransactionalSms"] = guest.Preferences.Receive_Transactional_SMS

	res["tags"] = strings.Join(guest.Tags, ",")
	return res
}
