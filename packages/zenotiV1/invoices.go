package zenotiv1

import "fmt"

func (c *Client) InvoicesGetDetails(invoiceId string) (Invoice, error) {
	res := Invoice{}

	_, body, err := c.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/invoices/" + invoiceId,
	}, &res)

	fmt.Println(string(body))
	return res, err
}
