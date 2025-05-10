package meta

import (
	"fmt"
	"strings"
	"time"
)

type (
	InsightsReq struct {
		AdAccountID   string
		From          time.Time
		To            time.Time
		TimeIncrement int
		Paging        Paging
		Fields        []string
	}
)

func (c *Client) Insights(req InsightsReq) (res MetaList[Insight], err error) {
	err = req.Validate()
	if err != nil {
		return
	}

	qParams := req.GetQParams()

	resp, err := fetch[[]Insight](reqParams{
		Method:   "GET",
		Endpoint: fmt.Sprintf("/%s/insights", req.AdAccountID),
		QParams:  qParams,
	}, c)

	if err != nil {
		return
	}

	res.Data = resp.Data
	return
}

func (r *InsightsReq) GetQParams() []queryParam {
	res := []queryParam{
		{
			Key: "time_range",
			Value: fmt.Sprintf(`{"since":"%s","until":"%s"}`,
				r.From.Format("2006-01-02"),
				r.To.Format("2006-01-02"),
			),
		},
	}

	if r.TimeIncrement > 0 {
		res = append(res, queryParam{
			Key:   "time_increment",
			Value: fmt.Sprintf("%d", r.TimeIncrement),
		})
	}

	if len(r.Fields) > 0 {
		res = append(res, queryParam{
			Key:   "fields",
			Value: strings.Join(r.Fields, ","),
		})
	}

	return res
}

func (r *InsightsReq) Validate() error {
	if r.AdAccountID == "" {
		return fmt.Errorf("adAccountID is required")
	}
	if r.From.IsZero() {
		return fmt.Errorf("from is required")
	}
	if r.To.IsZero() {
		return fmt.Errorf("to is required")
	}
	return nil
}
