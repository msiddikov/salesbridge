package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_ghl"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"context"
	"strconv"
	"strings"
)

var (
	// GHL Category
	ghlCategory = Category{
		Id:    "ghl",
		Name:  "GoHighLevel",
		Icon:  "ri:share-forward-2",
		Color: "#4C6FFF",
		Nodes: []Node{
			ghlTriggerOpportunityCreated,

			ghlCollectionOpportunities,

			ghlActionFindOpportunity,
			ghlActionCreateOpportunity,
			ghlActionUpdateOpportunity,
			ghlActionMergeOpportunity,

			ghlActionDeleteNotes,
			ghlActionUpdateLinkNote,
			ghlActionRegisterBookingNote,
			ghlActionRegisterSalesNote,

			ghlActionGetContact,
			ghlActionFindContact,
			ghlActionCreateContact,
			ghlActionUpdateContact,
		},
	}

	//////////////////////////////////////////////////
	//
	//                  Definitions
	//
	///////////////////////////////////////////////////

	//////////////////////////////////////////////////
	//                  Triggers
	///////////////////////////////////////////////////
	ghlTriggerOpportunityCreated = Node{
		Id:          "ghl.opportunity.created",
		Title:       "Opportunity Created",
		Description: "Triggers when a new opportunity is created in GoHighLevel.",
		Type:        NodeTypeTrigger,
		Icon:        "ri:file-add-line",
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: ghlOpportunityNodeFields,
			},
		},
	}

	//////////////////////////////////////////////////
	//                  Collections
	///////////////////////////////////////////////////

	ghlCollectionOpportunities = Node{
		Id:            "ghl.opportunity.collection",
		Title:         "Select Opportunities",
		Description:   "Selects a Collection of Opportunities from GoHighLevel.",
		Type:          NodeTypeCollection,
		CollectorFunc: ghlCollectionGetOpportunities,
		Icon:          "ri:stack-line",
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: ghlOpportunityNodeFields,
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "limit", Type: "number"},
			{Key: "pageFrom", Type: "number"},
			{Key: "pageTo", Type: "number"},
			{Key: "email", Type: "string"},
			{Key: "phone", Type: "string"},
			{Key: "query", Type: "string"},
			{Key: "pipelineId", Type: "string"},
			{Key: "createdAtFrom", Type: "datetime"},
			{Key: "createdAtTo", Type: "datetime"},
		},
	}

	//////////////////////////////////////////////////
	//                  Actions
	///////////////////////////////////////////////////

	//
	//////////////	Opportunities	////////////////////
	//

	ghlActionFindOpportunity = Node{
		Id:          "ghl.opportunity.find",
		Title:       "Find Opportunity",
		Description: "Finds an Opportunity in GoHighLevel by phone or email.",
		ExecFunc:    ghlFindOpportunityByEmailOrPhone,
		Type:        NodeTypeAction,
		Icon:        "ri:search-2-line",
		Kind:        "Opportunities",
		Color:       ColorAction,
		Ports: []NodePort{
			customPort("found", ghlOpportunityNodeFields),
			customPort("notFound", []NodeField{}),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "email", Type: "string"},
			{Key: "phone", Type: "string"},
			{Key: "pipelineId", Type: "string"},
		},
	}

	ghlActionCreateOpportunity = Node{
		Id:          "ghl.opportunity.create",
		Title:       "Create Opportunity",
		Description: "Creates Opportunity in GoHighLevel.",
		ExecFunc:    ghlCreateOpportunity,
		Type:        NodeTypeAction,
		Icon:        "ri:file-add-line",
		Kind:        "Opportunities",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(ghlOpportunityNodeFields),
			errorPort,
		},
		Fields: ghlOpportunityCreateFields,
	}

	ghlActionUpdateOpportunity = Node{
		Id:          "ghl.opportunity.update",
		Title:       "Update Opportunity",
		Description: "Updates Opportunity in GoHighLevel.",
		ExecFunc:    ghlActionsUpdateOpportunity,
		Type:        NodeTypeAction,
		Icon:        "ri:file-edit-line",
		Kind:        "Opportunities",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(ghlOpportunityNodeFields),
			errorPort,
		},
		Fields: ghlOpportunityUpdateFields,
	}

	ghlActionMergeOpportunity = Node{
		Id:          "ghl.opportunity.merge",
		Title:       "Merge Opportunities",
		Description: "Merges Opportunities in GoHighLevel.",
		ExecFunc:    mergeActionFunc,
		Type:        NodeTypeAction,
		Icon:        "ri:git-fork-line",
		Kind:        "Opportunities",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(ghlOpportunityNodeFields),
		},
		Fields: ghlOpportunityNodeFields,
	}

	//
	//////////////	Contacts	////////////////////
	//

	ghlActionFindContact = Node{
		Id:          "ghl.contact.find",
		Title:       "Find Contact",
		Description: "Finds a Contact in GoHighLevel by phone or email.",
		ExecFunc:    ghlFindContactByEmailOrPhone,
		Type:        NodeTypeAction,
		Icon:        "ri:search-2-line",
		Kind:        "Contacts",
		Color:       ColorAction,
		Ports: []NodePort{
			customPort("found", ghlContactFields),
			customPort("notFound", []NodeField{}),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "email", Type: "string"},
			{Key: "phone", Type: "string"},
		},
	}

	ghlActionCreateContact = Node{
		Id:          "ghl.contact.create",
		Title:       "Create Contact",
		Description: "Creates Contact in GoHighLevel.",
		ExecFunc:    ghlCreateContact,
		Type:        NodeTypeAction,
		Icon:        "ri:file-add-line",
		Kind:        "Contacts",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(ghlContactFields),
			errorPort,
		},
		Fields: ghlContactFields,
	}

	ghlActionUpdateContact = Node{
		Id:          "ghl.contact.update",
		Title:       "Update Contact",
		Description: "Updates Contact in GoHighLevel.",
		ExecFunc:    ghlUpdateContact,
		Type:        NodeTypeAction,
		Icon:        "ri:file-edit-line",
		Kind:        "Contacts",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(ghlContactFields),
			errorPort,
		},
		Fields: ghlContactFields,
	}

	//
	//////////////	Notes	////////////////////
	//

	ghlActionDeleteNotes = Node{
		Id:          "ghl.opportunity.notes.delete",
		Title:       "Delete Contact Notes",
		Description: "Deletes Contact Notes in GoHighLevel.",
		Type:        NodeTypeAction,
		ExecFunc:    ghlActionsDeleteNotes,
		Icon:        "ri:file-reduce-line",
		Kind:        "Notes",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name: "success",
				Payload: []NodeField{
					{Key: "message", Type: "string"},
				},
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "contactId", Label: "Contact ID", Type: "string"},
			{Key: "bookingNotes", Label: "Booking notes", Type: "bool"},
			{Key: "salesNotes", Label: "Sales notes", Type: "bool"},
			{Key: "linkNotes", Label: "Link notes", Type: "bool"},
		},
	}

	ghlActionUpdateLinkNote = Node{
		Id:          "ghl.opportunity.notes.updateLink",
		Title:       "Update Link Note",
		Description: "Updates Link Note in GoHighLevel.",
		ExecFunc:    ghlActionsUpdateLinkNote,
		Type:        NodeTypeAction,
		Icon:        "ri:file-edit-line",
		Kind:        "Notes",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name: "success",
				Payload: []NodeField{
					{Key: "message", Type: "string"},
				},
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "contactId", Label: "Contact ID", Type: "string"},
			{Key: "lazy", Label: "Lazy update", Type: "bool"},
		},
	}

	ghlActionRegisterBookingNote = Node{
		Id:          "ghl.opportunity.notes.registerBooking",
		Title:       "Register Booking Note",
		Description: "Registers Booking Note in GoHighLevel.",
		ExecFunc:    ghlActionsRegisterBookingNote,
		Type:        NodeTypeAction,
		Icon:        "ri:calendar-check-line",
		Kind:        "Notes",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name: "success",
				Payload: []NodeField{
					{Key: "wasAlreadyRegistered", Label: "Was already registered", Type: "bool"},
				},
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "contactId", Label: "Contact ID", Type: "string"},
			{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
			{Key: "date", Label: "Date", Type: "string"},
		},
	}

	ghlActionRegisterSalesNote = Node{
		Id:          "ghl.opportunity.notes.registerSales",
		Title:       "Register Sales Note",
		Description: "Registers Sales Note in GoHighLevel.",
		ExecFunc:    ghlActionsRegisterSalesNote,
		Type:        NodeTypeAction,
		Icon:        "ri:exchange-dollar-line",
		Kind:        "Notes",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name: "success",
				Payload: []NodeField{
					{Key: "wasAlreadyRegistered", Label: "Was already registered", Type: "bool"},
				},
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "contactId", Label: "Contact ID", Type: "string"},
			{Key: "invoiceId", Label: "Invoice ID", Type: "string"},
			{Key: "value", Label: "Value", Type: "string"},
			{Key: "date", Label: "Date", Type: "string"},
		},
	}

	//////////////////////////////////////////////////
	//                  Node Fields
	///////////////////////////////////////////////////

	ghlOpportunityNodeFields = []NodeField{
		{Key: "opportunityId", Type: "string"},
		{Key: "opportunityName", Type: "string"},
		{Key: "monetaryValue", Type: "string"},
		{Key: "source", Type: "string"},
		{Key: "status", Type: "string", SelectOptions: []string{"open", "won", "lost", "abandoned"}},
		{Key: "stageId", Type: "string"},
		{Key: "pipelineId", Type: "string"},
		{Key: "email", Type: "string"},
		{Key: "phone", Type: "string"},
		{Key: "contactId", Type: "string"},
		{Key: "contactFirstName", Type: "string"},
		{Key: "contactLastName", Type: "string"},
		{Key: "contactEmail", Type: "string"},
		{Key: "contactPhone", Type: "string"},
	}

	ghlOpportunityUpdateFields = []NodeField{
		{Key: "opportunityId", Type: "string"},
		{Key: "updateName", Label: "Update name", Type: "bool"},
		{Key: "name", Label: "Name", Type: "string"},
		{Key: "updateMonetaryValue", Label: "Update Monetary Value", Type: "bool"},
		{Key: "monetaryValue", Label: "Monetary Value", Type: "string"},
		{Key: "updateStatus", Label: "Update Status", Type: "bool"},
		{Key: "status", Label: "Status", Type: "string"},
		{Key: "updateStageId", Label: "Update Stage", Type: "bool"},
		{Key: "stageId", Label: "Stage ID", Type: "string"},
		{Key: "updatePipelineId", Label: "Update Pipeline", Type: "bool"},
		{Key: "pipelineId", Label: "Pipeline ID", Type: "string"},
		{Key: "updateAssignedTo", Label: "Update Assigned To", Type: "bool"},
		{Key: "assignedTo", Label: "Assigned To", Type: "string"},
	}

	ghlOpportunityCreateFields = []NodeField{
		{Key: "contactId", Type: "string"},
		{Key: "name", Label: "Name", Type: "string"},
		{Key: "monetaryValue", Label: "Monetary Value", Type: "string"},
		{Key: "status", Label: "Status", Type: "string", SelectOptions: []string{"open", "won", "lost", "abandoned"}},
		{Key: "stageId", Label: "Stage ID", Type: "string"},
		{Key: "pipelineId", Label: "Pipeline ID", Type: "string"},
		{Key: "assignedTo", Label: "Assigned To", Type: "string"},
	}

	ghlContactFields = []NodeField{
		{Key: "Id", Type: "string"},
		{Key: "email", Label: "Email", Type: "string"},
		{Key: "phone", Label: "Phone", Type: "string"},
		{Key: "firstName", Label: "First Name", Type: "string"},
		{Key: "lastName", Label: "Last Name", Type: "string"},
		{Key: "locationId", Label: "Location ID", Type: "string"},

		{Key: "address", Label: "Address", Type: "string"},
		{Key: "city", Label: "City", Type: "string"},
		{Key: "state", Label: "State", Type: "string"},
		{Key: "postalCode", Label: "Postal Code", Type: "string"},
		{Key: "country", Label: "Country", Type: "string"},

		{Key: "dndSms", Label: "DND SMS", Type: "string", SelectOptions: []string{"enabled", "disabled"}},
		{Key: "dndEmail", Label: "DND Email", Type: "string", SelectOptions: []string{"enabled", "disabled"}},

		{Key: "tags", Label: "Tags (comma separated)", Type: "string"},

		{Key: "source", Label: "Source", Type: "string"},
	}

	svc = runwayv2.Service{}
)

