package cerbo

import (
	"encoding/json"
	"fmt"
)

func (c *Client) FtnGetAllAvailableTypes() ([]PtNoteType, error) {
	wrapped := struct {
		Data map[string]PtNoteType
	}{}
	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/pt_notes",
	}, &wrapped)
	if err != nil {
		return nil, err
	}

	ptTypes := []PtNoteType{}

	for _, v := range wrapped.Data {
		ptTypes = append(ptTypes, v)
	}

	return ptTypes, nil
}

func (c *Client) FtnUpdateFreeTextNote(ptId string, noteTypeId uint, content string) error {

	body := map[string]interface{}{
		"note": content,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: fmt.Sprintf("/patients/%s/pt_notes/%d", ptId, noteTypeId),
		Body:     string(bodyBytes),
	}, nil)
	return err
}

func (c *Client) FtnGetFreeTextNote(ptId string, noteTypeId uint) (string, error) {
	wrapped := struct {
		Note string
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: fmt.Sprintf("/patients/%s/pt_notes/%d", ptId, noteTypeId),
	}, &wrapped)
	if err != nil {
		return "", err
	}
	return wrapped.Note, nil
}
