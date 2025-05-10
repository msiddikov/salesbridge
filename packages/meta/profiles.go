package meta

func (c *Client) Me() (User, error) {
	res, err := fetch[User](reqParams{
		Endpoint: "/me",
	}, c)
	if err != nil {
		return User{}, err
	}
	return res.Data, nil
}

func (c *Client) MyAdAccounts() ([]AdAccount, error) {
	res, err := fetch[[]AdAccount](reqParams{
		Endpoint: "/me/adaccounts",
	}, c)
	if err != nil {
		return []AdAccount{}, err
	}
	return res.Data, nil
}
