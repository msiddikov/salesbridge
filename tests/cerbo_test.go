package tests

import (
	"client-runaway-zenoti/packages/cerbo"
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
