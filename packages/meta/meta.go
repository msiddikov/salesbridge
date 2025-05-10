package meta

import (
	"time"
)

type (
	Service struct {
	}

	config struct {
		access_token string
		created      time.Time
	}

	Client struct {
		cfg     config
		service Service
	}
)

func (s *Service) NewClient(access_token string) (cli Client, err error) {
	cli = Client{
		cfg: config{
			access_token: access_token,
			created:      time.Now(),
		},
		service: *s,
	}
	return
}

func (c *Client) GetCfg() config {
	return c.cfg
}
