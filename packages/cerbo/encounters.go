package cerbo

import "encoding/json"

func (c *Client) GetEncounterTypes() ([]EncounterType, error) {
	wrapped := struct {
		Data []EncounterType `json:"data"`
	}{}
	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/encounter_types",
	}, &wrapped)
	if err != nil {
		return nil, err
	}

	return wrapped.Data, nil
}

func (c *Client) CreateEncounter(req EncounterCreateRequest) (Encounter, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return Encounter{}, err
	}

	res := Encounter{}
	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/encounters",
		Body:     string(payload),
	}, &res)
	if err != nil {
		return Encounter{}, err
	}

	return res, nil
}
