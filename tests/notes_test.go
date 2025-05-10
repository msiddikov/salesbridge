package tests

import (
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"testing"
)

func TestCreateNote(t *testing.T) {
	_, contact, client := GetTestingData()

	note := runwayv2.Note{
		ContactId: contact.Id,
		Body:      "Test note",
	}

	_, err := client.ContactsCreateNote(note)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNotes(t *testing.T) {
	_, contact, client := GetTestingData()

	notes, err := client.ContactsGetAllNotes(contact.Id)
	if err != nil {
		t.Error(err)
	}
	if len(notes) == 0 {
		t.Error("no notes found")
	}
}

func TestDeleteNotes(t *testing.T) {
	_, contact, client := GetTestingData()

	note := runwayv2.Note{
		ContactId: contact.Id,
		Id:        "CGiRDTLuBnYbYM6oMRFc",
	}

	err := client.ContactsDeleteNote(note)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateNotes(t *testing.T) {
	_, contact, client := GetTestingData()

	note := runwayv2.Note{
		ContactId: contact.Id,
		Id:        "CGiRDTLuBnYbYM6oMRFc",
		Body:      "Updated note",
	}

	_, err := client.ContactsUpdateNote(note)
	if err != nil {
		t.Error(err)
	}
}

func GetTestingData() (models.Location, runwayv2.Contact, runwayv2.Client) {
	loc := models.Location{}
	loc.Get("EgSKJqMXN6w7t2svxveY") // Client Runway test test

	svc := runway.GetSvc()
	client, err := svc.NewClientFromId(loc.Id)

	if err != nil {
		panic(err)
	}

	contacts, err := client.ContactsFind("m.siddikov@gmail.com")
	if err != nil {
		panic(err)
	}
	if len(contacts) == 0 {
		panic("no contacts found")
	}
	return loc, contacts[0], client
}
