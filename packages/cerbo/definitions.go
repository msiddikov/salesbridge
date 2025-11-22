package cerbo

import (
	"encoding/json"
	"time"
)

type (
	WebhookData struct {
		PracticeId string          `json:"practice_id"`
		EventType  string          `json:"event_type"`
		Path       string          `json:"path"`
		Data       json.RawMessage `json:"data"`
	}

	Schedule struct {
		Id                  string    `json:"id"`
		Start               CerboTime `json:"start"`
		End                 CerboTime `json:"end"`
		Title               string    `json:"title"`
		AppointmentStatus   string    `json:"appointment_status"`
		AppointmentLocation string    `json:"appointment_location"`
		AssignedProviders   []struct {
			Id    string `json:"id"`
			First string `json:"first"`
		} `json:"assigned_providers"`
		Patient Patient `json:"patient"`
	}

	Patient struct {
		Id        string `json:"id"`
		FirstName string `json:"first"`
		LastName  string `json:"last"`
		Dob       string `json:"dob"`
		Email     string `json:"email1"`
		Phone     string `json:"phone_mobile"`
		Provider  string `json:"provider"`
	}

	CerboTime struct {
		Time time.Time `json:"time,omitempty"`
	}
)

func (t *CerboTime) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		t.Time = time.Now()
		return
	}

	date, err := time.Parse(`"2006-01-02 15:04:05"`, string(b))
	if err == nil {
		t.Time = date
		return
	}

	date, err = time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		t.Time = date
		return
	}

	return err
}

func (t CerboTime) MarshalJSON() ([]byte, error) {
	parsed := t.Time.Format("\"2006-01-02 15:04:05\"")
	return []byte(parsed), nil
}
