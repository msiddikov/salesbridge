package cerbo

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
