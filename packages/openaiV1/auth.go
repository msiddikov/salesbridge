package openaiv1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

type (
	reqParams struct {
		Method   string
		Endpoint string
		Body     string
		QParams  []queryParam
		Headers  map[string]string
	}

	queryParam struct {
		Key   string
		Value string
	}
)

func (c *Client) fetch(r reqParams, data interface{}) (http.Response, []byte, error) {
	url := strings.TrimRight(c.cfg.baseURL, "/") + r.Endpoint

	headers := http.Header{}
	if r.Method != "GET" && r.Body != "" {
		headers.Add("Content-Type", "application/json")
	}
	if c.cfg.apiKey != "" {
		headers.Add("Authorization", "Bearer "+c.cfg.apiKey)
	}
	if c.cfg.organization != "" {
		headers.Add("OpenAI-Organization", c.cfg.organization)
	}
	if c.cfg.project != "" {
		headers.Add("OpenAI-Project", c.cfg.project)
	}
	for key, val := range r.Headers {
		headers.Add(key, val)
	}

	if len(r.QParams) > 0 {
		values := neturl.Values{}
		for _, v := range r.QParams {
			values.Add(v.Key, v.Value)
		}
		url = url + "?" + values.Encode()
	}

	client := &http.Client{}
	if c.cfg.clientTimeout > 0 {
		client.Timeout = c.cfg.clientTimeout
	}

	res, body, err := c.doRequestWithRetry(client, r.Method, url, r.Body, headers, 2)
	if err != nil {
		return http.Response{}, []byte{}, err
	}

	if res.StatusCode > 299 {
		return *res, body, fmt.Errorf("OPENAI>%s %s: HTTP error: %v %s", r.Method, r.Endpoint, res.StatusCode, string(body))
	}

	if data == nil {
		return *res, body, nil
	}
	if len(body) == 0 {
		return *res, body, nil
	}

	if err := json.Unmarshal(body, data); err != nil {
		return *res, body, err
	}

	return *res, body, nil
}

func (c *Client) doRequestWithRetry(client *http.Client, method, url, body string, headers http.Header, retries int) (*http.Response, []byte, error) {
	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		var bodyReader io.Reader
		if method != "GET" && body != "" {
			bodyReader = strings.NewReader(body)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, nil, err
		}
		req.Header = headers.Clone()

		res, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, nil, err
		}

		body, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			lastErr = err
			if attempt < retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return res, nil, err
		}

		if res.StatusCode == 429 || res.StatusCode == 500 || res.StatusCode == 502 || res.StatusCode == 503 || res.StatusCode == 504 {
			if attempt < retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
		}

		return res, body, nil
	}

	return nil, nil, lastErr
}
