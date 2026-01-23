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
			zenotiCollectionCollectionsNew,
			zenotiCollectionCollections,
			zenotiCollectionSales,
			zenotiActionMergeAppointment,
			zenotiActionMergeSales,
			zenotiActionGetInvoice,
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
		Data zenotiv1.Invoice `json:"data"`
	}

	var webhookBody WebhookBody
	if err := json.Unmarshal(WebhookBodyBytes, &webhookBody); err != nil {
		return err
	}

	locs := []models.Location{}
	err := db.DB.Where("zenoti_center_id = ?", webhookBody.Data.Invoice.Center_id).Find(&locs).Error
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

func zenotiCollectAppointments(ctx context.Context, fields map[string]interface{}, l models.Location) (collectionResult, error) {

	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	startDateString, ok := fields["dateFrom"].(string)
	if !ok || startDateString == "" {
		return collectionResult{}, fmt.Errorf("date from is required")
	}
	endDateString, ok := fields["dateTo"].(string)
	if !ok || endDateString == "" {
		return collectionResult{}, fmt.Errorf("date to is required")
	}

	includeNoShowAndCancels, ok := fields["includeNoShowAndCancels"].(bool)
	if !ok {
		includeNoShowAndCancels = false
	}

	startDate, err := parseTime(startDateString)
	if err != nil {
		return collectionResult{}, fmt.Errorf("invalid date from format")
	}
	endDate, err := parseTime(endDateString)
	if err != nil {
		return collectionResult{}, fmt.Errorf("invalid date to format")
	}

	if endDate.Before(startDate) {
		return collectionResult{}, fmt.Errorf("date to must be after date from")
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
		return collectionResult{}, nil
	}

	filterDateEnd := filterDateStart.AddDate(0, 0, 1)

	filter := zenotiv1.AppointmentFilter{
		StartDate:           filterDateStart,
		EndDate:             filterDateEnd,
		IncludeNoShowCancel: includeNoShowAndCancels,
	}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return collectionResult{}, err
	}

	apps, err := zenotiCli.AppointmentsListAppointments(filter)
	if err != nil {
		return collectionResult{}, err
	}

	res := []collectionItem{}
	for _, app := range apps {
		payload := mapZenotiAppointmentToNodePayload(app)
		res = append(res, collectionItem{payload: payload, countsFor: 1})
	}

	estimatedDaily := len(apps)
	if estimatedDaily == 0 {
		estimatedDaily = 1
	}

	return collectionResult{
		items:   res,
		total:   totalDays * estimatedDaily,
		hasMore: true,
	}, nil
}

// Collections
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

func zenotiCollectCollections(ctx context.Context, fields map[string]interface{}, l models.Location) (collectionResult, error) {

	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	startDateString, ok := fields["dateFrom"].(string)
	if !ok || startDateString == "" {
		return collectionResult{}, fmt.Errorf("date from is required")
	}
	endDateString, ok := fields["dateTo"].(string)
	if !ok || endDateString == "" {
		return collectionResult{}, fmt.Errorf("date to is required")
	}

	startDate, err := parseTime(startDateString)
	if err != nil {
		return collectionResult{}, fmt.Errorf("invalid date from format")
	}
	endDate, err := parseTime(endDateString)
	if err != nil {
		return collectionResult{}, fmt.Errorf("invalid date to format")
	}

	if endDate.Before(startDate) {
		return collectionResult{}, fmt.Errorf("date to must be after date from")
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
		return collectionResult{}, nil
	}

	filterDateEnd := filterDateStart

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return collectionResult{}, err
	}

	collections, err := zenotiCli.ReportsCollections(filterDateStart, filterDateEnd)
	if err != nil {
		return collectionResult{}, err
	}

	res := []collectionItem{}
	for _, collection := range collections {
		payload := mapZenotiCollectionToNodePayload(collection)
		res = append(res, collectionItem{payload: payload, countsFor: 1})
	}

	estimatedDaily := len(collections)
	if estimatedDaily == 0 {
		estimatedDaily = 1
	}

	return collectionResult{
		items:   res,
		total:   totalDays * estimatedDaily,
		hasMore: true,
	}, nil
}

