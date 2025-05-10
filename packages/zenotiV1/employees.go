package zenotiv1

import (
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func (c *Client) EmployeesListAllBlockOutTimes(start, end time.Time) ([]BlockOut, error) {
	res := struct {
		Block_out_times []BlockOut
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/centers/" + c.cfg.centerId + "/blockouttimes",
		QParams: []queryParam{
			{
				Key:   "start_date",
				Value: start.Format("2006-01-02"),
			},
			{
				Key:   "end_date",
				Value: end.Format("2006-01-02"),
			},
		},
	}, &res)

	return res.Block_out_times, err
}

func (c *Client) EmployeesListBlockOutTimes(employeeId string, start, end time.Time) ([]BlockOut, error) {
	res := struct {
		Block_out_times []BlockOut
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/employees/" + employeeId + "/blockouttimes",
		QParams: []queryParam{
			{
				Key:   "center_id",
				Value: c.cfg.centerId,
			},
			{
				Key:   "start_date",
				Value: start.Format("2006-01-02"),
			},
			{
				Key:   "end_date",
				Value: end.Format("2006-01-02"),
			},
		},
	}, &res)

	return res.Block_out_times, err
}

func (c *Client) EmployeesListBlockOutTimesAll(employeeId string, start, end time.Time) ([]BlockOut, error) {
	res := []BlockOut{}

	from := start
	to := end
	curTo := from.Add(10 * 24 * time.Hour)
	curTo = lvn.Ternary(curTo.After(to), to, curTo)

	for curTo.After(from) {
		boTimes, err := c.EmployeesListBlockOutTimes(employeeId, from, curTo)

		if err != nil {
			return nil, err
		}

		res = append(res, boTimes...)
		from = curTo
		curTo = from.Add(10 * 24 * time.Hour)
		curTo = lvn.Ternary(curTo.After(to), to.Add(-1*time.Second), curTo)
	}

	return res, nil
}