//////////////////////////////////////////////////
//
//                  Functions
//
///////////////////////////////////////////////////

//////////////////////////////////////////////////
//                  Contacts
///////////////////////////////////////////////////

func ghlFindContactByEmailOrPhone(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {

	email, emailOk := fields["email"].(string)
	phone, phoneOk := fields["phone"].(string)

	if !emailOk && !phoneOk || (email == "" && phone == "") {
		return errorPayload(nil, "either email or phone is required")
	}

	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	contacts, err := cli.ContactsFindByEmailPhone(email, phone)
	if err != nil {
		return errorPayload(err, "failed to find contacts by email or phone")
	}

	if len(contacts) == 0 {
		return customPayload("notFound", map[string]interface{}{})
	}

	return customPayload("found", mapGhlContactToNodePayload(contacts[0]))
}

func ghlCreateContact(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	contact := mapNodeFieldsToGhlContact(fields)
	contact.LocationId = l.Id

	createdContact, err := cli.ContactsCreate(contact)
	if err != nil {
		return errorPayload(err, "failed to create contact")
	}

	return successPayload(mapGhlContactToNodePayload(createdContact))
}

func ghlUpdateContact(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	contact := mapNodeFieldsToGhlContact(fields)
	if contact.Id == "" {
		return errorPayload(nil, "contact Id is required")
	}

	updatedContact, err := cli.ContactsUpdate(contact)
	if err != nil {
		return errorPayload(err, "failed to update contact")
	}

	return successPayload(mapGhlContactToNodePayload(updatedContact))
}

var ghlActionGetContact = Node{
	Id:          "ghlGetContact",
	Title:       "Get Contact by ID",
	Description: "Gets a Contact in GoHighLevel by ID.",
	ExecFunc:    ghlGetContactById,
	Type:        NodeTypeAction,
	Icon:        "ri:search-2-line",
	Kind:        "Contacts",
	Color:       ColorAction,
	Ports: []NodePort{
		successPort(ghlContactFields),
		errorPort,
	},
	Fields: []NodeField{
		{Key: "id", Type: "string"},
	},
}

func ghlGetContactById(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	id, ok := fields["id"].(string)
	if !ok || id == "" {
		return errorPayload(nil, "id is required")
	}

	contact, err := cli.ContactsGet(id)
	if err != nil {
		return errorPayload(err, "failed to get contact")
	}

	return successPayload(mapGhlContactToNodePayload(contact))
}

//////////////////////////////////////////////////
//                  Opportunities
///////////////////////////////////////////////////

func GhlTriggerOpportunityCreated(ctx context.Context, opportunity runwayv2.Opportunity, l models.Location) error {

	triggerInput := TriggerInput{
		LocationID:  l.Id,
		TriggerType: "ghl.opportunity.created",
		Port:        "out",
		Payload:     map[string]interface{}{},
	}
	triggerInput.Payload = mapGhlOpportunityToNodePayload(opportunity)

	err := StartAutomationsForTrigger(ctx, triggerInput)

	return err
}

func ghlCollectionGetOpportunities(ctx context.Context, fields map[string]interface{}, l models.Location) ([]map[string]interface{}, int, bool, error) {
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return nil, 0, false, err
	}

	page, ok := fields["page"].(float64)
	if !ok {
		page = 1
	}
	limit, ok := fields["limit"].(float64)
	if !ok {
		limit = 20
	}

	filter := runwayv2.OpportunitiesAdvancedFilter{
		LocationId: l.Id,
		Page:       int(page),
		Limit:      int(limit),
	}

	query, ok := fields["query"].(string)
	if ok && query != "" {
		filter.Query = query
	}

	fg := runwayv2.Filter{
		Group:   "AND",
		Filters: []runwayv2.Filter{},
	}

	email, ok := fields["email"].(string)
	if ok && email != "" {
		fg.Filters = append(fg.Filters, runwayv2.Filter{
			Field:    "email",
			Operator: "equals",
			Value:    email,
		})
	}

	phone, ok := fields["phone"].(string)
	if ok && phone != "" {
		fg.Filters = append(fg.Filters, runwayv2.Filter{
			Field:    "phone",
			Operator: "equals",
			Value:    phone,
		})
	}

	if pipelineId, ok := fields["pipelineId"].(string); ok && pipelineId != "" {
		fg.Filters = append(fg.Filters, runwayv2.Filter{
			Field:    "pipelineId",
			Operator: "equals",
			Value:    pipelineId,
		})
	}

	createdAtFrom, okFrom := fields["createdAtFrom"].(string)
	createdAtTo, okTo := fields["createdAtTo"].(string)
	if okFrom && okTo {
		fg.Filters = append(fg.Filters, runwayv2.Filter{
			Field:    "date_added",
			Operator: "range",
			ValueRange: map[string]string{
				"gte": createdAtFrom,
				"lte": createdAtTo,
			},
		})
	}

	opps, total, err := cli.OpportunitiesGetByPagination(filter)
	if err != nil {
		return nil, 0, false, err
	}

	res := []map[string]interface{}{}
	for _, opportunity := range opps {
		payload := mapGhlOpportunityToNodePayload(opportunity)
		res = append(res, payload)
	}

	return res, total, len(res) != 0, nil
}

