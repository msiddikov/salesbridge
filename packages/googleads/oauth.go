package googleads

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/metadata"
)

// Connection holds Google Ads OAuth tokens and metadata. It mirrors persisted data without tying to DB models.
type Connection struct {
	ID          uint
	ProfileID   uint
	DisplayName string
	Email       string

	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
}

var (
	errMissingConfig = errors.New("googleads oauth: client id/secret and redirect url are required")
	scopeEmail       = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile     = "https://www.googleapis.com/auth/userinfo.profile"
)

func (s Service) oauthConfig() (*oauth2.Config, error) {
	if s.ClientID == "" || s.ClientSecret == "" || s.RedirectURL == "" {
		return nil, errMissingConfig
	}
	var scope []string
	if s.Scope != "" {
		scope = strings.Fields(s.Scope)
	}
	if len(scope) == 0 {
		scope = []string{scopeAdwords, scopeEmail, scopeProfile}
	} else {
		ensure := func(val string) {
			for _, s := range scope {
				if s == val {
					return
				}
			}
			scope = append(scope, val)
		}
		ensure(scopeAdwords)
	}
	return &oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		RedirectURL:  s.RedirectURL,
		Scopes:       scope,
		Endpoint:     google.Endpoint,
	}, nil
}

// AuthURL builds the OAuth consent URL. The caller supplies a state that encodes profile/user context.
func (s Service) AuthURL(state string) (string, error) {
	cfg, err := s.oauthConfig()
	if err != nil {
		return "", err
	}
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

// ExchangeCode trades an authorization code for tokens and persists them via SaveConnection.
// The caller provides base metadata (profile/customer/login IDs, display name).
func (s Service) ExchangeCode(ctx context.Context, code string, base Connection) (Connection, error) {
	cfg, err := s.oauthConfig()
	if err != nil {
		return Connection{}, err
	}
	if code == "" {
		return Connection{}, errors.New("googleads oauth: code is required")
	}
	if s.SaveConnection == nil {
		return Connection{}, errors.New("googleads oauth: SaveConnection callback is not set")
	}

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return Connection{}, fmt.Errorf("exchange code: %w", err)
	}

	base.AccessToken = tok.AccessToken
	base.RefreshToken = tok.RefreshToken
	base.TokenExpiry = tok.Expiry

	if email, name, _ := fetchUserInfo(ctx, tok); email != "" {
		base.Email = email
		if base.DisplayName == "" && name != "" {
			base.DisplayName = fmt.Sprintf("%s's google ads", name)
		}
	}

	return s.SaveConnection(base)
}

// fetchUserInfo tries to get the user's email and name from Google userinfo endpoint. Best-effort only.
func fetchUserInfo(ctx context.Context, tok *oauth2.Token) (string, string, error) {
	if tok == nil || tok.AccessToken == "" {
		return "", "", errors.New("token missing")
	}
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("userinfo status %d", resp.StatusCode)
	}
	var body struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", "", err
	}
	return body.Email, body.Name, nil
}

// RefreshIfNeeded refreshes tokens if expired/near expiry using the stored refresh token.
func (s Service) RefreshIfNeeded(ctx context.Context, conn Connection) (Connection, error) {
	if conn.RefreshToken == "" {
		return conn, errors.New("googleads oauth: refresh token missing")
	}
	if s.UpdateConnectionTokens == nil {
		return conn, errors.New("googleads oauth: UpdateConnectionTokens callback is not set")
	}

	leeway := s.tokenExpiryLeeway
	if leeway == 0 {
		leeway = 60 * time.Second
	}

	nowFn := s.now
	if nowFn == nil {
		nowFn = time.Now
	}

	if conn.AccessToken != "" && nowFn().Add(leeway).Before(conn.TokenExpiry) {
		return conn, nil
	}

	cfg, err := s.oauthConfig()
	if err != nil {
		return conn, err
	}

	src := cfg.TokenSource(ctx, &oauth2.Token{
		RefreshToken: conn.RefreshToken,
		Expiry:       conn.TokenExpiry,
		AccessToken:  conn.AccessToken,
	})

	tok, err := src.Token()
	if err != nil {
		return conn, fmt.Errorf("refresh token: %w", err)
	}

	conn.AccessToken = tok.AccessToken
	if tok.RefreshToken != "" {
		conn.RefreshToken = tok.RefreshToken
	}
	conn.TokenExpiry = tok.Expiry

	if err := s.UpdateConnectionTokens(conn); err != nil {
		return conn, fmt.Errorf("persist refreshed tokens: %w", err)
	}

	return conn, nil
}

// WithHeaders returns a context containing the required Google Ads headers.
// It refreshes tokens if needed before returning.
func (s Service) WithHeaders(ctx context.Context, conn Connection) (context.Context, Connection, error) {
	updatedConn, err := s.RefreshIfNeeded(ctx, conn)
	if err != nil {
		return ctx, conn, err
	}

	mdPairs := []string{
		"authorization", "Bearer " + updatedConn.AccessToken,
	}
	devToken := ""
	if devToken == "" {
		devToken = developerToken
	}
	if devToken != "" {
		mdPairs = append(mdPairs, "developer-token", devToken)
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs(mdPairs...)), updatedConn, nil
}
