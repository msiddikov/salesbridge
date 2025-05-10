package zenoti

import (
	cmn "client-runaway-zenoti/internal/common"
	"client-runaway-zenoti/internal/config"
	testdata "client-runaway-zenoti/internal/tests/data"
	"client-runaway-zenoti/internal/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func GetAppointments(l config.Location, start, end time.Time) ([]types.Appointment, error) {
	empty := []types.Appointment{}

	params := []cmn.QueryParams{
		{Key: "center_id", Value: l.Zenoti.CenterId},
		{Key: "start_date", Value: start.Format("2006-01-02")},
		{Key: "end_date", Value: end.Format("2006-01-02")},
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/appointments",
		Api:      l.Zenoti.Api,
		QParams:  params,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := []types.Appointment{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result, nil
}

func GetCenters(api string) ([]types.Center, error) {
	empty := []types.Center{}
	res, err := cmn.Req(cmn.ReqParams{
		Method:   "GET",
		Endpoint: "/centers",
		Api:      api,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Centers []types.Center
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Centers, nil
}

func GetSales(l config.Location, start, end time.Time) ([]types.Sale, error) {
	empty := []types.Sale{}

	params := []cmn.QueryParams{
		{Key: "center_id", Value: l.Zenoti.CenterId},
		{Key: "start_date", Value: start.Format("2006-01-02")},
		{Key: "end_date", Value: end.Format("2006-01-02")},
	}

	res, err := cmn.Req(cmn.ReqParams{
		Method:   "GET",
		Endpoint: "/sales/salesreport",
		Api:      l.Zenoti.Api,
		QParams:  params,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Center_sales_report []types.Sale
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Center_sales_report, nil
}

func fillOutSales(c []types.Collection, api string) []types.Collection {

	for k, v := range c {
		g, err := GetGuest(v.Guest_id, api)
		if err != nil {
			fmt.Println(err)
			continue
		}
		c[k].Guest = g

		c[k].Total_collection = 0
		for _, p := range v.Items {
			c[k].Total_collection += p.Final_sale_price
		}
	}
	return c
}

func GetCollections(l config.Location, start, end time.Time) ([]types.Collection, error) {
	if config.IsTesting() {
		return fillOutSales(testdata.Collections(), l.Zenoti.Api), nil
	}
	empty := []types.Collection{}

	params := []cmn.QueryParams{
		{Key: "start_date", Value: start.Format("2006-01-02")},
		{Key: "end_date", Value: end.Format("2006-01-02")},
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/Centers/" + l.Zenoti.CenterId + "/collections_report",
		Api:      l.Zenoti.Api,
		QParams:  params,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Collections_report []types.Collection
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return fillOutSales(result.Collections_report, l.Zenoti.Api), nil
}

func GetGuestAppointments(id, api string) ([]types.Appointment, error) {
	empty := []types.Appointment{}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/guests/" + id + "/appointments",
		Api:      api,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Appointments []types.Appointment
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Appointments, nil
}

func GetGuest(id, api string) (types.GuestWithId, error) {
	empty := types.GuestWithId{}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/guests/" + id,
		Api:      api,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := types.GuestWithId{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result, nil
}

func TrimPhoneNumber(p string) string {
	p = strings.Replace(p, "+1", "", -1)
	p = strings.Replace(p, "+", "", -1)
	p = strings.Replace(p, " ", "", -1)
	return p
}

func CreateGuest(info types.Personal_info, l config.Location) (types.GuestWithId, error) {
	info.Mobile_phone.Number = TrimPhoneNumber(info.Mobile_phone.Number)
	empty := types.GuestWithId{}

	foundGuest, err := SearchGuest(info.Mobile_phone.Number, info.Email, l)

	if err == nil {
		return foundGuest, nil
	}

	if err != nil && err.Error() != "Not found" {
		return empty, err
	}

	// Testing
	// info.First_name = "TEST " + info.First_name

	// foundGuest, err = SearchGuest(info.Mobile_phone.Number, info.Email, l)

	// if err == nil {
	// 	return foundGuest, nil
	// }
	// if err != nil && err.Error() != "Not found" {
	// 	return empty, err
	// }
	// end testing

	body := struct {
		Center_Id     string              `json:"center_id"`
		Personal_Info types.Personal_info `json:"personal_info"`
		Referral      types.ReferralInfo  `json:"referral"`
	}{
		Center_Id:     l.Zenoti.CenterId,
		Personal_Info: info,
		Referral:      types.GetDefaultReferral(),
	}

	b, err := json.Marshal(body)
	if err != nil {
		return empty, err
	}
	bString := string(b)
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "POST",
		Endpoint: "/guests",
		Api:      l.Zenoti.Api,
		Body:     bString,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := types.GuestWithId{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result, nil
}

func SearchGuest(phone, email string, l config.Location) (types.GuestWithId, error) {

	g, err := searchGuestQuery("email", email, l)

	if err == nil {
		return g, nil
	}

	return searchGuestQuery("phone", phone, l)
}

func searchGuestQuery(query, value string, l config.Location) (types.GuestWithId, error) {
	empty := types.GuestWithId{}

	params := []cmn.QueryParams{
		{Key: "center_id", Value: l.Zenoti.CenterId},
		{Key: query, Value: value},
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/guests/search",
		Api:      l.Zenoti.Api,
		QParams:  params,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Guests []types.GuestWithId
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	if len(result.Guests) > 0 {
		return result.Guests[0], nil
	}

	return empty, fmt.Errorf("Not found")
}

func searchGuestQueryString(query, value string, l config.Location) (string, error) {

	params := []cmn.QueryParams{
		{Key: "center_id", Value: l.Zenoti.CenterId},
		{Key: query, Value: value},
		{Key: "expand", Value: "referral"},
		{Key: "expand", Value: "tags"},
		{Key: "expand", Value: "address_info"},
		{Key: "expand", Value: "preferences"},
		{Key: "expand", Value: "primary_employee"},
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/guests/search",
		Api:      l.Zenoti.Api,
		QParams:  params,
	})
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	result := struct {
		Guests []types.GuestWithId
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}

	if len(result.Guests) > 0 {
		return string(data), nil
	}

	return "", fmt.Errorf("Not found")
}

func UpdateGuest(email, phone string, l config.Location) error {
	guestInfo, err := searchGuestQueryString("email", email, l)
	if err != nil && err.Error() == "Not found" {
		guestInfo, err = searchGuestQueryString("phone", phone, l)
	}
	if err != nil && err.Error() == "Not found" {
		return fmt.Errorf("Not found")
	}
	if err != nil {
		return err
	}

	if strings.Contains(guestInfo, l.Zenoti.ReferralId) {
		return nil
	}

	id := gjson.Get(guestInfo, "guests.0.id").Str
	newBody, err := sjson.Set(gjson.Get(guestInfo, "guests.0").Raw, "referral", types.GetDefaultReferral())
	if err != nil {
		println(err)
		return err
	}

	_, err = cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "PUT",
		Endpoint: "/guests/" + id,
		Api:      l.Zenoti.Api,
		Body:     newBody,
	})
	if err != nil {
		println(err)
		return err
	}
	return nil
}

// get number of members for a location
func GetMembersNo(l config.Location) (int, error) {
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "Z",
		Method:   "GET",
		Endpoint: "/centers/" + l.Zenoti.CenterId + "/members",
		Api:      l.Zenoti.Api,
		QParams: []cmn.QueryParams{
			{Key: "status", Value: "Active"},
		},
	})
	if err != nil {
		return 0, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	result := struct {
		Page_info struct {
			Total int
		}
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return 0, err
	}

	return result.Page_info.Total, nil
}
