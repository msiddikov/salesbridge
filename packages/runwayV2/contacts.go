package runwayv2

import (
	"encoding/json"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func (a *Client) ContactsFind(query string) ([]Contact, error) {
	res := struct {
		Contacts []Contact
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/contacts/",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
			{
				Key:   "query",
				Value: query,
			},
		},
	}, &res)
	return res.Contacts, err
}

func (a *Client) ContactsGet(contactId string) (Contact, error) {
	res := struct {
		Contact Contact
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/contacts/" + contactId,
	}, &res)
	return res.Contact, err
}

func (a *Client) ContactsCreate(contact Contact) (Contact, error) {
	contact.LocationId = a.cfg.locationId
	res := struct {
		Contact Contact
	}{}

	body, err := json.Marshal(contact)
	if err != nil {
		return res.Contact, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/contacts/",
		Body:     string(body),
	}, &res)
	return res.Contact, err

}

// First searches by email, if not found searches by phone
// If not found neither by email nor by phone creates a new contact
func (a *Client) ContactsFirstOrCreate(contact Contact) (Contact, error) {
	contacts, err := a.ContactsFindByEmailPhone(contact.Email, contact.Phone)
	if err != nil {
		return Contact{}, err
	}

	// if found by email or phone, return the first one
	if len(contacts) != 0 {
		return contacts[0], nil
	}

	// if not found, create a new contact
	contact, err = a.ContactsCreate(contact)
	return contact, err
}

// Returns all contacts with the given email or phone
func (a *Client) ContactsFindByEmailPhone(email, phone string) ([]Contact, error) {
	contacts := []Contact{}
	err := error(nil)

	// search by email
	if email != "" {
		contacts, err = a.ContactsFind(email)
		if err != nil {
			return []Contact{}, err
		}
	}

	// search by phone
	if phone != "" {
		conts, err := a.ContactsFind(phone)
		if err != nil {
			return contacts, err
		}

		// append to contacts if not already there
		for _, cont := range conts {
			isFound := false
			for _, c := range contacts {
				if c.Id == cont.Id {
					isFound = true
					break
				}
			}
			if !isFound {
				contacts = append(contacts, cont)
			}
		}
	}

	return contacts, nil
}

// Finds duplicates
func (a *Client) ContactsDuplicates(email, phone string) (Contact, error) {
	res := struct {
		Contact Contact
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/contacts/search/duplicate",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
			{
				Key:   "email",
				Value: email,
			},
			{
				Key:   "number",
				Value: phone,
			},
		},
	}, &res)
	return res.Contact, err
}

func (a *Client) ContactsAddToWorkflow(contactId, workflowId string, eventTime time.Time) error {
	bodyStruct := struct {
		EventStartTime string `json:"eventStartTime"`
	}{
		EventStartTime: eventTime.Format("2006-01-02T15:04:05-07:00"),
	}
	body, _ := json.Marshal(bodyStruct)
	_, _, err := a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/contacts/" + contactId + "/workflow/" + workflowId,
		Body:     string(body),
	}, nil)

	return err
}

func (a *Client) ContactsUpdate(contact Contact) (Contact, error) {
	contactId := contact.Id
	contact.Id = ""
	contact.LocationId = ""
	res := struct {
		Contact Contact
	}{}

	body, err := lvn.Marshal(contact)
	if err != nil {
		return res.Contact, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "PUT",
		Endpoint: "/contacts/" + contactId,
		Body:     string(body),
	}, &res)
	return res.Contact, err
}

func (c *Contact) UpdateCustomFields(customFields []CustomFieldValue) {
	for _, cf := range customFields {
		isFound := false
		for k, v := range c.CustomFields {
			if v.Id == cf.Id {
				c.CustomFields[k] = cf
				isFound = true
				break
			}
		}

		if !isFound {
			c.CustomFields = append(c.CustomFields, cf)
		}
	}
}

func (a *Client) ContactsGetAllNotes(contactId string) ([]Note, error) {
	res := struct {
		Notes []Note
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/contacts/" + contactId + "/notes",
	}, &res)
	return res.Notes, err
}

func (a *Client) ContactsCreateNote(note Note) (Note, error) {
	res := struct {
		Note Note
	}{}

	body, err := lvn.MarshalSelected(note, "body", "userId")

	if err != nil {
		return res.Note, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/contacts/" + note.ContactId + "/notes",
		Body:     string(body),
	}, &res)
	return res.Note, err
}

func (a *Client) ContactsUpdateNote(note Note) (Note, error) {
	res := struct {
		Note Note
	}{}

	body, err := lvn.MarshalSelected(note, "body", "userId")

	if err != nil {
		return res.Note, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "PUT",
		Endpoint: "/contacts/" + note.ContactId + "/notes/" + note.Id,
		Body:     string(body),
	}, &res)
	return res.Note, err
}

func (a *Client) ContactsDeleteNote(note Note) error {
	res := struct {
		Succeded bool
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "DELETE",
		Endpoint: "/contacts/" + note.ContactId + "/notes/" + note.Id,
	}, &res)
	return err
}
