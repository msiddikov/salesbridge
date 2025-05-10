package zenotiv1

import (
	"time"
)

type (
	config struct {
		locationId string
		centerId   string
		apiKey     string
		created    time.Time
		host       string
	}

	Client struct {
		cfg     config
		service Service
	}
)

const (
	hostProd  = "https://api.zenoti.com/v1"
	hostStage = "https://api.zenotistage.com/v1"
)

var (
	// apis that should be called in staging environment
	stagingApis = []string{
		"d77e4caa974f404f800463494b46a6fd768b5874cb334a86851a1dc65760d645",
	}
)

func NewClient(locationId, centerId, apiKey string) (cli Client, err error) {

	cli = Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
			host:       hostProd,
		},
	}

	// if api is in staging environment
	for _, api := range stagingApis {
		if apiKey == api {
			cli.cfg.host = hostStage
		}
	}

	return
}

func (c *Client) GetCfg() config {
	return c.cfg
}
