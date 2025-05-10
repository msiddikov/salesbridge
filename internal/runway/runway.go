package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"errors"

	"gorm.io/gorm"
)

func TradeAccessCode(code string) (locationId string, err error) {
	client, _ := svc.NewClientFromId("")

	err = client.AuthByAccessCode(code)

	if err != nil {
		return
	}
	recentTokens = append(recentTokens, client)
	locationId = client.GetLocationId()
	return
}

func TradeAccessCodeAndSaveTokens(code string) (locationId string, err error) {
	client, _ := svc.NewClientFromId("")

	err = client.AuthByAccessCode(code)

	if err != nil {
		return
	}

	locationId = client.GetLocationId()
	loc := models.Location{}
	loc.Get(locationId)

	loc.SaveGhlTokens(client.GetTokens())
	return
}

func LocationInfo(locationId string) (models.LocationInfo, error) {
	info := models.LocationInfo{}
	loc := models.Location{}
	err := db.DB.Where("id=?", locationId).First(&loc).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return info, err
	}
	loc.ZenotiApi = "***"
	cli := runwayv2.Client{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cli, err = FindRecentToken(locationId)
		if err != nil {
			return info, err
		}
		locInfo, err := cli.LocationsInfo()
		if err != nil {
			return info, err
		}
		loc.Name = locInfo.Name
		loc.Id = locInfo.Id

	} else {
		cli, _ = svc.NewClientFromId(locationId)
	}

	info.Pipelines, err = cli.OpportunitiesGetPipelines()

	if err != nil {
		return info, err
	}

	info.Workflows, err = cli.WorkflowsGet()
	info.Location = loc
	return info, err
}

func GetOauthLink(endpoint string) string {
	return svc.GetOauthLink(endpoint)
}
