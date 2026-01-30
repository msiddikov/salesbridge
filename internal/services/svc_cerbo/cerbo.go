package svc_cerbo

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/cerbo"
	"fmt"
	"strings"
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

func CreateEncounter(location models.Location, req cerbo.EncounterCreateRequest) (cerbo.Encounter, error) {
	cli, err := clientForLocation(location)
	if err != nil {
		return cerbo.Encounter{}, err
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
