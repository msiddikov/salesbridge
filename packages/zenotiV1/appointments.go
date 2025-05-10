package zenotiv1

import (
	"strconv"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

type (
	AppointmentFilter struct {
		StartDate           time.Time `json:"start_date"`
		EndDate             time.Time `json:"end_date"`
		TherapistId         string    `json:"therapist_id"`
		IncludeNoShowCancel bool      `json:"include_no_show_cancel"`
	}
)

func (c *Client) AppointmentsListAppointments(filter AppointmentFilter) ([]Appointment, error) {
	res := []Appointment{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/appointments",
		QParams: append(filter.GetQParams(),
			queryParam{
				Key:   "center_id",
				Value: c.cfg.centerId,
			},
		),
	}, &res)

	return res, err
}

func (c *Client) AppointmentsListAllAppointments(filter AppointmentFilter) ([]Appointment, error) {
	res := []Appointment{}
	from := filter.StartDate
	to := filter.EndDate
	curTo := from.Add(10 * 24 * time.Hour)
	curTo = lvn.Ternary(curTo.After(to), to, curTo)

	for curTo.After(from) {
		fltr := filter
		fltr.StartDate = from
		fltr.EndDate = curTo
		apts, err := c.AppointmentsListAppointments(fltr)

		if err != nil {
			return nil, err
		}

		res = append(res, apts...)
		from = curTo
		curTo = from.Add(10 * 24 * time.Hour)
		curTo = lvn.Ternary(curTo.After(to), to.Add(-1*time.Second), curTo)
	}

	return res, nil
}

func (c *Client) AppointmentsGetDetails(id string) ([]Appointment, error) {
	res := []Appointment{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/appointments/" + id,
	}, &res)
	return res, err
}

// Helper functions _________________________________________________________

func (f *AppointmentFilter) GetQParams() []queryParam {
	res := []queryParam{
		{
			Key:   "start_date",
			Value: f.StartDate.Format("2006-01-02"),
		},
		{
			Key:   "end_date",
			Value: f.EndDate.Format("2006-01-02"),
		},
	}
	if f.TherapistId != "" {
		res = append(res, queryParam{
			Key:   "therapist_id",
			Value: f.TherapistId,
		})
	}
	if f.IncludeNoShowCancel {
		res = append(res, queryParam{
			Key:   "include_no_show_cancel",
			Value: strconv.FormatBool(f.IncludeNoShowCancel),
		})
	}
	return res
}
