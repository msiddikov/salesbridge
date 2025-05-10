package meta

import (
	"fmt"
	"testing"
	"time"
)

func TestInsights(t *testing.T) {
	cli := getTestClient()
	req := InsightsReq{
		AdAccountID:   "act_574247221284392",
		From:          time.Now().AddDate(0, -12, 0),
		To:            time.Now(),
		TimeIncrement: 1,
		Fields:        []string{"impressions", "spend"},
		Paging: Paging{
			Limit: 10,
		},
	}

	res, err := cli.Insights(req)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