// Invoices

var zenotiActionGetInvoice = Node{
	Id:          "zenoti.invoice.get",
	Title:       "Get Zenoti Invoice",
	Description: "Gets an invoice in Zenoti by invoice ID.",
	ExecFunc:    zenotiActionGetInvoiceById,
	Type:        NodeTypeAction,
	Icon:        "ri:form",
	Color:       ColorAction,
	Ports: []NodePort{
		successPort(zenotiInvoiceNodeFields),
		errorPort,
	},
	Fields: []NodeField{
		{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
	},
}

func zenotiActionGetInvoiceById(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	invoiceId, ok := fields["invoiceId"].(string)
	if !ok || invoiceId == "" {
		return errorPayload(nil, "invoiceId is required and must be a string")
	}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)

	if err != nil {
		return errorPayload(err, "failed to create zenoti client")
	}

	invoice, err := zenotiCli.InvoicesGetDetails(invoiceId)
	if err != nil {
		return errorPayload(err, "failed to get invoice by ID")
	}

	return successPayload(mapZenotiInvoiceToNodePayload(invoice))
}

// Sales
var (
	zenotiCollectionSales = Node{
		Id:            "zenoti.collection.sales",
		Title:         "Select Sales",
		Description:   "Selects a Set of Sales from Zenoti.",
		CollectorFunc: zenotiCollectSales,
		Type:          NodeTypeCollection,
		Icon:          "ri:stack",
		Color:         ColorDefault,
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: zenotiCollectionSalesNodeFields,
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "limit", Label: "Limit", Type: "number"},
			{Key: "pageFrom", Label: "Page From", Type: "number"},
			{Key: "pageTo", Label: "Page To", Type: "number"},
			{Key: "startDate", Label: "Start Date", Type: "datetime", Required: true},
			{Key: "endDate", Label: "End Date", Type: "datetime", Required: true},
			{Key: "invoiceStatuses", Label: "Invoice Statuses", Type: "[]string", SelectOptions: []string{"closed", "open"}},
			{Key: "itemTypes", Label: "Item Types", Type: "[]string", SelectOptions: zenotiv1.GetAllItemTypes()},
			{Key: "paymentTypes", Label: "Payment Types", Type: "[]string", SelectOptions: zenotiv1.GetAllPaymentTypes()},
			{Key: "saleTypes", Label: "Sale Types", Type: "[]string", SelectOptions: zenotiv1.GetAllSaleTypes()},
		},
	}

	zenotiCollectionSalesNodeFields = []NodeField{
		{Key: "centerId", Label: "Center ID", Type: "string"},
		{Key: "centerName", Label: "Center Name", Type: "string"},
		{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
		{Key: "saleDate", Label: "Sale Date", Type: "string"},
		{Key: "closedDate", Label: "Closed Date", Type: "string"},

		{Key: "collected", Label: "Collected", Type: "number"},
		{Key: "discount", Label: "Discount", Type: "number"},
		{Key: "redeemed", Label: "Redeemed", Type: "number"},
		{Key: "taxableRedemption", Label: "Taxable Redemption", Type: "number"},
		{Key: "salesExTax", Label: "Sales Ex Tax", Type: "number"},
		{Key: "salesExcludingRedemption", Label: "Sales Excluding Redemption", Type: "number"},
		{Key: "salesIncTax", Label: "Sales Inc Tax", Type: "number"},
		{Key: "status", Label: "Status", Type: "string"},

		{Key: "itemIds", Label: "Item IDs", Type: "string"},
		{Key: "itemNames", Label: "Item Names", Type: "string"},
		{Key: "itemTypes", Label: "Item Types", Type: "string"},

		{Key: "guestId", Label: "Guest ID", Type: "string"},
		{Key: "guestName", Label: "Guest Name", Type: "string"},
	}
)