func ghlActionsUpdateOpportunity(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	res := make(map[string]map[string]interface{})
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}
	id, ok := fields["opportunityId"].(string)
	if !ok || id == "" {
		return errorPayload(nil, "opportunityId is required")
	}

	req, err := mapGhlUpdateOpportunityFieldsToGhlUpdateOpportunityReq(fields)
	if err != nil {
		return errorPayload(err, "failed to map update opportunity fields")
	}
	opp, err := cli.OpportunitiesUpdate(req, id)
	if err != nil {
		return errorPayload(err, "failed to update opportunity")
	}

	res["success"] = mapGhlOpportunityToNodePayload(opp)
	return res
}

func ghlFindOpportunityByEmailOrPhone(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {

	email, emailOk := fields["email"].(string)
	phone, phoneOk := fields["phone"].(string)
	pipelineId, pipelineIdOk := fields["pipelineId"].(string)

	if !emailOk && !phoneOk || (email == "" && phone == "") {
		return errorPayload(nil, "either email or phone is required")
	}

	if !pipelineIdOk || pipelineId == "" {
		return errorPayload(nil, "pipelineId is required")
	}

	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}
	opps, err := cli.OpportunitiesFindByEmailPhone(email, phone, pipelineId)
	if err != nil {
		return errorPayload(err, "failed to find opportunities by email or phone")
	}

	if len(opps) == 0 {
		return customPayload("notFound", map[string]interface{}{})
	}

	return customPayload("found", mapGhlOpportunityToNodePayload(opps[0]))
}

