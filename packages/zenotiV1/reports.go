package zenotiv1

import (
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func (c *Client) ReportsCollections(from, to time.Time) ([]Collection, error) {
	res := struct {
		Collections_report []Collection
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/Centers/" + c.cfg.centerId + "/collections_report",
		QParams: []queryParam{
			{
				Key:   "start_date",
				Value: from.Format("2006-01-02"),
			},
			{
				Key:   "end_date",
				Value: to.Format("2006-01-02"),
			},
		},
	}, &res)
	return res.Collections_report, err
}

func (c *Client) ReportsAllCollections(from, to time.Time) ([]Collection, error) {
	res := []Collection{}

	curTo := from.Add(6 * 24 * time.Hour)
	curTo = lvn.Ternary(curTo.After(to), to, curTo)

	for curTo.After(from) {

		collections, err := c.ReportsCollections(from, curTo)

		if err != nil {
			return nil, err
		}

		res = append(res, collections...)
		from = curTo
		curTo = from.Add(6 * 24 * time.Hour)
		curTo = lvn.Ternary(curTo.After(to), to.Add(-1*time.Second), curTo)
	}

	return res, nil
}
