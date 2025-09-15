package meta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response struct based on Meta Ads Insights API
type AdSpendResponse struct {
	Data []struct {
		Spend string `json:"spend"`
	} `json:"data"`
}

// GetMetaAdSpend fetches ad spend for the given account and date range
func GetMetaAdSpend(accessToken, adAccountID, startDate, endDate string) (string, error) {
	// Build URL
	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/%s/insights?fields=spend&time_range={\"since\":\"%s\",\"until\":\"%s\"}&access_token=%s",
		adAccountID, startDate, endDate, accessToken,
	)

	// Send request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("error response: %s", string(body))
	}

	// Parse response
	var result AdSpendResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return "0", nil
	}

	return result.Data[0].Spend, nil
}
