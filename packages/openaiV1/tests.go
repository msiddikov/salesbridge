package openaiv1

import (
	"os"
	"testing"
)

func getClient(t *testing.T) Client {
	t.Helper()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}

	service := Service{
		APIKey:       apiKey,
		BaseURL:      os.Getenv("OPENAI_BASE_URL"),
		Project:      os.Getenv("OPENAI_PROJECT"),
		Organization: os.Getenv("OPENAI_ORG"),
	}

	cli, err := service.NewClient("")
	if err != nil {
		t.Fatalf("client init: %v", err)
	}

	return cli
}