func mapFieldsToZenotiSalesFilter(fields map[string]interface{}) (zenotiv1.SalesAccrualFilter, error) {
	var filter zenotiv1.SalesAccrualFilter

	startDateStr, ok := fields["startDate"].(string)
	if !ok || startDateStr == "" {
		return filter, fmt.Errorf("startDate is required")
	}
	endDateStr, ok := fields["endDate"].(string)
	if !ok || endDateStr == "" {
		return filter, fmt.Errorf("endDate is required")
	}

	startDate, err := parseTime(startDateStr)
	if err != nil {
		return filter, fmt.Errorf("invalid startDate format")
	}
	endDate, err := parseTime(endDateStr)
	if err != nil {
		return filter, fmt.Errorf("invalid endDate format")
	}

	filter.Start_date = zenotiv1.ZenotiTime{Time: startDate}
	filter.End_date = zenotiv1.ZenotiTime{Time: endDate}

	if v, ok := fields["invoiceStatuses"].([]interface{}); ok {
		for _, status := range v {
			if statusStr, ok := status.(string); ok {
				status := int(-1)
				switch strings.ToLower(statusStr) {
				case "closed":
					status = 4
				case "open":
					status = 0
				}

				filter.Invoice_statuses = append(filter.Invoice_statuses, zenotiv1.InvoiceStatus(status))
			}
		}
	}

	if v, ok := fields["itemTypes"].([]interface{}); ok {
		for _, itemType := range v {
			if itemTypeStr, ok := itemType.(string); ok {
				filter.Item_types = append(filter.Item_types, zenotiv1.ItemType(itemTypeStr))
			}
		}
	}
	if v, ok := fields["paymentTypes"].([]interface{}); ok {
		for _, paymentType := range v {
			if paymentTypeStr, ok := paymentType.(string); ok {
				pt, _ := zenotiv1.GetPaymentType(paymentTypeStr)
				filter.Payment_types = append(filter.Payment_types, pt)
			}
		}
	}
	if v, ok := fields["saleTypes"].([]interface{}); ok {
		for _, saleType := range v {
			if saleTypeStr, ok := saleType.(string); ok {
				st, _ := zenotiv1.GetSaleType(saleTypeStr)
				filter.Sale_types = append(filter.Sale_types, st)
			}
		}
	}

	return filter, nil
}
func mapZenotiSaleToNodePayload(sale zenotiv1.SalesDetails) map[string]interface{} {
	res := make(map[string]interface{})
	res["centerId"] = sale.Center_id
	res["centerName"] = sale.Center_name
	res["invoiceId"] = sale.Invoice_id
	res["saleDate"] = sale.Sale_date.Time.Format(time.RFC3339)
	res["closedDate"] = sale.Invoice_closed_date.Time.Format(time.RFC3339)
	res["saleDate"] = sale.Sale_date.Time.Format(time.RFC3339)

	res["collected"] = sale.Collected
	res["discount"] = sale.Discount
	res["redeemed"] = sale.Redeemed
	res["taxableRedemption"] = sale.Taxable_redemption
	res["salesExTax"] = sale.Sales_ex_tax
	res["salesExcludingRedemption"] = sale.Sales_excluding_redemption
	res["salesIncTax"] = sale.Sales_inc_tax
	res["status"] = sale.Status

	res["itemIds"] = []string{sale.Item_id}
	res["itemNames"] = []string{sale.Item_name}
	res["itemTypes"] = []string{string(sale.Item_type)}

	res["guestId"] = sale.Guest_id
	res["guestName"] = sale.Guest_name

	return res
}

