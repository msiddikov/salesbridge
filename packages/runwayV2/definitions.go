package runwayv2

import "time"

type (
	Contact struct {
		Id         string `json:"id,omitempty"`
		Email      string `json:"email,omitempty"`
		Phone      string `json:"phone"`
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Gender     string `json:"gender,omitempty"`
		LocationId string `json:"locationId,omitempty"`

		Address1   string `json:"address1,omitempty"`
		City       string `json:"city,omitempty"`
		State      string `json:"state,omitempty"`
		PostalCode string `json:"postalCode,omitempty"`
		Country    string `json:"country,omitempty"`

		Tags         []string           `json:"tags,omitempty"`
		CustomFields []CustomFieldValue `json:"customFields,omitempty"`

		Source string `json:"source,omitempty"` // e.g. "sales-bridge", "zenoti"
	}

	CustomFieldValue struct {
		Id          string
		Name        string
		Field_value any
	}

	Meta struct {
		StartAfterId string
		StartAfter   int
		Total        int
	}

	Note struct {
		Id        string    `json:"id,omitempty"`
		Body      string    `json:"body"`
		ContactId string    `json:"contactId,omitempty"`
		UserId    string    `json:"userId,omitempty"`
		DateAdded time.Time `json:"dateAdded,omitempty"`
	}

	Pipeline struct {
		Id     string
		Name   string
		Stages []struct {
			Id   string
			Name string
		}
	}

	Opportunity struct {
		Id              string
		Name            string
		MonetaryValue   float64
		Status          OpportunityStatus
		Contact         Contact
		PipelineStageId string
		PipelineId      string
		CreatedAt       time.Time
		AssignedTo      string
	}

	OpportunityChangeParams struct {
		Title       string  `json:"title"`
		StageId     string  `json:"stageId"`
		Status      string  `json:"status"`
		MonetyValue float64 `json:"monetaryValue"`
	}

	Location struct {
		Id        string
		CompanyId string
		Name      string
		Address   string
		City      string
	}

	Conversation struct {
		Id         string
		LocationId string
		ContactId  string
	}

	Workflow struct {
		Id     string
		Name   string
		Status string
	}

	CustomField struct {
		Id       string `json:",omitempty"`
		Name     string
		DataType string
	}

	NotFoundErr error

	Calendar struct {
		Id         string
		Name       string
		LocationId string
		GroupId    string
		IsActive   bool
	}

	Event struct {
		Id                string
		Name              string
		LocationId        string
		CalendarId        string
		StartTime         time.Time
		EndTime           time.Time
		Title             string
		ContactId         string
		AppointmentStatus string
	}

	OpportunityStatus string

	BlockSlot struct {
		Id            string    `json:"id,omitempty"`
		LocationId    string    `json:"locationId,omitempty"`
		CalendarId    string    `json:"calendarId,omitempty"`
		StartTime     time.Time `json:"startTime,omitempty"`
		EndTime       time.Time `json:"endTime,omitempty"`
		Title         string    `json:"title,omitempty"`
		CalendarNotes string    `json:"calendarNotes,omitempty"`
	}
)

const (
	OpportunityStatusOpen      OpportunityStatus = "open"
	OpportunityStatusWon       OpportunityStatus = "won"
	OpportunityStatusLost      OpportunityStatus = "lost"
	OpportunityStatusAbandoned OpportunityStatus = "abandoned"
	OpportunityStatusAll       OpportunityStatus = "all"
)

func (m *Meta) isZero() bool {
	return *m == Meta{}
}
