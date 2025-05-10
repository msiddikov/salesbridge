package cmn

import (
	"client-runaway-zenoti/internal/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	ReqParams struct {
		Platform string
		Method   string
		Endpoint string
		Body     string
		QParams  []QueryParams
		Api      string
	}

	QueryParams struct {
		Key   string
		Value string
	}
)

var (
	zenotiCoolTill = time.Now()
	zenotiHost     = "https://api.zenoti.com/v1"
	runawayHost    = "https://rest.gohighlevel.com/v1"
	runaway2Host   = "https://api.msgsndr.com"
)

func Req(p ReqParams) (http.Response, error) {
	if p.Platform != "Z" && p.Platform != "R" && p.Platform != "R2" { // Z- zenoti and R- runaway
		return http.Response{}, fmt.Errorf("%s No platform defined", p.Endpoint)
	}

	host := ""
	isZenoti := false
	//isRunaway = false
	switch p.Platform {
	case "Z":
		host = zenotiHost
		isZenoti = true
	case "R":
		host = runawayHost
	case "R2":
		host = runaway2Host
	}

	// Cool down Zenoti
	if isZenoti && time.Now().Before(zenotiCoolTill) {
		fmt.Printf("Cooling down for %v\n", zenotiCoolTill.Sub(time.Now()))
		time.Sleep(zenotiCoolTill.Sub(time.Now()))
	}

	url := host + p.Endpoint
	client := &http.Client{}
	bodyReader := strings.NewReader(p.Body)

	req, err := http.NewRequest(p.Method, url, bodyReader)
	if p.Method == "GET" {
		req, err = http.NewRequest(p.Method, url, nil)
	}
	if err != nil {
		return http.Response{}, err
	}

	q := req.URL.Query()
	for _, v := range p.QParams {
		q.Add(v.Key, v.Value)
	}
	req.URL.RawQuery = strings.Replace(q.Encode(), "%40", "@", -1)

	switch p.Platform {
	case "Z":
		req.Header.Add("Authorization", "apikey "+p.Api)
	case "R":
		req.Header.Add("Authorization", "Bearer "+p.Api)
	}

	if p.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	if isZenoti && string(res.Header.Get("RateLimit-Remaining")) == "1" {
		resetIn, _ := strconv.ParseFloat(string(res.Header.Get("RateLimit-Reset")), 32)
		zenotiCoolTill = time.Now().Add(time.Duration(resetIn+2) * time.Second)
	}

	if res.StatusCode == 429 {
		fmt.Println("cooling down runaway for 40s...")
		time.Sleep(40 * time.Second)
		return Req(p)
	}

	if res.StatusCode == 502 {
		fmt.Println("Got 502 retrying in 5s...")
		time.Sleep(5 * time.Second)
		return Req(p)
	}

	if res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return http.Response{}, err
		}
		return *res, fmt.Errorf("%s: HTTP error: %v %s ", p.Endpoint, res.StatusCode, string(body))
	}

	return *res, nil
}

func GetLog(apts []types.Appointment) []string {
	log := []string{}

	for _, s := range apts {
		log = append(log, s.Id)
	}

	return log
}

func NotifySlack(channel, msg string) {
	webhook := ""
	switch channel {
	case "info":
		webhook = "T011KG2JUS2/B04EN9MBSLA/yTDH5ZrwZ3ROtPhWaXA94Me2"
	case "urgent":
		webhook = "T011KG2JUS2/B04MQPJ4D1N/BZg8yFThqlndeXKMsf84sLbW"
	case "critical":
		webhook = "T011KG2JUS2/B04LY4UTUDT/LZ1UdMHqgt7NFYJpq4qiBCkk"
	default:
		webhook = "T011KG2JUS2/B04EN9MBSLA/yTDH5ZrwZ3ROtPhWaXA94Me2"
	}
	body := struct {
		Text string `json:"text"`
	}{
		Text: msg,
	}
	data, _ := json.Marshal(body)
	http.Post("https://hooks.slack.com/services/"+webhook, "applicaton/json", strings.NewReader(string(data)))
}
