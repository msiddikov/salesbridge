package runwayv2

func (a *Client) LocationsInfo() (Location, error) {
	res := struct {
		Location Location
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/locations/" + a.cfg.locationId,
	}, &res)

	return res.Location, err
}
