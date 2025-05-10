package meta

import (
	"net/url"
)

type (
	User struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	AdAccount struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	Paging struct {
		Limit    int        `json:"limit"`
		Type     PagingType `json:"paging_type"`
		Previous string     `json:"previous"`
		Next     string     `json:"next"`
	}

	MetaList[T any] struct {
		Data   []T    `json:"data"`
		Paging Paging `json:"paging"`
	}

	Insight map[string]string

	PagingType string
)

const (
	CursorBasedPagination PagingType = "cursor"
	TimeBasedPagination   PagingType = "time"
	OffsetBasedPagination PagingType = "offset"
)

func (p *Paging) getNextParams() ([]queryParam, error) {
	if p.Next == "" {
		return []queryParam{}, nil
	}
	res := []queryParam{}

	u, err := url.Parse(p.Next)
	if err != nil {
		return res, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return res, err
	}

	if q.Get("limit") != "" {
		res = append(res, queryParam{
			Key:   "limit",
			Value: q.Get("limit"),
		})
	}

	if q.Get("after") != "" {
		res = append(res, queryParam{
			Key:   "after",
			Value: q.Get("after"),
		})
	}

	if q.Get("until") != "" {
		res = append(res, queryParam{
			Key:   "until",
			Value: q.Get("until"),
		})
	}

	return res, nil
}

func (p *Paging) getPreviousParams() ([]queryParam, error) {
	if p.Previous == "" {
		return []queryParam{}, nil
	}
	res := []queryParam{}

	u, err := url.Parse(p.Next)
	if err != nil {
		return res, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return res, err
	}

	if q.Get("limit") != "" {
		res = append(res, queryParam{
			Key:   "limit",
			Value: q.Get("limit"),
		})
	}

	if q.Get("before") != "" {
		res = append(res, queryParam{
			Key:   "before",
			Value: q.Get("before"),
		})
	}

	if q.Get("since") != "" {
		res = append(res, queryParam{
			Key:   "since",
			Value: q.Get("since"),
		})
	}

	return res, nil
}
