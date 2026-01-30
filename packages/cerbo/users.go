package cerbo

func (c *Client) GetUsers() (users []User, err error) {
	response := UsersListResponse{}
	_, _, err = c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/users",
		QParams: []queryParam{
			{
				Key:   "limit",
				Value: "100",
			},
		},
	}, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *Client) GetUser(userId string) (user User, err error) {

	_, _, err = c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/users/" + userId,
	}, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