func zenotiCollectSales(ctx context.Context, fields map[string]interface{}, l models.Location) (collectionResult, error) {
	filter, err := mapFieldsToZenotiSalesFilter(fields)
	if err != nil {
		return collectionResult{}, err
	}
	limit, ok := fields["limit"].(float64)
	if !ok || limit < 1 {
		limit = 50
	}
	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	pageInfo := zenotiv1.PageInfo{
		Size: int(limit),
		Page: int(page),
	}

	filter.Center_ids = []string{l.ZenotiCenterId}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return collectionResult{}, err
	}

	sales, pageInfo, err := zenotiCli.ReportsSalesAccrual(filter, pageInfo)
	if err != nil {
		return collectionResult{}, err
	}

	// fetch next page in case if the current invoice is split across pages
	page2 := pageInfo
	page2.Page = page2.Page + 1

	sales2, _, err := zenotiCli.ReportsSalesAccrual(filter, page2)
	if err != nil {
		return collectionResult{}, err
	}

	invoices := map[string][]zenotiv1.SalesDetails{}
	for _, sale := range sales {
		invoices[sale.Invoice_id] = append(invoices[sale.Invoice_id], sale)
	}
	for _, sale := range sales2 {
		_, exists := invoices[sale.Invoice_id]
		if exists {
			invoices[sale.Invoice_id] = append(invoices[sale.Invoice_id], sale)
		}
	}

	res := []collectionItem{}
	for _, salesDetails := range invoices {
		if len(salesDetails) == 0 {
			continue
		}
		resLine := mapZenotiSaleToNodePayload(salesDetails[0])

		for i := 1; i < len(salesDetails); i++ {
			// Append item details
			resLine["itemIds"] = append(resLine["itemIds"].([]string), salesDetails[i].Item_id)
			resLine["itemNames"] = append(resLine["itemNames"].([]string), salesDetails[i].Item_name)
			resLine["itemTypes"] = append(resLine["itemTypes"].([]string), string(salesDetails[i].Item_type))

			// Aggregate amounts
			resLine["collected"] = resLine["collected"].(float64) + salesDetails[i].Collected
			resLine["discount"] = resLine["discount"].(float64) + salesDetails[i].Discount
			resLine["redeemed"] = resLine["redeemed"].(float64) + salesDetails[i].Redeemed
			resLine["taxableRedemption"] = resLine["taxableRedemption"].(float64) + salesDetails[i].Taxable_redemption
			resLine["salesExTax"] = resLine["salesExTax"].(float64) + salesDetails[i].Sales_ex_tax
			resLine["salesExcludingRedemption"] = resLine["salesExcludingRedemption"].(float64) + salesDetails[i].Sales_excluding_redemption
			resLine["salesIncTax"] = resLine["salesIncTax"].(float64) + salesDetails[i].Sales_inc_tax
		}

		// stringify amounts
		resLine["collected"] = fmt.Sprintf("%.2f", resLine["collected"].(float64))
		resLine["discount"] = fmt.Sprintf("%.2f", resLine["discount"].(float64))
		resLine["redeemed"] = fmt.Sprintf("%.2f", resLine["redeemed"].(float64))
		resLine["taxableRedemption"] = fmt.Sprintf("%.2f", resLine["taxableRedemption"].(float64))
		resLine["salesExTax"] = fmt.Sprintf("%.2f", resLine["salesExTax"].(float64))
		resLine["salesExcludingRedemption"] = fmt.Sprintf("%.2f", resLine["salesExcludingRedemption"].(float64))
		resLine["salesIncTax"] = fmt.Sprintf("%.2f", resLine["salesIncTax"].(float64))

		res = append(res, collectionItem{payload: resLine, countsFor: len(salesDetails)})
	}

	return collectionResult{
		items:   res,
		total:   pageInfo.Total,
		hasMore: pageInfo.Page*pageInfo.Size < pageInfo.Total,
	}, nil

}

// Collections - new report

