package zenotiv1

import (
	"encoding/json"
	"strconv"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func (c *Client) ReportsCollections(from, to time.Time) ([]Collection, error) {
	res := struct {
		Collections_report []Collection
	}{}

	_, _, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/Centers/" + c.cfg.centerId + "/collections_report",
		QParams: []queryParam{
			{
				Key:   "start_date",
				Value: from.Format("2006-01-02"),
			},
			{
				Key:   "end_date",
				Value: to.Format("2006-01-02"),
			},
		},
	}, &res)
	return res.Collections_report, err
}

func (c *Client) ReportsAllCollections(from, to time.Time) ([]Collection, error) {
	res := []Collection{}

	curTo := from.Add(6 * 24 * time.Hour)
	curTo = lvn.Ternary(curTo.After(to), to, curTo)

	for curTo.After(from) {

		collections, err := c.ReportsCollections(from, curTo)

		if err != nil {
			return nil, err
		}

		res = append(res, collections...)
		from = curTo
		curTo = from.Add(6 * 24 * time.Hour)
		curTo = lvn.Ternary(curTo.After(to), to.Add(-1*time.Second), curTo)
	}

	return res, nil
}

//
// Sales Accrual Report
//

type SalesAccrualFilter struct {
	Start_date       ZenotiTime      `json:"start_date"`
	End_date         ZenotiTime      `json:"end_date"`
	Center_ids       []string        `json:"center_ids"`
	Invoice_statuses []InvoiceStatus `json:"invoice_statuses"` // 0-Open, 4 - Closed
	Item_types       []ItemType      `json:"item_types"`       // "Service", "Product", "Membership","GiftCard", "Package"
	Payment_types    []PaymentType   `json:"payment_types"`
	Sale_types       []SaleType      `json:"sale_types"`
}

type SalesDetails struct {
	Center_id           string
	Center_name         string
	Invoice_id          string
	Invoice_closed_date ZenotiTime
	Sale_date           ZenotiTime

	Collected                  float64
	Discount                   float64
	Redeemed                   float64
	Taxable_redemption         float64
	Sales_ex_tax               float64
	Sales_excluding_redemption float64
	Sales_inc_tax              float64
	Status                     string

	Item_id   string
	Item_name string
	Item_type ItemType
	Price     float64
	Qty       int

	Guest_id   string
	Guest_name string
}

func (c *Client) ReportsSalesAccrual(filter SalesAccrualFilter, pager PageInfo) ([]SalesDetails, PageInfo, error) {
	res := struct {
		Sales     []SalesDetails `json:"sales"`
		Page_info PageInfo       `json:"page_info"`
		Error     string         `json:"error"`
	}{}

	bodyBytes, err := json.Marshal(filter)
	if err != nil {
		return nil, PageInfo{}, err
	}

	_, _, err = c.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/reports/sales/accrual_basis/flat_file",
		QParams: []queryParam{
			{
				Key:   "page",
				Value: strconv.Itoa(pager.Page),
			},
			{
				Key:   "size",
				Value: strconv.Itoa(pager.Size),
			},
		},
		Body: string(bodyBytes),
	}, &res)
	return res.Sales, res.Page_info, err
}
