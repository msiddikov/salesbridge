package runwayv2

import "strings"

func (a *Client) TriggerFire(url string, body string) (err error) {
	req := reqParams{
		Method:   "POST",
		Endpoint: url,
		Body:     body,
	}

	_, _, err = a.fetchUrl(url, req, nil)

	if strings.Contains(err.Error(), " is inactive.") {
		return nil // Ignore inactive triggers
	}

	return err
}
