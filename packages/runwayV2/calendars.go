package runwayv2

import "encoding/json"

type (
	CalendarCreateReq struct {
		Name          string                 `json:"name"`
		Notifications []CalendarNotification `json:"notifications"`
		LocationId    string                 `json:"locationId"`
		TeamMembers   []CalendarTeamMember   `json:"teamMembers"`
		OpenHours     []CalendarHours        `json:"openHours"`
	}

	CalendarNotification struct {
		Type                string   `json:"type"`
		ShouldSendToContact struct{} `json:"shouldSendToContact"`
		ShouldSendToUser    struct{} `json:"shouldSendToUser"`
		SelectedUsers       string   `json:"selectedUsers"`
		TemplateId          string   `json:"templateId"`
	}
	CalendarTeamMember struct {
		UserId          string `json:"userId"`
		Priority        int    `json:"priority"`
		MeetingLocation string `json:"meetingLocation"`
	}
	CalendarHours struct {
		OpenHour    int `json:"openHour"`
		OpenMinute  int `json:"openMinute"`
		CloseHour   int `json:"closeHour"`
		CloseMinute int `json:"closeMinute"`
	}

	CalendarCreateBlockSlotReq struct {
		LocationId string `json:"locationId"`
		CalendarId string `json:"calendarId"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		Title      string `json:"title"`
	}
)

func (a *Client) CalendarsGet() ([]Calendar, error) {
	res := struct {
		Calendars []Calendar
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/calendars/",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
		},
	}, &res)
	return res.Calendars, err
}

func (a *Client) CalendarGetEvents(calendarId string) ([]Event, error) {
	res := struct {
		Events []Event
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/calendars/" + calendarId + "/events",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
		},
	}, &res)
	return res.Events, err
}

func (a *Client) CalendarCreate(req CalendarCreateReq) (Calendar, error) {
	res := struct {
		Calendar Calendar
	}{}

	bd, err := json.Marshal(req)

	if err != nil {
		return res.Calendar, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/calendars/",
		Body:     string(bd),
	}, &res)

	return res.Calendar, err
}

func (a *Client) CalendarsDelete(id string) error {
	_, _, err := a.fetch(reqParams{
		Method:   "DELETE",
		Endpoint: "/calendars/" + id,
	}, nil)
	return err
}

func (a *Client) CalendarCreateBlockSlot(req BlockSlot) (BlockSlot, error) {
	res := BlockSlot{}
	bd, err := json.Marshal(req)

	if err != nil {
		return res, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/calendars/events/block-slots",
		Body:     string(bd),
	}, &res)

	return res, err
}

func (a *Client) CalendarDeleteEvent(eventId string) error {
	_, _, err := a.fetch(reqParams{
		Method:   "DELETE",
		Endpoint: "/calendars/events/" + eventId,
	}, nil)

	return err
}

func (a *Client) CalendarEditBlockSlot(id string, req BlockSlot) (BlockSlot, error) {
	res := BlockSlot{}
	bd, err := json.Marshal(req)

	if err != nil {
		return res, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "PUT",
		Endpoint: "/calendars/events/block-slots/" + id,
		Body:     string(bd),
	}, &res)

	return res, err
}
