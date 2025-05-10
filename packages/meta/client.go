package meta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const host = "https://graph.facebook.com/v17.0"

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
	resp[T any] struct {
		Req     *http.Request
		Res     *http.Response
		RawBody []byte
		Data    T
		Paging  Paging
	}

	SourceWithPaging[T any] struct {
		Data   T
		Paging Paging
	}
)

var (
	coolDownTill = time.Now()
)

func fetch[T any](r reqParams, c *Client) (resp[T], error) {

	// Cool down
	if time.Now().Before(coolDownTill) {
		fmt.Printf("Cooling down meta for %v\n", time.Until(coolDownTill))
		time.Sleep(time.Until(coolDownTill))
	}

	result := resp[T]{}
	// Create body
	if r.Method == "" {
		r.Method = "GET"
	}

	//New request
	r.QParams = append(r.QParams,
		queryParam{
			Key:   "access_token",
			Value: c.cfg.access_token,
		},
	)
	addFields[T](&r.QParams)

	url := host + r.Endpoint
	client := &http.Client{}
	bodyReader := strings.NewReader(r.Body)

	req, err := http.NewRequest(r.Method, url, bodyReader)
	if err != nil {
		return result, err
	}

	q := req.URL.Query()
	for _, v := range r.QParams {
		q.Add(v.Key, v.Value)
	}
	req.URL.RawQuery = strings.Replace(q.Encode(), "%40", "@", -1)

	if r.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	// Read rate limit
	if string(res.Header.Get("RateLimit-Remaining")) == "1" {
		resetIn, _ := strconv.ParseFloat(string(res.Header.Get("RateLimit-Reset")), 32)
		coolDownTill = time.Now().Add(time.Duration(resetIn+2) * time.Second)
	}

	if res.StatusCode == 429 {
		fmt.Println("cooling down meta for 40s...")
		time.Sleep(40 * time.Second)
		return fetch[T](r, c)
	}

	if res.StatusCode == 502 {
		fmt.Println("Got 502 retrying meta in 5s...")
		time.Sleep(5 * time.Second)
		return fetch[T](r, c)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	result.Req = req
	result.Res = res
	result.RawBody = body

	if res.StatusCode > 299 {
		return result, fmt.Errorf("META>%s %s: HTTP error: %v %s ", r.Method, r.Endpoint, res.StatusCode, string(body))
	}

	if len(body) == 0 {
		return result, nil
	}

	if hasPaging[T]() {
		err = json.Unmarshal(body, &result)
		result.Req = req
		result.Res = res
		result.RawBody = body
	} else {
		err = json.Unmarshal(body, &result.Data)
	}

	if err != nil {
		return result, err
	}

	return result, nil
}

func getFields[T any]() (fields string) {
	res := []string{}
	var oI T
	t := reflect.TypeOf(oI)
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map || t.Kind() == reflect.Chan || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	o := reflect.New(t).Interface()

	p, _ := json.Marshal(o)
	var m map[string]interface{}
	json.Unmarshal(p, &m)

	for k := range m {
		res = append(res, k)
	}

	return strings.Join(res, ",")
}

// adds fields param, generated from type, to query params if not exists
func addFields[T any](params *[]queryParam) {
	for _, p := range *params {
		if p.Key == "fields" {
			return
		}
	}
	*params = append(*params, queryParam{
		Key:   "fields",
		Value: getFields[T](),
	})
}

func hasPaging[T any]() bool {
	var t T
	switch any(t).(type) {
	case []AdAccount:
		return true
	case []Insight:
		return true
	default:
		return false
	}
}