var zenotiCollectionCollectionsNew = Node{
	Id:            "zenoti.collection.collectionNew",
	Title:         "Select Collections with new report",
	Description:   "Selects a Set of Collections from Zenoti from new report.",
	CollectorFunc: zenotiCollectCollectionsNew,
	Type:          NodeTypeCollection,
	Icon:          "ri:stack",
	Color:         ColorDefault,
	Ports: []NodePort{
		{
			Name:    "out",
			Payload: zenotiCollectionDetailNodeFields,
		},
		errorPort,
	},
	Fields: zenotiCollectionsDetailsFilterNodeFields,
}

func mapFieldsToZenotiCollectionsFilter(fields map[string]interface{}) (zenotiv1.CollectionsFilter, error) {
	var filter zenotiv1.CollectionsFilter

	startDateStr, ok := fields["startDate"].(string)
	if !ok || startDateStr == "" {
		return filter, fmt.Errorf("startDate is required")
	}
	endDateStr, ok := fields["endDate"].(string)
	if !ok || endDateStr == "" {
		return filter, fmt.Errorf("endDate is required")
	}

	startDate, err := parseTime(startDateStr)
	if err != nil {
		return filter, fmt.Errorf("invalid startDate format")
	}
	endDate, err := parseTime(endDateStr)
	if err != nil {
		return filter, fmt.Errorf("invalid endDate format")
	}

	filter.Start_date = zenotiv1.ZenotiTime{Time: startDate}
	filter.End_date = zenotiv1.ZenotiTime{Time: endDate}

	if v, ok := fields["invoiceStatuses"].([]interface{}); ok {
		for _, status := range v {
			if statusStr, ok := status.(string); ok {
				status := int(-1)
				switch strings.ToLower(statusStr) {
				case "closed":
					status = 4
				case "open":
					status = 0
				}

				filter.Invoice_statuses = append(filter.Invoice_statuses, zenotiv1.InvoiceStatus(status))
			}
		}
	}
	if v, ok := fields["paymentTypes"].([]interface{}); ok {
		for _, paymentType := range v {
			if paymentTypeStr, ok := paymentType.(string); ok {
				pt, _ := zenotiv1.GetPaymentType(paymentTypeStr)
				filter.Payment_types = append(filter.Payment_types, pt)
			}
		}
	}
	if v, ok := fields["saleTypes"].([]interface{}); ok {
		for _, saleType := range v {
			if saleTypeStr, ok := saleType.(string); ok {
				st, _ := zenotiv1.GetSaleType(saleTypeStr)
				filter.Sale_types = append(filter.Sale_types, st)
			}
		}
	}

	return filter, nil
}

var zenotiCollectionsDetailsFilterNodeFields = []NodeField{
	{Key: "limit", Label: "Limit", Type: "number"},
	{Key: "pageFrom", Label: "Page From", Type: "number"},
	{Key: "pageTo", Label: "Page To", Type: "number"},
	{Key: "startDate", Label: "Start Date", Type: "datetime"},
	{Key: "endDate", Label: "End Date", Type: "datetime"},
	{Key: "invoiceStatuses", Label: "Invoice Statuses", Type: "[]string", SelectOptions: []string{"closed", "open"}},
	{Key: "paymentTypes", Label: "Payment Types", Type: "[]string", SelectOptions: zenotiv1.GetAllPaymentTypes()},
	{Key: "saleTypes", Label: "Sale Types", Type: "[]string", SelectOptions: zenotiv1.GetAllSaleTypes()},
}

func mapZenotiCollectionDetailsToNodePayload(collection zenotiv1.CollectionDetails) map[string]interface{} {
	res := make(map[string]interface{})
	res["centerId"] = collection.Center_id
	res["centerName"] = collection.Center_name
	res["invoiceId"] = collection.Invoice_id
	res["invoiceClosedDate"] = collection.Invoice_closed_date.Time.Format(time.RFC3339)

	res["totalPaid"] = fmt.Sprintf("%.2f", collection.Total_paid)

	res["guestId"] = collection.Guest_id
	res["guestName"] = collection.Guest_name

	return res
}

