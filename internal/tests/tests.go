package tests

import (
	config "client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func CheckSales(api string) {

	centers, _ := zenoti.GetCenters(api)

	end := time.Now().Add(0 - 143*24*time.Hour)
	start := time.Now().Add(0 - 150*24*time.Hour)

	for end.Before(time.Now()) {

		fmt.Printf("Start: %s\n  End: %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
		for _, v := range centers {
			sales, _ := zenoti.GetSales(config.Location{Zenoti: struct {
				Api                 string
				CenterId            string
				Url                 string
				ReferralId          string
				Appointments        []string
				CollectionsSyncDay  time.Time
				AppointmentsSyncDay time.Time
			}{Api: api, CenterId: v.Id}}, start, end)
			fmt.Printf("Center %s %v sales\n", v.Id, len(sales))
		}
		start = start.Add(7 * 24 * time.Hour)
		end = end.Add(7 * 24 * time.Hour)
	}

}

func GetContact() {
	url := "https://rest.gohighlevel.com/v1/contacts/lookup?email=neanea126@yahoo.com&phone=%252B1+2076894988"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhdGlvbl9pZCI6InVpOU1LUmdnWEpNcmNVVk9pUmE4IiwiY29tcGFueV9pZCI6IlpGTnNvc25tWjNYZ3ZQYnF0YUt2IiwidmVyc2lvbiI6MSwiaWF0IjoxNjQzMDUzMTYzOTUzLCJzdWIiOiJ6YXBpZXIifQ.zFyKkORKLg50A41moAZh5bdH0FsWQy9bsGZupWTwRcs")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
