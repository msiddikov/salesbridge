package zenotiv1

import "fmt"

type (
	CenterServicesFilter struct {
		OnlyAddOns     bool   `json:"only_add_ons"`
		CatalogEnabled bool   `json:"catalog_enabled"`
		CategoryId     string `json:"category_id"`
		Page           int    `json:"page"`
		Size           int    `json:"size"`
	}
)

func CentersListAll(api string) ([]Center, error) {

	client, err := NewClient("", "", api)

	if err != nil {
		return nil, err
	}

	res := struct {
		Centers []Center
	}{}

	_, _, err = client.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/centers",
	}, &res)

	return res.Centers, err
}

func (c *Client) CenterServicesGet(filter CenterServicesFilter) ([]Service, PageInfo, error) {

	res := struct {
		Services []Service
		PageInfo PageInfo `json:"page_info"`
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/Centers/" + c.cfg.centerId + "/services",
		QParams:  filter.GetQParams(),
	}, &res)
	return res.Services, res.PageInfo, err

}

// CenterServicesGetAll returns all services for a center ignores the page and size fields in the filter
func (c Client) CenterServicesGetAll(filter CenterServicesFilter) ([]Service, error) {

	result := []Service{}

	filter.Page = 1
	filter.Size = 100

	for {
		services, pageInfo, err := c.CenterServicesGet(filter)
		if err != nil {
			return nil, err
		}

		result = append(result, services...)
		if pageInfo.Page*pageInfo.Size >= pageInfo.Total {
			break
		}
		filter.Page++
	}

	return result, nil

}

func (c *Client) CenterRoomsGet() ([]Room, error) {

	res := struct {
		Rooms []Room
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/Centers/" + c.cfg.centerId + "/rooms",
	}, &res)

	return res.Rooms, err

}

//////////////////////////////////////////////////////////////////////////
//
// 	Helper functions _________________________________________________________
//
//////////////////////////////////////////////////////////////////////////

func (f CenterServicesFilter) GetQParams() []queryParam {
	res := []queryParam{
		queryParam{
			Key:   "only_add_ons",
			Value: fmt.Sprint(f.OnlyAddOns),
		},
		queryParam{
			Key:   "catalog_enabled",
			Value: fmt.Sprint(f.CatalogEnabled),
		},
		queryParam{
			Key:   "category_id",
			Value: f.CategoryId,
		},
		queryParam{
			Key:   "page",
			Value: fmt.Sprint(f.Page),
		},
		queryParam{
			Key:   "size",
			Value: fmt.Sprint(f.Size),
		},
	}
	return res
}
