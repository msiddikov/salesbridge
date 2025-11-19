package automator

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	// Zentoti Category
	zenotiCategory = Category{
		Id:    "zenoti",
		Name:  "Zenoti",
		Icon:  "ri:health-book-line",
		Color: "#f557c5ff",
		Nodes: []Node{
			zenotiTriggerGuestCreated,
			zenotiTriggerAppointmentCreated,
			zenotiTriggerInvoiceClosed,
			zenotiCollectionAppointments,
			zenotiCollectionCollections,
			zenotiActionMergeAppointment,
			zenotiActionMergeSales,
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

	//////////////////////////////////////////////////
	//                  Collections
	///////////////////////////////////////////////////

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

	zenotiAppointmentNodeFields = []NodeField{
		{Key: "appointmentId", Label: "Appointment ID", Type: "string"},
		{Key: "guestId", Label: "Guest ID", Type: "string"},
		{Key: "guestFirstName", Label: "Guest First Name", Type: "string"},
		{Key: "guestLastName", Label: "Guest Last Name", Type: "string"},
		{Key: "guestEmail", Label: "Guest Email", Type: "string"},
		{Key: "guestPhone", Label: "Guest Phone", Type: "string"},
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

//////////////////////////////////////////////////
//                  Triggers
///////////////////////////////////////////////////

var zenotiTriggerAppointmentCreated = Node{
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

var zenotiTriggerInvoiceClosed = Node{
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

func ZenotiTriggerInvoiceClosed(ctx context.Context, WebhookBodyBytes []byte) error {
	type WebhookBody struct {
		Data zenotiv1.WebhookDataPayload `json:"data"`
	}

	var webhookBody WebhookBody
	if err := json.Unmarshal(WebhookBodyBytes, &webhookBody); err != nil {
		return err
	}

	locs := []models.Location{}
	err := db.DB.Where("zenoti_center_id = ?", webhookBody.Data.Invoice.Center_Id).Find(&locs).Error
	if err != nil {
		return err
	}

	for _, loc := range locs {

		res := mapZenotiInvoiceClosedWebhookToNodePayload(webhookBody.Data)
		triggerInput := TriggerInput{
			LocationID:  loc.Id,
			TriggerType: "zenoti.invoice.closed",
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

//////////////////////////////////////////////////
//                  Guests
///////////////////////////////////////////////////

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
//                  collectors
///////////////////////////////////////////////////

// Appointments
var zenotiCollectionAppointments = Node{
	Id:            "zenoti.appointment.collection",
	Title:         "Select Appointments",
	Description:   "Selects a Collection of Appointments from Zenoti.",
	CollectorFunc: zenotiCollectAppointments,
	Type:          NodeTypeCollection,
	Icon:          "ri:stack",
	Color:         ColorDefault,
	Ports: []NodePort{
		{
			Name:    "out",
			Payload: zenotiAppointmentNodeFields,
		},
		errorPort,
	},
	Fields: []NodeField{
		{Key: "dateFrom", Type: "datetime"},
		{Key: "dateTo", Type: "datetime"},
		{Key: "includeNoShowAndCancels", Type: "bool"},
	},
}

func zenotiCollectAppointments(ctx context.Context, fields map[string]interface{}, l models.Location) ([]map[string]interface{}, int, bool, error) {

	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	startDateString, ok := fields["dateFrom"].(string)
	if !ok || startDateString == "" {
		return nil, 0, false, fmt.Errorf("date from is required")
	}
	endDateString, ok := fields["dateTo"].(string)
	if !ok || endDateString == "" {
		return nil, 0, false, fmt.Errorf("date to is required")
	}

	includeNoShowAndCancels, ok := fields["includeNoShowAndCancels"].(bool)
	if !ok {
		includeNoShowAndCancels = false
	}

	startDate, err := parseTime(startDateString)
	if err != nil {
		return nil, 0, false, fmt.Errorf("invalid date from format")
	}
	endDate, err := parseTime(endDateString)
	if err != nil {
		return nil, 0, false, fmt.Errorf("invalid date to format")
	}

	if endDate.Before(startDate) {
		return nil, 0, false, fmt.Errorf("date to must be after date from")
	}

	loc := startDate.Location()
	endDate = endDate.In(loc)

	startDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)

	totalDays := int(endDay.Sub(startDay)/(24*time.Hour)) + 1
	if totalDays < 1 {
		totalDays = 1
	}

	fields["pageTo"] = float64(totalDays)

	offset := int(page) - 1
	if offset < 0 {
		offset = 0
	}

	filterDateStart := startDay.AddDate(0, 0, offset)
	if filterDateStart.After(endDay) {
		return nil, 0, false, nil
	}

	filterDateEnd := filterDateStart.AddDate(0, 0, 1)

	filter := zenotiv1.AppointmentFilter{
		StartDate:           filterDateStart,
		EndDate:             filterDateEnd,
		IncludeNoShowCancel: includeNoShowAndCancels,
	}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return nil, 0, false, err
	}

	apps, err := zenotiCli.AppointmentsListAppointments(filter)
	if err != nil {
		return nil, 0, false, err
	}

	res := []map[string]interface{}{}
	for _, app := range apps {
		payload := mapZenotiAppointmentToNodePayload(app)
		res = append(res, payload)
	}

	estimatedDaily := len(apps)
	if estimatedDaily == 0 {
		estimatedDaily = 1
	}

	return res, totalDays * estimatedDaily, true, nil
}

// Sales
var zenotiCollectionCollections = Node{
	Id:            "zenoti.collection.collection",
	Title:         "Select Collections",
	Description:   "Selects a Set of Collections from Zenoti.",
	CollectorFunc: zenotiCollectCollections,
	Type:          NodeTypeCollection,
	Icon:          "ri:stack",
	Color:         ColorDefault,
	Ports: []NodePort{
		{
			Name:    "out",
			Payload: zenotiCollectionNodeFields,
		},
		errorPort,
	},
	Fields: []NodeField{
		{Key: "dateFrom", Type: "datetime"},
		{Key: "dateTo", Type: "datetime"},
	},
}

func zenotiCollectCollections(ctx context.Context, fields map[string]interface{}, l models.Location) ([]map[string]interface{}, int, bool, error) {

	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	startDateString, ok := fields["dateFrom"].(string)
	if !ok || startDateString == "" {
		return nil, 0, false, fmt.Errorf("date from is required")
	}
	endDateString, ok := fields["dateTo"].(string)
	if !ok || endDateString == "" {
		return nil, 0, false, fmt.Errorf("date to is required")
	}

	startDate, err := parseTime(startDateString)
	if err != nil {
		return nil, 0, false, fmt.Errorf("invalid date from format")
	}
	endDate, err := parseTime(endDateString)
	if err != nil {
		return nil, 0, false, fmt.Errorf("invalid date to format")
	}

	if endDate.Before(startDate) {
		return nil, 0, false, fmt.Errorf("date to must be after date from")
	}

	loc := startDate.Location()
	endDate = endDate.In(loc)

	startDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)

	totalDays := int(endDay.Sub(startDay)/(24*time.Hour)) + 1
	if totalDays < 1 {
		totalDays = 1
	}

	fields["pageTo"] = float64(totalDays)

	offset := int(page) - 1
	if offset < 0 {
		offset = 0
	}

	filterDateStart := startDay.AddDate(0, 0, offset)
	if filterDateStart.After(endDay) {
		return nil, 0, false, nil
	}

	filterDateEnd := filterDateStart.AddDate(0, 0, 1)

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return nil, 0, false, err
	}

	collections, err := zenotiCli.ReportsCollections(filterDateStart, filterDateEnd)
	if err != nil {
		return nil, 0, false, err
	}

	res := []map[string]interface{}{}
	for _, collection := range collections {
		payload := mapZenotiCollectionToNodePayload(collection)
		res = append(res, payload)
	}

	estimatedDaily := len(collections)
	if estimatedDaily == 0 {
		estimatedDaily = 1
	}

	return res, totalDays * estimatedDaily, true, nil
}

var (
	zenotiActionMergeAppointment = Node{
		Id:          "zenoti.appointment.merge",
		Title:       "Merge Appointments",
		Description: "Merges Appointments in Zenoti.",
		ExecFunc:    mergeActionFunc,
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(zenotiAppointmentNodeFields),
		},
		Fields: zenotiAppointmentNodeFields,
	}

	zenotiActionMergeSales = Node{
		Id:          "zenoti.sales.merge",
		Title:       "Merge Sales",
		Description: "Merges Sales in Zenoti.",
		ExecFunc:    mergeActionFunc,
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(zenotiInvoiceNodeFields),
		},
		Fields: zenotiInvoiceNodeFields,
	}
)

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

func mapZenotiAppointmentToNodePayload(appt zenotiv1.Appointment) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = appt.Id
	res["guestId"] = appt.Guest.Id
	res["guestFirstName"] = appt.Guest.First_name
	res["guestLastName"] = appt.Guest.Last_name
	res["guestEmail"] = appt.Guest.Email
	res["invoiceId"] = appt.Invoice_id
	res["serviceIds"] = appt.Service.Id
	res["serviceNames"] = appt.Service.Name
	res["date"] = appt.Start_time.Time.Format("2006-01-02")

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

var zenotiInvoiceNodeFields = []NodeField{
	{Key: "id", Label: "Id", Type: "string"},
	{Key: "date", Label: "Date", Type: "string"},
	{Key: "amount", Label: "Amount", Type: "string"},

	{Key: "guestId", Label: "Guest ID", Type: "string"},
	{Key: "guestFirstName", Label: "Guest First Name", Type: "string"},
	{Key: "guestLastName", Label: "Guest Last Name", Type: "string"},
	{Key: "guestEmail", Label: "Guest Email", Type: "string"},
	{Key: "guestPhone", Label: "Guest Phone", Type: "string"},
}

func mapZenotiInvoiceClosedWebhookToNodePayload(data zenotiv1.WebhookDataPayload) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = data.Invoice.Id
	res["date"] = data.Invoice.Invoice_Date.Time.Format("2006-01-02 15:04:05")
	res["amount"] = toString(data.Invoice.Total_Price.Sum_Total)

	res["guestId"] = data.Invoice.Guest.Id
	res["guestFirstName"] = data.Invoice.Guest.First_Name
	res["guestLastName"] = data.Invoice.Guest.Last_Name
	res["guestEmail"] = data.Invoice.Guest.Email
	res["guestPhone"] = data.Invoice.Guest.Mobile_Phone

	return res
}

var zenotiCollectionNodeFields = []NodeField{
	{Key: "InvoiceId", Label: "Invoice Id", Type: "string"},
	{Key: "date", Label: "Date", Type: "string"},
	{Key: "amount", Label: "Amount", Type: "string"},

	{Key: "guestId", Label: "Guest ID", Type: "string"},
}

func mapZenotiCollectionToNodePayload(collection zenotiv1.Collection) map[string]interface{} {
	res := make(map[string]interface{})
	res["InvoiceId"] = collection.Invoice_id
	res["date"] = collection.Created_Date.Time.Format("2006-01-02 15:04:05")
	res["amount"] = toString(collection.Total_collection)

	res["guestId"] = collection.Guest_id

	return res
}
