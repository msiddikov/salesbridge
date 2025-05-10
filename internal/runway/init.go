package runway

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var recentTokens []runwayv2.Client

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
)

func init() {
	svc = runwayv2.Service{
		SaveTokens:   saveToken,
		GetTokens:    getToken,
		ClientId:     "633bc5167ea65f59c51c6ab2-l8u3eg1p",
		ClientSecret: "112e0c5b-760e-4715-92b6-699bdef2cd06",
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
