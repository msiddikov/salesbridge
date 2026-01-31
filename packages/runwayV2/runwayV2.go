package runwayv2

import (
	"fmt"
	"time"
)

type (
	Service struct {
		SaveTokens   func(locationId, accessToken, refreshToken string) error
		GetTokens    func(locationId string) (accessToken, refreshToken string, err error)
		ServerDomain string
		ClientId     string
		ClientSecret string
		Scope        string
	}

	config struct {
		locationId   string
		accessToken  string
		refreshToken string
		accessCode   string
		created      time.Time
	}

	Client struct {
		cfg     config
		service Service
	}
)

func (s *Service) NewClientFromId(locationId string) (cli Client, err error) {

	cli = Client{
		cfg: config{
			locationId: locationId,
			created:    time.Now(),
		},
		service: *s,
	}

	if locationId != "" {
		cli.cfg.accessToken, cli.cfg.refreshToken, err = s.GetTokens(locationId)
	}

	return
}

func (s *Service) NewClient(locationId, accessToken, refreshToken string) (cli Client, err error) {

	if accessToken == "" || refreshToken == "" {
		accessToken, refreshToken, err = s.GetTokens(locationId)
		if err != nil {
			return
		}
	}

	cli = Client{
		cfg: config{
			locationId:   locationId,
			accessToken:  accessToken,
			refreshToken: refreshToken,
			created:      time.Now(),
		},
		service: *s,
	}

	return
}

func (c *Client) GetLocationId() string {
	return c.cfg.locationId
}

func (c *Client) GetTokens() (accessToken string, refreshToken string) {
	return c.cfg.accessToken, c.cfg.refreshToken
}

func (c *Client) GetCreated() time.Time {
	return c.cfg.created
}

func (s *Service) GetOauthLink(redirectEndPoint, state string) string {
	url := fmt.Sprintf(`https://marketplace.gohighlevel.com/oauth/chooselocation?response_type=code&redirect_uri=%s&client_id=%s&scope=%s&state=%s`,
		s.ServerDomain+redirectEndPoint,
		s.ClientId,
		s.Scope,
		state,
	)
	return url
}
