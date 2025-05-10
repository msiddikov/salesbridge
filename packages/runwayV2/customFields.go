package runwayv2

import lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"

func (a *Client) CustomFieldsGet() ([]CustomField, error) {
	res := struct {
		CustomFields []CustomField
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/locations/" + a.cfg.locationId + "/customFields",
	}, &res)
	return res.CustomFields, err
}

func (a *Client) CustomFieldsCreate(cf CustomField) (CustomField, error) {
	res := struct {
		CustomField CustomField
	}{}

	body, err := lvn.Marshal(cf)
	if err != nil {
		return res.CustomField, err
	}

	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/locations/" + a.cfg.locationId + "/customFields",
		Body:     string(body),
	}, &res)
	return res.CustomField, err
}

// Searches by Name and DataType, creates new custom field if not found
func (a *Client) CustomFieldsFirstOrCreate(customField CustomField) (CustomField, error) {
	fields, err := a.CustomFieldsGet()

	if err != nil {
		return CustomField{}, err
	}

	for _, f := range fields {
		if f.Name == customField.Name && f.DataType == customField.DataType {
			return f, nil
		}
	}

	return a.CustomFieldsCreate(customField)
}
