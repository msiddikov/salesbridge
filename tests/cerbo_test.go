package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/services/svc_cerbo"
	"client-runaway-zenoti/packages/cerbo"
	"fmt"
	"testing"
)

var (
	subdomain = "medmatrixemr"
	cerboUser = "pk_hello123"
	cerboPass = "sk_hWhAqg_HaCesGMt_VfrxSyOt_FaGX"
)

func TestGetPatient(t *testing.T) {
	patientId := "1884"

	cCli, err := cerbo.NewClient(subdomain, cerboUser, cerboPass)
	if err != nil {
		t.Fatalf("Failed to create Cerbo client: %v", err)
	}
	patient, err := cCli.GetPatient(patientId)
	if err != nil {
		t.Fatalf("Failed to get patient: %v", err)
	}
	t.Logf("Patient: %+v", patient)
}

func TestGetUser(t *testing.T) {
	userId := "39"

	cCli, err := cerbo.NewClient(subdomain, cerboUser, cerboPass)
	if err != nil {
		t.Fatalf("Failed to create Cerbo client: %v", err)
	}
	user, err := cCli.GetUser(userId)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	t.Logf("User: %+v", user)
}

func TestGetNoteTypes(t *testing.T) {
	cCli := getCerboClient("Med Matrix")

	noteTypes, err := cCli.FtnGetAllAvailableTypes()
	if err != nil {
		t.Fatalf("Failed to get note types: %v", err)
	}
	t.Logf("Note Types: %+v", noteTypes)

}

func TestGetFreeTextNote(t *testing.T) {
	cCli := getCerboClient("Med Matrix")

	noteContent, err := cCli.FtnGetFreeTextNote("1884", 2)
	if err != nil {
		t.Fatalf("Failed to get free text note: %v", err)
	}
	fmt.Printf("Note Content: %s", noteContent)
}

func TestUpdateFreeTextNote(t *testing.T) {
	cCli := getCerboClient("Med Matrix")

	err := cCli.FtnUpdateFreeTextNote("1884", 2, "Updated note content from test")
	if err != nil {
		t.Fatalf("Failed to update free text note: %v", err)
	}
}

func getCerboClient(locName string) cerbo.Client {
	l := models.Location{}

	err := db.DB.Where("name=?", locName).Preload("CerboApiObj").Find(&l).Error
	if err != nil {
		panic(err)
	}

	cCli, err := cerbo.NewClient(l.CerboApiObj.Subdomain, l.CerboApiObj.Username, l.CerboApiObj.ApiKey)
	if err != nil {
		panic(err)
	}

	return cCli

}

func TestUpdateSectionFreeTextNotes(t *testing.T) {
	locName := "Med Matrix"
	patientID := "1884"
	noteTypeID := uint(1)
	sectionHeader := "Test Section 333"
	sectionContent := "This is the content for the test freaet."

	l := models.Location{}

	err := db.DB.Where("name=?", locName).Preload("CerboApiObj").Find(&l).Error
	if err != nil {
		t.Fatalf("Failed to find location: %v", err)
	}

	err = svc_cerbo.UpsertSectionIntoFreeTextNote(l, patientID, noteTypeID, sectionHeader, sectionContent)
	if err != nil {
		t.Fatalf("Failed to upsert section into free text note: %v", err)
	}
}
