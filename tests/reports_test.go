package tests

import (
	"client-runaway-zenoti/internal/runway"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

var (
	stats = []string{
		"Expenses",
		// "Sales",
		// "SalesNo",
		// "ROI",
		// "NewLeads",
		// "Bookings",
		// "NoShows",
		// "ShowNoSale",
		// "LeadsConv",
		// "BookingsConv",
		// "ShowRate",
		// "Rank",
		// "MembershipConv",
	}

	locIds = []string{
		"jEbtFz8wbaNI1WWDDFIb",
		// "TknDfZc7WZF0cq5YE0du",
	}

	tags = []string{
		// "Facial | Lead",
		// "Body Contouring | Lead",
		// "march members",
		// "Membership | Lead",
		// "Morpheus | Lead",
		// "Email | Lead",
		// "Email List ",
		// "Walk in | Lead",
	}

	from = time.Now().Add(-30 * 24 * time.Hour)
)

func TestReportStats(t *testing.T) {
	for _, stat := range stats {
		res, err := runway.GetStats(stat, from, time.Now(), locIds, tags)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s: %v\n", stat, res)
	}
}

func TestReportDetails(t *testing.T) {
	details, err := runway.GetOpportunitiesForLocations(from, time.Now(), locIds, tags, "")
	if err != nil {
		t.Error(err)
	}

	bytes, err := runway.GetDetails(details)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("reportDetails.xlsx", bytes, 0644)
	if err != nil {
		t.Error(err)
	}
}
