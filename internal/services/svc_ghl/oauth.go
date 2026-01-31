package svc_ghl

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"fmt"
	"net/http"
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func TradeAccessCodeAndSaveTokens(code string, profileId uint) (locationId string, err error) {
	client, _ := svc.NewClientFromId("")

	err = client.AuthByAccessCode(code)

	if err != nil {
		return
	}

	locationId = client.GetLocationId()
	loc := models.Location{}
	loc.Get(locationId)
	if loc.Id == "" || loc.Name == "" {
		locInfo, err := client.LocationsInfo()
		if err != nil {
			return locationId, err
		}
		loc.Name = locInfo.Name
		loc.Id = locInfo.Id
		loc.ProfileID = profileId

		err = db.DB.Create(&loc).Error
		if err != nil {
			return locationId, err
		}
	}

	err = loc.SaveGhlTokens(client.GetTokens())
	if err != nil {
		return locationId, err
	}
	return
}

func GetGhlOauthLink(c *gin.Context) {
	profileId := c.MustGet("user").(models.User).ProfileID
	oauthLink := svc.GetOauthLink("/auth/hl", fmt.Sprintf("%d", profileId))

	c.Data(lvn.Res(200, oauthLink, "OK"))
}

func AddLocationAuthHandler(c *gin.Context) {
	code := c.Query("code")

	profileIds := c.Query("state")

	profileId, err := strconv.Atoi(profileIds)
	if err != nil {
		panic(err)
	}

	locId, err := TradeAccessCodeAndSaveTokens(code, uint(profileId))
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusTemporaryRedirect, config.Confs.Settings.AppDomain+"/oauth/callback?locationId="+locId)
}
