package runwayv2

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	OpportunitiesFilter struct {
		StartDate  time.Time
		EndDate    time.Time
		Status     OpportunityStatus
		Limit      int
		PipelineId string
		StageId    string
		ContactId  string
		Order      string
		Query      string
	}

	OpportunitiesSearchFilter struct {
		LocationId      string
		ContactId       string
		PipelineId      string
		PipelineStageId string
	}

	OpportunityUpdateReq struct {
		PipelineId      string            `json:"pipelineId"`
		Name            string            `json:"name"`
		PipelineStageId string            `json:"pipelineStageId"`
		Status          OpportunityStatus `json:"status"`
		MonetaryValue   float64           `json:"monetaryValue"`
		AssignedTo      string            `json:"assignedTo"`
	}

	OpportunityCreateReq struct {
		ContactId       string            `json:"contactId"`
		LocationId      string            `json:"locationId"`
		PipelineId      string            `json:"pipelineId"`
		Name            string            `json:"name"`
		PipelineStageId string            `json:"pipelineStageId"`
		Status          OpportunityStatus `json:"status"`
		MonetaryValue   float64           `json:"monetaryValue"`
		AssignedTo      string            `json:"assignedTo"`
	}
)

func (a *Client) OpportunitiesGetPipelines() ([]Pipeline, error) {
	res := struct {
		Pipelines []Pipeline
	}{Pipelines: []Pipeline{}}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/opportunities/pipelines",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
		},
	}, &res)
	if err != nil {
		return []Pipeline{}, err
	}
	return res.Pipelines, nil
}

func (a *Client) OpportunitiesGetAll(filter OpportunitiesFilter) ([]Opportunity, error) {
	hasLimit := filter.Limit != 0
	if !hasLimit {
		filter.Limit = 100
	}
	res := []Opportunity{}
	meta := Meta{}
	fetchedAll := false

	for !fetchedAll {
		opps, meta1, err := a.getOpportunitiesByMeta(filter, meta)
		meta = meta1

		if err != nil {
			return res, err
		}

		res = append(res, opps...)

		fetchedAll = meta.StartAfter == 0
		if hasLimit {
			break
		}
	}

	return res, nil
}

func (a *Client) OpportunitiesGet(id string) (Opportunity, error) {
	res := struct {
		Opportunity Opportunity
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/opportunities/" + id,
	}, &res)
	if err != nil {
		return Opportunity{}, err
	}
	return res.Opportunity, nil
}

func (a *Client) OpportunitiesSearch(filter OpportunitiesSearchFilter) ([]Opportunity, Meta, error) {
	res := struct {
		Opportunities []Opportunity
		Meta          Meta
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/opportunities/search",
		QParams:  filter.GetQParams(),
	}, &res)

	if err != nil {
		return []Opportunity{}, Meta{}, err
	}

	return res.Opportunities, res.Meta, nil
}

// Finds all opportunities with the given email or phone
func (a *Client) OpportunitiesFindByEmailPhone(email, phone, pipelineId string) ([]Opportunity, error) {
	contacts, err := a.ContactsFindByEmailPhone(email, phone)
	if err != nil {
		return []Opportunity{}, err
	}

	var opportunities []Opportunity
	for _, contact := range contacts {
		ops, _, err := a.OpportunitiesSearch(OpportunitiesSearchFilter{
			ContactId:  contact.Id,
			PipelineId: pipelineId,
			LocationId: a.cfg.locationId,
		})
		if err != nil {
			return []Opportunity{}, err
		}
		opportunities = append(opportunities, ops...)
	}

	return opportunities, nil
}

func (a *Client) OpportunitiesCreate(req OpportunityCreateReq) (Opportunity, error) {
	res := struct {
		Opportunity Opportunity
	}{}

	body, err := json.Marshal(req)
	if err != nil {
		return Opportunity{}, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/opportunities/",
		Body:     string(body),
	}, &res)
	if err != nil {
		return Opportunity{}, err
	}
	return res.Opportunity, nil
}

func (a *Client) OpportunitiesFirstOrCreate(filter OpportunitiesSearchFilter, req OpportunityCreateReq) (Opportunity, error) {
	opps, _, err := a.OpportunitiesSearch(filter)
	if err != nil {
		return Opportunity{}, err
	}

	if len(opps) == 0 {
		return a.OpportunitiesCreate(req)
	}

	return opps[0], nil
}

func (a *Client) OpportunitiesUpdate(req OpportunityUpdateReq, id string) (Opportunity, error) {
	res := struct {
		Opportunity Opportunity
	}{}

	body, err := json.Marshal(req)
	if err != nil {
		return Opportunity{}, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "PUT",
		Endpoint: "/opportunities/" + id,
		Body:     string(body),
	}, &res)
	if err != nil {
		return Opportunity{}, err
	}
	return res.Opportunity, nil
}

func (a *Client) getOpportunitiesByMeta(filter OpportunitiesFilter, meta Meta) ([]Opportunity, Meta, error) {
	res := struct {
		Opportunities []Opportunity
		Meta          Meta
	}{}
	params := filter.GetQParams()

	if !meta.isZero() {
		params = append(params,
			queryParam{
				Key:   "startAfter",
				Value: fmt.Sprint(meta.StartAfter),
			},
			queryParam{
				Key:   "startAfterId",
				Value: meta.StartAfterId,
			},
		)
	}

	params = append(params, queryParam{
		Key:   "location_id",
		Value: a.cfg.locationId,
	})

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/opportunities/search",
		QParams:  params,
	}, &res)

	return res.Opportunities, res.Meta, err

}

// Helper functions _________________________________________________________

func (f *OpportunitiesFilter) GetQParams() []queryParam {
	params := []queryParam{}

	if !f.StartDate.IsZero() {
		params = append(params, queryParam{
			Key:   "date",
			Value: f.StartDate.Format("01-02-2006"),
		})
	}

	if !f.EndDate.IsZero() {
		params = append(params, queryParam{
			Key:   "endDate",
			Value: f.EndDate.Format("01-02-2006"),
		})
	}

	if f.Status != "" {
		params = append(params, queryParam{
			Key:   "status",
			Value: string(f.Status),
		})
	}

	if f.PipelineId != "" {
		params = append(params, queryParam{
			Key:   "pipeline_id",
			Value: f.PipelineId,
		})
	}

	if f.StageId != "" {
		params = append(params, queryParam{
			Key:   "pipeline_stage_id",
			Value: f.StageId,
		})
	}

	if f.ContactId != "" {
		params = append(params, queryParam{
			Key:   "contact_id",
			Value: f.ContactId,
		})
	}

	if f.Limit != 0 {
		params = append(params, queryParam{
			Key:   "limit",
			Value: fmt.Sprint(f.Limit),
		})
	}

	if f.Query != "" {
		params = append(params, queryParam{
			Key:   "q",
			Value: f.Query,
		})
	}

	return params
}
func (f *OpportunitiesSearchFilter) GetQParams() []queryParam {
	params := []queryParam{}

	if f.LocationId != "" {
		params = append(params, queryParam{
			Key:   "location_id",
			Value: f.LocationId,
		})
	}

	if f.ContactId != "" {
		params = append(params, queryParam{
			Key:   "contact_id",
			Value: f.ContactId,
		})
	}

	if f.PipelineId != "" {
		params = append(params, queryParam{
			Key:   "pipeline_id",
			Value: f.PipelineId,
		})
	}

	if f.PipelineStageId != "" {
		params = append(params, queryParam{
			Key:   "pipeline_stage_id",
			Value: f.PipelineStageId,
		})
	}

	return params
}
