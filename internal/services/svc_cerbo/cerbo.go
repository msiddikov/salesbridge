package svc_cerbo

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/cerbo"
	"fmt"
	"strings"
	"time"
)

func FindPatient(location models.Location, firstName, lastName, email, username, dob string) (*cerbo.Patient, error) {
	cli, err := clientForLocation(location)
	if err != nil {
		return nil, err
	}

	params := cerbo.PatientSearchParams{
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     strings.TrimSpace(email),
		Username:  strings.TrimSpace(username),
		Dob:       strings.TrimSpace(dob),
	}

	patients, err := cli.FindPatients(params)
	if err != nil {
		return nil, err
	}
	if len(patients) == 0 {
		return nil, nil
	}

	return &patients[0], nil
}

func CreateEncounter(location models.Location, patientID, encounterType, header, content string) (cerbo.Encounter, error) {
	cli, err := clientForLocation(location)
	if err != nil {
		return cerbo.Encounter{}, err
	}

	req := cerbo.EncounterCreateRequest{
		EncounterType: strings.TrimSpace(encounterType),
		Title:         strings.TrimSpace(header),
		Content:       strings.TrimSpace(content),
		PatientId:     strings.TrimSpace(patientID),
		DateOfService: time.Now().Format("2006-01-02"),
	}

	return cli.CreateEncounter(req)
}

func clientForLocation(location models.Location) (cerbo.Client, error) {
	if location.CerboApiObjId == 0 {
		return cerbo.Client{}, fmt.Errorf("cerbo api is not configured")
	}

	var api models.CerboApi
	if err := db.DB.First(&api, "id = ?", location.CerboApiObjId).Error; err != nil {
		return cerbo.Client{}, err
	}

	return cerbo.NewClient(api.Subdomain, api.Username, api.ApiKey)
}
