package runwayv2

func (a *Client) TriggerFire(url string, body string) (err error) {
	req := reqParams{
		Method:   "POST",
		Endpoint: url,
		Body:     body,
	}

	_, _, err = a.fetchUrl(url, req, nil)

	return err
}
