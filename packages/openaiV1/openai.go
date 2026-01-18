package openaiv1

import (
	"fmt"
	"time"
)

const defaultBaseURL = "https://api.openai.com/v1"

type (
	Service struct {
		APIKey        string
		BaseURL       string
		Organization  string
		Project       string
		ClientTimeout time.Duration
	}

	config struct {
		apiKey        string
		baseURL       string
		organization  string
		project       string
		created       time.Time
		clientTimeout time.Duration
	}

	Client struct {
		cfg     config
		service Service
	}
)

func NewClient(apiKey string) (Client, error) {
	svc := Service{}
	return svc.NewClient(apiKey)
}

func (s *Service) NewClient(apiKey string) (Client, error) {
	if apiKey == "" {
		apiKey = s.APIKey
	}
	if apiKey == "" {
		return Client{}, fmt.Errorf("openai api key is required")
	}

	baseURL := s.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	cli := Client{
		cfg: config{
			apiKey:        apiKey,
			baseURL:       baseURL,
			organization:  s.Organization,
			project:       s.Project,
			created:       time.Now(),
			clientTimeout: s.ClientTimeout,
		},
		service: *s,
	}

	return cli, nil
}

func (c *Client) GetCfg() config {
	return c.cfg
}

func (c *Client) GetCreated() time.Time {
	return c.cfg.created
}