var zenotiCollectionDetailNodeFields = []NodeField{
	{Key: "centerId", Label: "Center ID", Type: "string"},
	{Key: "centerName", Label: "Center Name", Type: "string"},
	{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
	{Key: "invoiceClosedDate", Label: "Invoice Closed Date", Type: "string"},

	{Key: "totalPaid", Label: "Total Paid", Type: "number"},

	{Key: "guestId", Label: "Guest ID", Type: "string"},
	{Key: "guestName", Label: "Guest Name", Type: "string"},
}

func zenotiCollectCollectionsNew(ctx context.Context, fields map[string]interface{}, l models.Location) (collectionResult, error) {
	filter, err := mapFieldsToZenotiCollectionsFilter(fields)
	if err != nil {
		return collectionResult{}, err
	}
	limit, ok := fields["limit"].(float64)
	if !ok || limit < 1 {
		limit = 50
	}
	page, ok := fields["page"].(float64)
	if !ok || page < 1 {
		page = 1
	}

	pageInfo := zenotiv1.PageInfo{
		Size: int(limit),
		Page: int(page),
	}

	filter.Centers.IDs = []string{l.ZenotiCenterId}

	zenotiCli, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApiObj.ApiKey)
	if err != nil {
		return collectionResult{}, err
	}

	collections, pageInfo, err := zenotiCli.ReportsNewCollections(filter, pageInfo)
	if err != nil {
		return collectionResult{}, err
	}

	res := []collectionItem{}
	for _, collection := range collections {
		payload := mapZenotiCollectionDetailsToNodePayload(collection)
		res = append(res, collectionItem{payload: payload, countsFor: 1})
	}

	return collectionResult{
		items:   res,
		total:   pageInfo.Total,
		hasMore: pageInfo.Page*pageInfo.Size < pageInfo.Total,
	}, nil

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

	{Key: "serviceIds", Label: "Service IDs", Type: "string"},
	{Key: "serviceNames", Label: "Service Names", Type: "string"},
}

func mapZenotiInvoiceClosedWebhookToNodePayload(data zenotiv1.Invoice) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = data.Invoice.Id
	res["date"] = data.Invoice.Invoice_date.Time.Format("2006-01-02 15:04:05")
	res["amount"] = toString(data.Invoice.Total_price.Sum_total)
	res["guestId"] = data.Invoice.Guest.Id
	res["guestFirstName"] = data.Invoice.Guest.First_name
	res["guestLastName"] = data.Invoice.Guest.Last_name
	res["guestEmail"] = data.Invoice.Guest.Email
	res["guestPhone"] = data.Invoice.Guest.Mobile_phone
	res["serviceIds"] = ""
	res["serviceNames"] = ""

	for _, item := range data.Invoice.Invoice_Items {
		res["serviceIds"] = res["serviceIds"].(string) + item.Id + ","
		res["serviceNames"] = res["serviceNames"].(string) + item.Name + ","
	}

	return res
}

func mapZenotiInvoiceToNodePayload(data zenotiv1.Invoice) map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = data.Invoice.Id
	res["date"] = data.Invoice.Invoice_date.Time.Format("2006-01-02 15:04:05")
	res["amount"] = toString(data.Invoice.Total_price.Sum_total)

	res["guestId"] = data.Guest.Id
	res["guestFirstName"] = data.Guest.First_name
	res["guestLastName"] = data.Guest.Last_name
	res["guestPhone"] = data.Guest.Mobile_phone

	res["serviceIds"] = ""
	res["serviceNames"] = ""

	for _, item := range data.Invoice_Items {
		res["serviceIds"] = res["serviceIds"].(string) + item.Id + ","
		res["serviceNames"] = res["serviceNames"].(string) + item.Name + ","
	}

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
