package cerbo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	config struct {
		locationId string
		subdomain  string
		username   string
		apiKey     string
		created    time.Time
	}

	Client struct {
		cfg config
	}
)

type (
	reqParams struct {
		Method   string
		Endpoint string
		Body     string
		QParams  []queryParam
	}

	queryParam struct {
		Key   string
		Value string
	}
)

func NewClient(subdomain, username, apiKey string) (cli Client, err error) {

	cli = Client{
		cfg: config{
			subdomain: subdomain,
			username:  username,
			apiKey:    apiKey,
			created:   time.Now(),
		},
	}

	return
}

const (
	hostFormat = "https://%s.md-hq.com/api/v1"
)

func (c *Client) GetCfg() config {
	return c.cfg
}

var (
	coolDownTill = time.Now()
)

func (a *Client) fetch(r reqParams, data interface{}) (http.Response, []byte, error) {

	// Cool down
	if time.Now().Before(coolDownTill) {
		fmt.Printf("Cooling down for %v\n", time.Until(coolDownTill))
		time.Sleep(time.Until(coolDownTill))
	}

	url := fmt.Sprintf(hostFormat, a.cfg.subdomain) + r.Endpoint
	client := &http.Client{}
	bodyReader := strings.NewReader(r.Body)

	req, err := http.NewRequest(r.Method, url, bodyReader)
	if r.Method == "GET" {
		req, err = http.NewRequest(r.Method, url, nil)
	}

	if err != nil {
		return http.Response{}, []byte{}, err
	}

	q := req.URL.Query()
	for _, v := range r.QParams {
		q.Add(v.Key, v.Value)
	}
	req.URL.RawQuery = strings.Replace(q.Encode(), "%40", "@", -1)

	req.SetBasicAuth(a.cfg.username, a.cfg.apiKey)

	if r.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return http.Response{}, []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return http.Response{}, []byte{}, err
	}

	// Read rate limit
	if string(res.Header.Get("RateLimit-Remaining")) == "1" {
		resetIn, _ := strconv.ParseFloat(string(res.Header.Get("RateLimit-Reset")), 32)
		coolDownTill = time.Now().Add(time.Duration(resetIn+2) * time.Second)
	}

	if res.StatusCode == 429 {
		fmt.Println("cooling down zenoti for 40s...")
		time.Sleep(40 * time.Second)
		return a.fetch(r, data)
	}

	if res.StatusCode == 502 {
		fmt.Println("Got 502 retrying in 5s...")
		time.Sleep(5 * time.Second)
		return a.fetch(r, data)
	}

	if res.StatusCode > 299 {
		return *res, body, fmt.Errorf("CERBO>%s %s: HTTP error: %v %s ", r.Method, r.Endpoint, res.StatusCode, string(body))
	}

	if data == nil {
		return http.Response{}, body, nil
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return http.Response{}, body, err
	}

	return *res, body, nil
}
