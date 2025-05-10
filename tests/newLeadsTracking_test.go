package tests

import (
	"client-runaway-zenoti/internal/runway"
	"fmt"
	"testing"
)

func TestLastOpportunityDate(t *testing.T) {
	date, err := runway.LastNewLeadDate("jEbtFz8wbaNI1WWDDFIb")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(date)
}

func TestNewLeadsTracking(t *testing.T) {
	runway.CheckForNewLeads()
}
