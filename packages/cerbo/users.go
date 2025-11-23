package cerbo

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