func ghlCreateOpportunity(ctx context.Context, fields map[string]interface{}, l models.Location) (payload map[string]map[string]interface{}) {
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	req, err := mapGhlCreateOpportunityFieldsToGhlCreateOpportunityReq(fields)
	req.LocationId = l.Id
	if err != nil {
		return errorPayload(err, "failed to map create opportunity fields")
	}

	opp, err := cli.OpportunitiesCreate(req)
	if err != nil {
		return errorPayload(err, "failed to create opportunity")
	}

	return successPayload(mapGhlOpportunityToNodePayload(opp))
}

//////////////////////////////////////////////////
//                  Notes
///////////////////////////////////////////////////

func ghlActionsDeleteNotes(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	res := make(map[string]map[string]interface{})
	cli, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return errorPayload(err, "failed to create GHL client")
	}

	contactId, ok := fields["contactId"].(string)
	if !ok || contactId == "" {
		return errorPayload(nil, "contactId is required")
	}

	notes, err := cli.ContactsGetAllNotes(contactId)
	if err != nil {
		return errorPayload(err, "failed to get contact notes")
	}

	deleteBookingNotes, _ := fields["bookingNotes"].(bool)
	deleteSalesNotes, _ := fields["salesNotes"].(bool)
	deleteLinkNotes, _ := fields["linkNotes"].(bool)

	for _, note := range notes {
		if deleteBookingNotes && strings.Contains(note.Body, "Bookings for this contact:") {
			err := cli.ContactsDeleteNote(note)
			if err != nil {
				return errorPayload(err, "failed to delete booking note")
			}
		}
		if deleteSalesNotes && strings.Contains(note.Body, "Collections for this contact:") {
			err := cli.ContactsDeleteNote(note)
			if err != nil {
				return errorPayload(err, "failed to delete sales note")
			}
		}
		if deleteLinkNotes && (strings.Contains(note.Body, "Please follow this link") || strings.Contains(note.Body, "Zenoti update:")) {
			err := cli.ContactsDeleteNote(note)
			if err != nil {
				return errorPayload(err, "failed to delete link note")
			}
		}
	}
	res["success"] = map[string]interface{}{
		"message": "Notes deleted successfully",
	}

	return res
}

