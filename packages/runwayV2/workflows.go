package runwayv2

func (a *Client) WorkflowsGet() ([]Workflow, error) {
	res := struct {
		Workflows []Workflow
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/workflows/",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
		},
	}, &res)

	return res.Workflows, err
}
