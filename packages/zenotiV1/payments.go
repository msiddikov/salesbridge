package zenotiv1

func (c *Client) GetPaymentOptions(guestId string) ([]PaymentOption, error) {

	res := struct {
		Accounts []PaymentOption
	}{}
	// Create a new HTTP GET request
	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/guests/" + guestId + "/accounts",
		QParams: []queryParam{
			{
				Key:   "center_id",
				Value: c.cfg.centerId,
			},
		},
	}, &res)
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}
