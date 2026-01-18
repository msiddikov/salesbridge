package cerbo

import (
	"fmt"
)

func (c *Client) GetPatient(patientId string) (patient Patient, err error) {

	_, _, err = c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/patients/" + patientId,
	}, &patient)
	if err != nil {
		return Patient{}, err
	}
	return patient, nil
}

func (c *Client) FindPatients(params PatientSearchParams) ([]Patient, error) {
	wrapped := struct {
		Data []Patient `json:"data"`
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/patients/search",
		QParams:  params.toQueryParams(),
	}, &wrapped)
	if err != nil {
		return nil, err
	}

	return wrapped.Data, nil
}

func (p PatientSearchParams) toQueryParams() []queryParam {
	params := []queryParam{}

	if p.FirstName != "" {
		params = append(params, queryParam{Key: "first_name", Value: p.FirstName})
	}
	if p.LastName != "" {
		params = append(params, queryParam{Key: "last_name", Value: p.LastName})
	}
	if p.Email != "" {
		params = append(params, queryParam{Key: "email", Value: p.Email})
	}
	if p.Phone != "" {
		params = append(params, queryParam{Key: "phone", Value: p.Phone})
	}
	if p.Dob != "" {
		params = append(params, queryParam{Key: "dob", Value: p.Dob})
	}
	if p.ProviderId != "" {
		params = append(params, queryParam{Key: "provider_id", Value: p.ProviderId})
	}
	if p.Tag != "" {
		params = append(params, queryParam{Key: "tag", Value: p.Tag})
	}
	if p.Username != "" {
		params = append(params, queryParam{Key: "username", Value: p.Username})
	}
	if p.Limit > 0 {
		params = append(params, queryParam{Key: "limit", Value: fmt.Sprintf("%d", p.Limit)})
	}
	if p.Offset > 0 {
		params = append(params, queryParam{Key: "offset", Value: fmt.Sprintf("%d", p.Offset)})
	}

	return params
}
