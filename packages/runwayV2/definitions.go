package runwayv2

import (
	"encoding/json"
	"time"
)

type (
	Contact struct {
		Id          string `json:"id,omitempty"`
		Email       string `json:"email,omitempty"`
		Phone       string `json:"phone,omitempty"`
		FirstName   string `json:"firstName,omitempty"`
		LastName    string `json:"lastName,omitempty"`
		DateOfBirth string `json:"dateOfBirth,omitempty"`
		Gender      string `json:"gender,omitempty"`
		LocationId  string `json:"locationId,omitempty"`

		Address1    string             `json:"address1,omitempty"`
		City        string             `json:"city,omitempty"`
		State       string             `json:"state,omitempty"`
		PostalCode  string             `json:"postalCode,omitempty"`
		Country     string             `json:"country,omitempty"`
		DndSettings ContactDndSettings `json:"dndSettings,omitempty"`

		Tags         []string           `json:"tags,omitempty"`
		CustomFields []CustomFieldValue `json:"customFields,omitempty"`

		Source string `json:"source,omitempty"` // e.g. "sales-bridge", "zenoti"
	}

	ContactDndSettings struct {
		Call     *ContactDndSetting `json:"Call,omitempty"`
		Sms      *ContactDndSetting `json:"SMS,omitempty"`
		Email    *ContactDndSetting `json:"Email,omitempty"`
		Whatsapp *ContactDndSetting `json:"Whatsapp,omitempty"`
		Gmb      *ContactDndSetting `json:"GMB,omitempty"`
		Fb       *ContactDndSetting `json:"FB,omitempty"`
	}

	ContactDndSetting struct {
		Status ContactDndStatus `json:"status,omitempty"`
	}

	ContactDndStatus string

	CustomFieldValue struct {
		Id          string
		Name        string
		Field_value any
	}

	Meta struct {
		StartAfterId string
		StartAfter   int
		Total        int
		Page         int
	}

	Pagination struct {
		Limit int `json:"limit,omitempty"`
		Page  int `json:"currentPage,omitempty"`
		Total int `json:"total,omitempty"`
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
		Source          string
		Status          OpportunityStatus
		Contact         Contact
		PipelineStageId string
		PipelineId      string
		CreatedAt       time.Time
		AssignedTo      string
		Attributions    []OpportunityAttribution
	}

	OpportunityAttribution struct {
		UtmGclid         string `json:"utmGclid,omitempty"`
		UtmSessionSource string `json:"utmSessionSource,omitempty"`
		Fbp              string `json:"fbp,omitempty"`
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

	TimeSlots struct {
		Dates struct {
			Slots []string `json:"slots,omitempty"`
		} `json:"_dates_,omitempty"`
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

	TriggerSubscriptionTriggerData struct {
		Id        string                                 `json:"id,omitempty"`
		Key       string                                 `json:"key,omitempty"`
		Filters   []TriggerSubscriptionTriggerDataFilter `json:"filters,omitempty"`
		EventType string                                 `json:"eventType,omitempty"`
		TargetUrl string                                 `json:"targetUrl,omitempty"`
	}

	TriggerSubscriptionEvent             string
	TriggerSubscriptionTriggerDataFilter struct {
		Field string   `json:"field,omitempty"`
		Value []string `json:"value,omitempty"`
	}

	TriggerSubscriptionExtra struct {
		LocationId string `json:"locationId,omitempty"`
		WorkflowId string `json:"workflowId,omitempty"`
	}

	TriggersSubscriptionMeta struct {
		Key     string `json:"key,omitempty"`
		Version string `json:"version,omitempty"`
	}

	TriggerSubscriptionData struct {
		TriggerData TriggerSubscriptionTriggerData `json:"triggerData,omitempty"`
		Extras      TriggerSubscriptionExtra       `json:"extras,omitempty"`
		Meta        TriggersSubscriptionMeta       `json:"meta,omitempty"`
	}

	Filter struct {
		Group    string `json:"group,omitempty"`
		Field    string `json:"field,omitempty"`
		Operator string `json:"operator,omitempty"`
		// Value is used when a single value is provided.
		Value string `json:"-"`
		// ValueRange is used when a range/object is provided for the value.
		ValueRange map[string]string `json:"-"`
		Filters    []Filter          `json:"filters,omitempty"`
	}

	Message struct {
		Id                     string   `json:"id,omitempty"`
		AltId                  string   `json:"altId,omitempty"`
		MessageType            string   `json:"messageType,omitempty"`
		LocationId             string   `json:"locationId,omitempty"`
		ContactId              string   `json:"contactId,omitempty"`
		ConversationId         string   `json:"conversationId,omitempty"`
		DateAdded              string   `json:"dateAdded,omitempty"`
		Body                   string   `json:"body,omitempty"`
		Direction              string   `json:"direction,omitempty"`
		Status                 string   `json:"status,omitempty"`
		ContentType            string   `json:"contentType,omitempty"`
		Attachments            []string `json:"attachments,omitempty"`
		UserId                 string   `json:"userId,omitempty"`
		ConversationProviderId string   `json:"conversationProviderId,omitempty"`
		From                   string   `json:"from,omitempty"`
		To                     string   `json:"to,omitempty"`
	}

	Appointment struct {
		Id                string `json:"id"`
		Address           string `json:"address"`
		Title             string `json:"title"`
		CalendarId        string `json:"calendarId"`
		ContactId         string `json:"contactId"`
		GroupId           string `json:"groupId"`
		AppointmentStatus string `json:"appointmentStatus"`
		AssignedUserId    string `json:"assignedUserId"`
		Notes             string `json:"notes"`
		Source            string `json:"source"`
		StartTime         string `json:"startTime"`
		EndTime           string `json:"endTime"`
		DateAdded         string `json:"dateAdded"`
		DateUpdated       string `json:"dateUpdated"`
	}
)

// MarshalJSON implements custom JSON marshaling for Filter so that the
// "value" key contains either the string Value (when non-empty) or the
// object ValueRange (when Value is empty). Other zero-value fields are
// omitted like normal `omitempty` behaviour.
func (f Filter) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if f.Group != "" {
		m["group"] = f.Group
	}
	if f.Field != "" {
		m["field"] = f.Field
	}
	if f.Operator != "" {
		m["operator"] = f.Operator
	}
	// value: prefer Value if non-empty, otherwise ValueRange if present
	if f.Value != "" {
		m["value"] = f.Value
	} else if len(f.ValueRange) > 0 {
		m["value"] = f.ValueRange
	}
	if len(f.Filters) > 0 {
		m["filters"] = f.Filters
	}
	return json.Marshal(m)
}

// UnmarshalJSON supports reading a Filter where the "value" key may be a
// string or an object. It will populate Value or ValueRange accordingly.
func (f *Filter) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// helper to unmarshal string fields
	var s string
	if v, ok := raw["group"]; ok {
		if err := json.Unmarshal(v, &s); err == nil {
			f.Group = s
		}
	}
	if v, ok := raw["field"]; ok {
		if err := json.Unmarshal(v, &s); err == nil {
			f.Field = s
		}
	}
	if v, ok := raw["operator"]; ok {
		if err := json.Unmarshal(v, &s); err == nil {
			f.Operator = s
		}
	}

	// value can be either string or object
	if v, ok := raw["value"]; ok {
		// try string first
		var str string
		if err := json.Unmarshal(v, &str); err == nil {
			f.Value = str
		} else {
			// try object
			var obj map[string]string
			if err := json.Unmarshal(v, &obj); err == nil {
				f.ValueRange = obj
			} else {
				// unknown type: ignore
			}
		}
	}

	if v, ok := raw["filters"]; ok {
		var ff []Filter
		if err := json.Unmarshal(v, &ff); err == nil {
			f.Filters = ff
		}
	}

	return nil
}

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

const (
	ContactDndStatusActive    ContactDndStatus = "active"
	ContactDndStatusInactive  ContactDndStatus = "inactive"
	ContactDndStatusPermanent ContactDndStatus = "permanent"
)

const (
	TriggerSubscriptionEventCreate TriggerSubscriptionEvent = "CREATED"
	TriggerSubscriptionEventUpdate TriggerSubscriptionEvent = "UPDATED"
	TriggerSubscriptionEventDelete TriggerSubscriptionEvent = "DELETED"
)
