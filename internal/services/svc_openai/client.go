package svc_openai

import (
	"client-runaway-zenoti/internal/config"
	openaiv1 "client-runaway-zenoti/packages/openaiV1"
	"fmt"
)

func openAIClient() (openaiv1.Client, error) {
	settings := config.Confs.Settings
	if settings.OpenAIAPIKey == "" {
		return openaiv1.Client{}, fmt.Errorf("ai api key is required")
	}

	service := openaiv1.Service{
		APIKey:       settings.OpenAIAPIKey,
		BaseURL:      settings.OpenAIBaseURL,
		Organization: settings.OpenAIOrganization,
		Project:      settings.OpenAIProject,
	}

	return service.NewClient("")
}
