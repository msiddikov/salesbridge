package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ExchangeResponse is the JSON returned by Meta's OAuth token exchange endpoint.
type ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds (typically 5,184,000 ≈ 60 days)
}

// ExchangeForLongLivedToken exchanges a short-lived user token for a long-lived one.
// Docs: GET https://graph.facebook.com/v21.0/oauth/access_token
func ExchangeForLongLivedToken(ctx context.Context, appID, appSecret, shortLivedToken string) (ExchangeResponse, error) {
	if appID == "" || appSecret == "" || shortLivedToken == "" {
		return ExchangeResponse{}, errors.New("appID, appSecret, and shortLivedToken are required")
	}

	base := "https://graph.facebook.com/v21.0/oauth/access_token"

	q := url.Values{}
	q.Set("grant_type", "fb_exchange_token")
	q.Set("client_id", appID)
	q.Set("client_secret", appSecret)
	q.Set("fb_exchange_token", shortLivedToken)

	u, _ := url.Parse(base)
	u.RawQuery = q.Encode()

	// Robust HTTP client with timeouts
	httpClient := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        100,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var fbErr map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&fbErr)
		return ExchangeResponse{}, fmt.Errorf("meta error: status=%d body=%v", resp.StatusCode, fbErr)
	}

	var out ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ExchangeResponse{}, fmt.Errorf("decode response: %w", err)
	}
	if out.AccessToken == "" {
		return ExchangeResponse{}, errors.New("empty access_token in response")
	}
	return out, nil
}
