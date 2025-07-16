package runwayv2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const refreshTokenUrl = "https://api.msgsndr.com/oauth/token"
const host = "https://api.msgsndr.com"

type (
	OauthRes struct {
		Access_token  string
		Token_type    string
		Refresh_token string
		Scope         string
		UserType      string
		LocationId    string
	}

	reqParams struct {
		Method   string
		Endpoint string
		Body     string
		Version  string
		QParams  []queryParam
	}

	queryParam struct {
		Key   string
		Value string
	}
)

func (a *Client) AuthByAccessCode(code string) (err error) {
	a.cfg.accessCode = code
	return a.refreshAccessToken()
}

func (a *Client) refreshAccessToken() error {
	var grandType = ""
	if a.cfg.refreshToken != "" {
		grandType = "refresh_token"
	} else if a.cfg.accessCode != "" {
		grandType = "authorization_code"
	} else {
		err := fmt.Errorf("accessCode or refreshToken are needed")
		return err
	}

	hc := http.Client{}

	form := url.Values{}
	form.Add("client_id", a.service.ClientId)
	form.Add("client_secret", a.service.ClientSecret)
	form.Add("grant_type", grandType)
	form.Add("code", a.cfg.accessCode)
	form.Add("refresh_token", a.cfg.refreshToken)
	req, err := http.NewRequest("POST", refreshTokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := hc.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		fmt.Println(string(data))
		return fmt.Errorf(string(data))
	}

	resp := OauthRes{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	a.cfg.accessToken = resp.Access_token
	a.cfg.refreshToken = resp.Refresh_token
	a.cfg.locationId = resp.LocationId
	fmt.Println(resp.Access_token)
	fmt.Println(resp.Scope)
	return nil
}

func (p *reqParams) getVersion() string {
	if p.Version == "" {
		return "2021-04-15"
	}
	return p.Version
}

func (a *Client) fetch(r reqParams, data interface{}) (http.Response, []byte, error) {

	url := host + r.Endpoint
	return a.fetchUrl(url, r, data)
}

func (a *Client) fetchUrl(url string, r reqParams, data interface{}) (http.Response, []byte, error) {

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

	req.Header.Add("Authorization", "Bearer "+a.cfg.accessToken)

	req.Header.Add("Version", r.getVersion())

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

	if res.StatusCode == 401 && (strings.Contains(string(body), "access token has expired") || strings.Contains(string(body), "access token is invalid") || strings.Contains(string(body), "Invalid JWT")) {
		// need update the token
		err := a.UpdateToken()
		if err != nil {
			return *res, body, fmt.Errorf("RUNWAY_V2>%s %s: Unable to refresh the access token for %s: %s ", r.Method, r.Endpoint, a.cfg.locationId, err.Error())
		}
		return a.fetch(r, data)
	}

	if res.StatusCode == 429 {
		fmt.Println("cooling down runawayV2 for 3s...")
		time.Sleep(3 * time.Second)
		return a.fetch(r, data)
	}

	if res.StatusCode == 502 {
		fmt.Println("Got 502 retrying in 5s...")
		time.Sleep(5 * time.Second)
		return a.fetch(r, data)
	}

	if res.StatusCode > 299 {
		return *res, body, fmt.Errorf("RUNWAY_V2>%s %s: HTTP error: %v %s ", r.Method, r.Endpoint, res.StatusCode, string(body))
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

func (a *Client) UpdateToken() error {
	fmt.Println("Updating token for " + a.cfg.locationId)
	err := a.refreshAccessToken()
	if err != nil {
		return err
	}
	return a.saveTokens()
}

func (a *Client) saveTokens() error {
	return a.service.SaveTokens(a.cfg.locationId, a.cfg.accessToken, a.cfg.refreshToken)
}
