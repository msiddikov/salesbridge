package openaiv1

import "testing"

func TestAssistantsCreate(t *testing.T) {
	client := getClient(t)

	assistant, err := client.AssistantsCreate(AssistantRequest{
		Model:        "gpt-4.1-mini",
		Name:         "SalesBridge Test Assistant",
		Instructions: "Be concise.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if assistant.ID == "" {
		t.Error("expected assistant id")
	}
}
