package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"client-runaway-zenoti/internal/config"
)

// Client handles communication with Grafana Loki
type Client struct {
	httpClient *http.Client
	endpoint   string
	userID     string
	apiKey     string
	dbName     string
}

var defaultClient *Client

// Init initializes the default Grafana client with configuration
func Init() {
	cfg := config.Confs.Grafana
	defaultClient = NewClient(cfg.Endpoint, cfg.UserID, cfg.APIKey, config.Confs.DB.DbName)
}

// NewClient creates a new Grafana Loki client
func NewClient(endpoint, userID, apiKey, dbName string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		endpoint:   endpoint,
		userID:     userID,
		apiKey:     apiKey,
		dbName:     dbName,
	}
}

// SendLog sends a log payload to Grafana Loki
func (c *Client) SendLog(logs logsPayload) error {
	logsByte, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(logsByte))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.userID, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		return fmt.Errorf("failed to send log, status code: %d: %s", resp.StatusCode, body.String())
	}

	return nil
}

// SendLogMessage sends a structured log message to Grafana Loki
func (c *Client) SendLogMessage(lm LogMessage) error {
	logs := logsPayload{
		Streams: []logEntry{
			{
				Stream: map[string]string{
					"application":   "ghl-zenoti",
					"database_name": c.dbName,
					"channel":       lm.Channel,
					"locationName":  lm.LocationName,
					"locationId":    lm.LocationId,
				},
				Values: [][2]string{
					{strconv.FormatInt(time.Now().UnixNano(), 10), lm.Msg},
				},
			},
		},
	}
	return c.SendLog(logs)
}

// Notify sends a notification log message
func (c *Client) Notify(locName, locId, channel, msg string) error {
	return c.SendLogMessage(LogMessage{
		Msg:          msg,
		Channel:      channel,
		LocationName: locName,
		LocationId:   locId,
	})
}

// Package-level functions that use the default client

// SendLogMessage sends a log message using the default client
func SendLogMessage(lm LogMessage) error {
	if defaultClient == nil {
		return fmt.Errorf("grafana client not initialized, call Init() first")
	}
	return defaultClient.SendLogMessage(lm)
}

// Notify sends a notification using the default client
func Notify(locName, locId, channel, msg string) {
	if defaultClient == nil {
		log.Println("grafana: client not initialized, call Init() first")
		return
	}
	if err := defaultClient.Notify(locName, locId, channel, msg); err != nil {
		log.Printf("grafana: failed to send log message: %v", err)
	}
}
