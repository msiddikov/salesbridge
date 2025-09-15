package meta

import "testing"

func TestGetAdSpend(t *testing.T) {
	client := getTestClient()

	spend, err := GetMetaAdSpend(client.cfg.access_token, "act_552470297607717", "2025-09-01", "2025-09-30")
	if err != nil {
		t.Fatalf("failed to get ad spend: %v", err)
	}

	t.Logf("Ad spend: %s", spend)
}
