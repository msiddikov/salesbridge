package openaiv1

import "testing"

func TestResponsesCreate(t *testing.T) {
	client := getClient(t)

	resp, err := client.ResponsesCreate(ResponseRequest{
		Model: "gpt-4.1-mini",
		Input: "Say hello in one sentence.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID == "" {
		t.Error("expected response id")
	}
}
