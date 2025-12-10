package googleads

import (
	"testing"
)

func TestAuthURLRequiresConfig(t *testing.T) {
	svc := Service{}
	if _, err := svc.AuthURL("state"); err == nil {
		t.Fatalf("expected error when config is missing")
	}
}

func TestAuthURLBuilds(t *testing.T) {
	svc := Service{
		ClientID:     "client",
		ClientSecret: "secret",
		RedirectURL:  "https://example.com/cb",
	}
	u, err := svc.AuthURL("abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == "" {
		t.Fatalf("expected auth url")
	}
}
