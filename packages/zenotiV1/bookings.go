package zenotiv1

import (
	"encoding/json"
	"fmt"
)

type (
	BookingReqGuestsRoom struct {
		Id string `json:"id"`
	}
	BookingReqGuestsItems struct {
		Item      BookingReqItem       `json:"item"`
		Room      BookingReqGuestsRoom `json:"room"`
		Therapist BookingReqTherapist  `json:"therapist"`
	}
	BookingReqGuest struct {
		Id    string                  `json:"id"`
		Items []BookingReqGuestsItems `json:"items"`
	}
	BookingReqItem struct {
		Id string `json:"Id"`
	}
	BookingReqTherapist struct {
		Id string `json:"Id"`
	}

	BookingReq struct {
		CenterId               string            `json:"center_id,omitempty"`
		Date                   ZenotiDate        `json:"date,omitempty"`
		Guests                 []BookingReqGuest `json:"guests,omitempty"`
		IsOnlyCatalogEmployees bool              `json:"is_only_catalog_employees"`
	}
)

func (c *Client) BookingsCreate(req BookingReq) (Booking, error) {

	book := Booking{}
	req.CenterId = c.cfg.centerId

	body, err := json.Marshal(req)
	if err != nil {
		return book, err
	}

	_, respBody, err := c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/bookings",
		QParams: []queryParam{
			{
				Key:   "is_double_booking_enabled",
				Value: "true",
			},
		},
		Body: string(body),
	}, &book)

	if book.Error.StatusCode != 0 {
		return book, fmt.Errorf("error: %s", string(book.Error.Message))
	}
	fmt.Println(string(respBody))

	return book, err
}

func (c *Client) BookingsReserve(bookingId, slot string) (Reservation, error) {
	res := Reservation{}
	_, _, err := c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/bookings/" + bookingId + "/slots/reserve",
		Body:     fmt.Sprintf(`{"slot_time": "%s"}`, slot),
	}, &res)

	if res.Error.StatusCode != 0 {
		return res, fmt.Errorf("error: %s", string(res.Error.Message))
	}

	return res, err
}

func (c *Client) BookingsConfirm(bookingId string) (Reservation, error) {
	res := Reservation{}
	_, _, err := c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/bookings/" + bookingId + "/slots/confirm",
	}, &res)

	if res.Error.StatusCode != 0 {
		return res, fmt.Errorf("error: %s", string(res.Error.Message))
	}
	return res, err
}

func (c *Client) BookWithConfirm(req BookingReq) (Booking, error) {
	book, err := c.BookingsCreate(req)
	if err != nil {
		return book, err
	}

	_, err = c.BookingsReserve(book.Id, req.Date.Time.Format("2006-01-02T15:04:05"))
	if err != nil {
		return book, err
	}

	_, err = c.BookingsConfirm(book.Id)
	if err != nil {
		return book, err
	}

	return book, nil
}
