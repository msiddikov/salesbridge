package svc_ghl

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	svc runwayv2.Service

	scopes = []string{
		"locations.readonly",
		"opportunities.readonly",
		"opportunities.write",
		"contacts.readonly",
		"conversations/message.write",
		"contacts.write",
		"conversations.readonly",
		"conversations.write",
		"conversations/message.readonly",
		"workflows.readonly",
		"locations/customFields.readonly",
		"locations/customFields.write",
		"calendars.readonly",
		"calendars.write",
		"calendars/events.readonly",
		"calendars/events.write",
	}
	recentTokens []runwayv2.Client
)

func init() {
	svc = runwayv2.Service{
		SaveTokens:   saveToken,
		GetTokens:    getToken,
		ClientId:     "681efbdb042c3d8b62429177-mahw24ly",
		ClientSecret: "5166c393-c76d-46e9-b8d0-6d8140c28456",
		Scope:        strings.Join(scopes, " "),
		ServerDomain: config.Confs.Settings.SrvDomain,
	}
	models.Svc = &svc
}

func getToken(locationId string) (accessToken, refreshToken string, err error) {
	loc := models.Location{}
	err = loc.Get(locationId)
	if err != nil {
		return
	}

	return loc.GetGhlTokens()
}

func saveToken(locationId, accessToken, refreshToken string) error {
	loc := models.Location{}

	err := loc.Get(locationId)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return loc.SaveGhlTokens(accessToken, refreshToken)
}

func FindRecentToken(locationId string) (runwayv2.Client, error) {
	var client runwayv2.Client
	for _, v := range recentTokens {
		if v.GetLocationId() == locationId {
			client = v
		}
	}
	if client.GetLocationId() != locationId {
		return runwayv2.Client{}, fmt.Errorf("recent token not found")

	}
	return client, nil
}

func GetSvc() runwayv2.Service {
	return svc
}