func ghlActionsUpdateLinkNote(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {

	contactId, ok := fields["contactId"].(string)
	if !ok || contactId == "" {
		return errorPayload(nil, "contactId is required")
	}

	lazy, _ := fields["lazy"].(bool)

	err := svc_ghl.UpdateLinkNote(contactId, l, !lazy)
	if err != nil {
		return errorPayload(err, "failed to update link note")
	}

	res := map[string]interface{}{
		"message": "Link note updated successfully",
	}

	return successPayload(res)
}

func ghlActionsRegisterBookingNote(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {

	contactId, ok := fields["contactId"].(string)
	if !ok || contactId == "" {
		return errorPayload(nil, "contactId is required")
	}

	invoiceId, ok := fields["invoiceId"].(string)
	if !ok || invoiceId == "" {
		return errorPayload(nil, "invoiceId is required")
	}

	dateStr, ok := fields["date"].(string)
	if !ok || dateStr == "" {
		return errorPayload(nil, "date is required")
	}

	date, err := parseTime(dateStr)
	if err != nil {
		return errorPayload(err, "failed to parse date")
	}

	wasAlreadyRegistered, err := svc_ghl.RegisterBooking(contactId, invoiceId, date, l)
	if err != nil {
		return errorPayload(err, "failed to register booking note")
	}

	res := map[string]interface{}{
		"wasAlreadyRegistered": wasAlreadyRegistered,
	}

	return successPayload(res)
}

func ghlActionsRegisterSalesNote(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {

	contactId, ok := fields["contactId"].(string)
	if !ok || contactId == "" {
		return errorPayload(nil, "contactId is required")
	}

	invoiceId, ok := fields["invoiceId"].(string)
	if !ok || invoiceId == "" {
		return errorPayload(nil, "invoiceId is required")
	}

	valueStr, ok := fields["value"].(string)
	if !ok || valueStr == "" {
		return errorPayload(nil, "value is required")
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return errorPayload(err, "failed to parse value")
	}

	dateStr, ok := fields["date"].(string)
	if !ok || dateStr == "" {
		return errorPayload(nil, "date is required")
	}

	date, err := parseTime(dateStr)
	if err != nil {
		return errorPayload(err, "failed to parse date")
	}

	wasAlreadyRegistered, err := svc_ghl.RegisterCollection(contactId, invoiceId, value, date, l)
	if err != nil {
		return errorPayload(err, "failed to register sales note")
	}

	res := map[string]interface{}{
		"wasAlreadyRegistered": wasAlreadyRegistered,
	}

	return successPayload(res)
}

//////////////////////////////////////////////////
//
//                  Mappers
//
///////////////////////////////////////////////////

func mapGhlOpportunityToNodePayload(opportunity runwayv2.Opportunity) map[string]interface{} {
	payload := map[string]interface{}{}
	payload["opportunityId"] = opportunity.Id
	payload["opportunityName"] = opportunity.Name
	payload["monetaryValue"] = opportunity.MonetaryValue
	payload["source"] = opportunity.Source
	payload["status"] = string(opportunity.Status)
	payload["stageId"] = opportunity.PipelineStageId
	payload["pipelineId"] = opportunity.PipelineId
	payload["email"] = opportunity.Contact.Email
	payload["phone"] = opportunity.Contact.Phone
	payload["contactId"] = opportunity.Contact.Id
	payload["contactFirstName"] = opportunity.Contact.FirstName
	payload["contactLastName"] = opportunity.Contact.LastName
	payload["contactEmail"] = opportunity.Contact.Email
	payload["contactPhone"] = opportunity.Contact.Phone

	return payload
}

func mapGhlUpdateOpportunityFieldsToGhlUpdateOpportunityReq(fields map[string]interface{}) (runwayv2.OpportunityUpdateReq, error) {
	req := runwayv2.OpportunityUpdateReq{}
	var err error

	if updateName, ok := fields["updateName"].(bool); ok {
		req.UpdateName = updateName
	}
	if name, ok := fields["name"].(string); ok {
		req.Name = name
	}
	if updateMonetaryValue, ok := fields["updateMonetaryValue"].(bool); ok {
		req.UpdateMonetaryValue = updateMonetaryValue
	}
	if monetaryValue, ok := fields["monetaryValue"].(string); ok {
		req.MonetaryValue, err = strconv.ParseFloat(monetaryValue, 64)
		if err != nil {
			return req, err
		}
	}
	if updateStatus, ok := fields["updateStatus"].(bool); ok {
		req.UpdateStatus = updateStatus
	}
	if status, ok := fields["status"].(string); ok {
		req.Status = runwayv2.OpportunityStatus(status)
	}
	if updateStageId, ok := fields["updateStageId"].(bool); ok {
		req.UpdatePipelineStageId = updateStageId
	}
	if stageId, ok := fields["stageId"].(string); ok {
		req.PipelineStageId = stageId
	}
	if updatePipelineId, ok := fields["updatePipelineId"].(bool); ok {
		req.UpdatePipelineId = updatePipelineId
	}
	if pipelineId, ok := fields["pipelineId"].(string); ok {
		req.PipelineId = pipelineId
	}
	if updateAssignedTo, ok := fields["updateAssignedTo"].(bool); ok {
		req.UpdateAssignedTo = updateAssignedTo
	}
	if assignedTo, ok := fields["assignedTo"].(string); ok {
		req.AssignedTo = assignedTo
	}

	return req, nil
}

func mapGhlCreateOpportunityFieldsToGhlCreateOpportunityReq(fields map[string]interface{}) (runwayv2.OpportunityCreateReq, error) {
	req := runwayv2.OpportunityCreateReq{}
	var err error

	if name, ok := fields["name"].(string); ok {
		req.Name = name
	}
	if monetaryValue, ok := fields["monetaryValue"].(string); ok {
		req.MonetaryValue, err = strconv.ParseFloat(monetaryValue, 64)
		if err != nil {
			return req, err
		}
	}
	if status, ok := fields["status"].(string); ok {
		req.Status = runwayv2.OpportunityStatus(status)
	}
	if stageId, ok := fields["stageId"].(string); ok {
		req.PipelineStageId = stageId
	}
	if pipelineId, ok := fields["pipelineId"].(string); ok {
		req.PipelineId = pipelineId
	}
	if contactId, ok := fields["contactId"].(string); ok {
		req.ContactId = contactId
	}
	if assignedTo, ok := fields["assignedTo"].(string); ok {
		req.AssignedTo = assignedTo
	}

	return req, nil
}

func mapGhlContactToNodePayload(contact runwayv2.Contact) map[string]interface{} {
	payload := map[string]interface{}{}
	payload["Id"] = contact.Id
	payload["email"] = contact.Email
	payload["phone"] = contact.Phone
	payload["firstName"] = contact.FirstName
	payload["lastName"] = contact.LastName
	payload["locationId"] = contact.LocationId

	payload["address"] = contact.Address1
	payload["city"] = contact.City
	payload["state"] = contact.State
	payload["postalCode"] = contact.PostalCode
	payload["country"] = contact.Country

	if contact.DndSettings.Sms != nil {
		payload["dndSms"] = contact.DndSettings.Sms.Status
	}
	if contact.DndSettings.Email != nil {
		payload["dndEmail"] = contact.DndSettings.Email.Status
	}

	payload["tags"] = strings.Join(contact.Tags, ", ")

	payload["source"] = contact.Source

	return payload
}

func mapNodeFieldsToGhlContact(fields map[string]interface{}) runwayv2.Contact {
	contact := runwayv2.Contact{}

	if id, ok := fields["Id"].(string); ok {
		contact.Id = id
	}
	if email, ok := fields["email"].(string); ok {
		contact.Email = email
	}
	if phone, ok := fields["phone"].(string); ok {
		contact.Phone = phone
	}
	if firstName, ok := fields["firstName"].(string); ok {
		contact.FirstName = firstName
	}
	if lastName, ok := fields["lastName"].(string); ok {
		contact.LastName = lastName
	}
	if locationId, ok := fields["locationId"].(string); ok {
		contact.LocationId = locationId
	}

	if address, ok := fields["address"].(string); ok {
		contact.Address1 = address
	}
	if city, ok := fields["city"].(string); ok {
		contact.City = city
	}
	if state, ok := fields["state"].(string); ok {
		contact.State = state
	}
	if postalCode, ok := fields["postalCode"].(string); ok {
		contact.PostalCode = postalCode
	}
	if country, ok := fields["country"].(string); ok {
		contact.Country = country
	}

	if dndSms, ok := fields["dndSms"].(string); ok {
		contact.DndSettings.Sms.Status = runwayv2.ContactDndStatus(dndSms)
	}
	if dndEmail, ok := fields["dndEmail"].(string); ok {
		contact.DndSettings.Email.Status = runwayv2.ContactDndStatus(dndEmail)
	}

	if tags, ok := fields["tags"].(string); ok {
		contact.Tags = strings.Split(tags, ",")
		for i := range contact.Tags {
			contact.Tags[i] = strings.TrimSpace(contact.Tags[i])
		}
	}

	if source, ok := fields["source"].(string); ok {
		contact.Source = source
	}

	return contact

}

var ()
