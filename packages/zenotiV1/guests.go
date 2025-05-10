package zenotiv1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	GuestFilter struct {
		Page               int
		Size               int
		SearchAcrossCenter bool
		CenterId           string
		ZipCode            string
		UserName           string
		UserCode           string
		Phone              string
		Email              string
		FirstName          string
		LastName           string
		Tags               string
	}

	GuestAppointmentsFilter struct {
		GuestId   string
		Page      int
		Size      int
		StartDate time.Time
		EndDate   time.Time
	}
)

// Creates a new guest.
func (c *Client) GuestsCreate(guest Guest) (Guest, error) {

	guest.Center_id = c.cfg.centerId
	guest.Personal_info.Mobile_phone.Number = TrimPhoneNumber(guest.Personal_info.Mobile_phone.Number)
	guest.Personal_info.Mobile_phone.Country_code = 225

	res := Guest{}
	body, err := json.Marshal(guest)
	if err != nil {
		return res, err
	}

	// create the guest
	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/guests/",
		Body:     string(body),
	}, &res)

	return res, err
}

// Updates the guest referral info.
// If the guest is not found, it returns an error.
// If the guest is found, it updates the guest referral info and returns nil.
func (c *Client) GuestsUpdate(email, phone string) error {
	// search for the guest by email
	guests, err := c.GuestsGetByPhoneEmail(phone, email)
	if err != nil {
		return err
	}

	// if still not found, return an error
	if len(guests) == 0 {
		return fmt.Errorf("guest not found")
	}

	// update the guest referral info
	return c.GuestsUpdateById(guests[0].Id)
}

// Searches for the guest by email and phone.
// If the guest is not found, it returns a "guest not found" error.
func (c *Client) GuestsGetByPhoneEmail(phone, email string) ([]Guest, error) {
	guests := []Guest{}
	err := error(nil)

	phone = TrimPhoneNumber(phone)

	// search for the guest by email
	if email != "" {
		guests, _ = c.GuestSearch(GuestFilter{
			Email: email,
		})
		if err != nil {
			return []Guest{}, err
		}
	}

	// if not found, search by phone
	if len(guests) == 0 && phone != "" {
		guests, err = c.GuestSearch(GuestFilter{
			Phone: phone,
		})
		if err != nil {
			return []Guest{}, err
		}
	}

	// if still not found, return an error
	if len(guests) == 0 {
		return []Guest{}, fmt.Errorf("guest not found")
	}

	return guests, nil
}

// It takes the guest info in full and edits it as a json string.
func (c *Client) GuestsUpdateById(id string) error {
	// get the guest info
	_, body, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/guests/" + id,
	}, nil)
	if err != nil {
		return err
	}

	// edit the guest info
	newBody := string(body)
	if err != nil {
		return err
	}

	// update the guest info
	_, _, err = c.fetch(reqParams{
		Method:   "PUT",
		Endpoint: "/guests/" + id,
		Body:     newBody,
	}, nil)

	return err
}

func (c *Client) GuestSearch(filter GuestFilter) ([]Guest, error) {
	res := struct {
		Guests []Guest
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/guests/search",
		QParams:  filter.GetQParams(),
	}, &res)
	return res.Guests, err
}

func (c *Client) GuestsListAppointments(filter GuestAppointmentsFilter) ([]Appointment, PageInfo, error) {
	res := struct {
		Appointments []Appointment
		Page_info    PageInfo
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/guests/" + filter.GuestId + "/appointments",
		QParams:  filter.GetQParams(),
	}, &res)
	return res.Appointments, res.Page_info, err
}

func (c *Client) GuestsGetById(id string) (Guest, error) {
	res := Guest{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/guests/" + id,
	}, &res)

	return res, err
}

// Helper functions _________________________________________________________

func (f *GuestFilter) GetQParams() []queryParam {
	res := []queryParam{}
	if f.Page != 0 {
		res = append(res, queryParam{
			Key:   "page",
			Value: strconv.Itoa(f.Page),
		})
	}
	if f.Size != 0 {
		res = append(res, queryParam{
			Key:   "size",
			Value: strconv.Itoa(f.Size),
		})
	}
	if f.SearchAcrossCenter {
		res = append(res, queryParam{
			Key:   "search_across_center",
			Value: strconv.FormatBool(f.SearchAcrossCenter),
		})
	}
	if f.CenterId != "" {
		res = append(res, queryParam{
			Key:   "center_id",
			Value: f.CenterId,
		})
	}
	if f.ZipCode != "" {
		res = append(res, queryParam{
			Key:   "zip_code",
			Value: f.ZipCode,
		})
	}
	if f.UserName != "" {
		res = append(res, queryParam{
			Key:   "user_name",
			Value: f.UserName,
		})
	}
	if f.UserCode != "" {
		res = append(res, queryParam{
			Key:   "user_code",
			Value: f.UserCode,
		})
	}
	if f.Phone != "" {
		res = append(res, queryParam{
			Key:   "phone",
			Value: f.Phone,
		})
	}
	if f.Email != "" {
		res = append(res, queryParam{
			Key:   "email",
			Value: f.Email,
		})
	}
	if f.FirstName != "" {
		res = append(res, queryParam{
			Key:   "first_name",
			Value: f.FirstName,
		})
	}
	if f.LastName != "" {
		res = append(res, queryParam{
			Key:   "last_name",
			Value: f.LastName,
		})
	}
	if f.Tags != "" {
		res = append(res, queryParam{
			Key:   "tags",
			Value: f.Tags,
		})
	}
	return res
}

func (f *GuestAppointmentsFilter) GetQParams() []queryParam {
	res := []queryParam{
		{
			Key:   "guest_id",
			Value: f.GuestId,
		},
	}
	if f.Page != 0 {
		res = append(res, queryParam{
			Key:   "page",
			Value: strconv.Itoa(f.Page),
		})
	}

	if f.Size != 0 {
		res = append(res, queryParam{
			Key:   "size",
			Value: strconv.Itoa(f.Size),
		})
	}

	if !f.StartDate.IsZero() {
		res = append(res, queryParam{
			Key:   "start_date",
			Value: f.StartDate.Format("2006-01-02"),
		})
	}

	if !f.EndDate.IsZero() {
		res = append(res, queryParam{
			Key:   "end_date",
			Value: f.EndDate.Format("2006-01-02"),
		})
	}

	return res
}

func TrimPhoneNumber(p string) string {
	p = strings.Replace(p, "+1", "", -1)
	p = strings.Replace(p, "+", "", -1)
	p = strings.Replace(p, " ", "", -1)
	return p
}
