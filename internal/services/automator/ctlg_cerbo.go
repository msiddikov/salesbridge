package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_cerbo"
	"client-runaway-zenoti/packages/cerbo"
	"context"
	"fmt"
	"strings"
)

var (
	// Cerbo Category
	cerboCategory = Category{
		Id:    "cerbo",
		Name:  "Cerbo",
		Icon:  "ri:stethoscope-line",
		Color: "#0EA5E9",
		Nodes: []Node{
			cerboActionFindPatient,
			cerboActionCreateEncounter,
		},
	}

	cerboActionFindPatient = Node{
		Id:          "cerbo.patient.find",
		Title:       "Find Patient",
		Description: "Finds a patient in Cerbo by matching fields.",
		ExecFunc:    cerboFindPatient,
		Type:        NodeTypeAction,
		Icon:        "ri:search-2-line",
		Color:       ColorAction,
		Ports: []NodePort{
			customPort("found", cerboPatientFields),
			customPort("notFound", []NodeField{}),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "first_name", Type: "string"},
			{Key: "last_name", Type: "string"},
			{Key: "email", Type: "string"},
			{Key: "username", Type: "string"},
			{Key: "dob", Type: "string"},
		},
	}

	cerboActionCreateEncounter = Node{
		Id:          "cerbo.encounter.create",
		Title:       "Create Encounter",
		Description: "Creates a new encounter on a patient chart.",
		ExecFunc:    cerboCreateEncounter,
		Type:        NodeTypeAction,
		Icon:        "ri:file-add-line",
		Color:       ColorAction,
		Ports: []NodePort{
			successPort(cerboEncounterFields),
			errorPort,
		},
		Fields: []NodeField{
			{Key: "patient_id", Type: "string", Required: true},
			{Key: "encounterType", Type: "string", Required: true},
			{Key: "header", Type: "string"},
			{Key: "content", Type: "string"},
		},
	}

	cerboPatientFields = []NodeField{
		{Key: "id", Type: "string"},
		{Key: "first_name", Type: "string"},
		{Key: "last_name", Type: "string"},
		{Key: "email", Type: "string"},
		{Key: "phone", Type: "string"},
		{Key: "dob", Type: "string"},
	}

	cerboEncounterFields = []NodeField{
		{Key: "id", Type: "string"},
		{Key: "patient_id", Type: "string"},
		{Key: "encounter_type", Type: "string"},
		{Key: "encounter_date", Type: "string"},
		{Key: "status", Type: "string"},
	}
)

func cerboFindPatient(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	first_name, ok := fields["first_name"].(string)
	if !ok {
		first_name = ""
	}
	last_name, ok := fields["last_name"].(string)
	if !ok {
		last_name = ""
	}
	email, ok := fields["email"].(string)
	if !ok {
		email = ""
	}
	username, ok := fields["username"].(string)
	if !ok {
		username = ""
	}
	dob, ok := fields["dob"].(string)
	if !ok {
		dob = ""
	}
	patient, err := svc_cerbo.FindPatient(
		l,
		first_name,
		last_name,
		email,
		username,
		dob,
	)

	if err != nil {
		return errorPayload(err, "failed to find patient")
	}
	if patient == nil {
		return customPayload("notFound", map[string]interface{}{})
	}

	return customPayload("found", mapCerboPatientToPayload(*patient))
}

func cerboCreateEncounter(ctx context.Context, fields map[string]interface{}, l models.Location) map[string]map[string]interface{} {
	patientID := strings.TrimSpace(fmt.Sprint(fields["patient_id"]))
	if patientID == "" {
		return errorPayload(nil, "patient_id is required")
	}

	encounterType := strings.TrimSpace(fmt.Sprint(fields["encounterType"]))
	if encounterType == "" {
		return errorPayload(nil, "encounterType is required")
	}

	header := strings.TrimSpace(fmt.Sprint(fields["header"]))
	content := strings.TrimSpace(fmt.Sprint(fields["content"]))

	encounter, err := svc_cerbo.CreateEncounter(l, patientID, encounterType, header, content)
	if err != nil {
		return errorPayload(err, "failed to create encounter")
	}

	return successPayload(mapCerboEncounterToPayload(encounter))
}

func mapCerboPatientToPayload(p cerbo.Patient) map[string]interface{} {
	return map[string]interface{}{
		"id":         p.Id,
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"email":      p.Email,
		"phone":      p.Phone,
		"dob":        p.Dob,
	}
}

func mapCerboEncounterToPayload(e cerbo.Encounter) map[string]interface{} {
	return map[string]interface{}{
		"id":             e.Id,
		"patient_id":     e.PatientId,
		"encounter_type": e.EncounterType,
		"encounter_date": e.EncounterDate,
		"status":         e.Status,
	}
}
