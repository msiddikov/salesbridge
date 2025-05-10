package webServer

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/types"
	"client-runaway-zenoti/internal/zenoti"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"errors"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func setSettingsRoutes(router *gin.Engine) {
	router.GET("/settings/locations", allLocations)
	router.GET("/settings/locations/:id", locationInfo)
	router.POST("/settings/locations/:id", saveLocation)
	router.GET("/settings/oauthLinks", getRunwayAuthLink)
	router.GET("/settings/workflows/:id", getWorkflows)
	router.GET("/settings/zenoti-centers/:api", getCenters)
	router.GET("/settings/zenoti-centers/:api/:centerid/services", getCenterServices)
}

func getWorkflows(c *gin.Context) {
	data, err := runway.GetWorkflows(c.Param("id"))

	if err != nil {
		panic(err)
	}

	c.Data(lvn.Res(200, data, "ok"))
}

func allLocations(c *gin.Context) {
	locations := []types.Location{}

	dbLocations := []models.Location{}
	db.DB.Find(&dbLocations)

	for _, l := range dbLocations {
		locations = append(locations, types.Location{
			Name:         l.Name,
			Id:           l.Id,
			IsIntegrated: true})
	}

	c.Data(lvn.Res(200, locations, "OK"))
}

func locationInfo(c *gin.Context) {
	id := c.Param("id")

	info, err := runway.LocationInfo(id)
	if err != nil {
		panic(err)
	}

	c.Data(lvn.Res(200, info, "OK"))
}

func getRunwayAuthLink(c *gin.Context) {
	res := struct {
		Url    string
		Update string
	}{
		Url:    runway.GetOauthLink("/auth/gohighlevel"),
		Update: runway.GetOauthLink("/auth/gohighlevel/update"),
	}
	c.Data(lvn.Res(200, res, "OK"))
}

func saveLocation(c *gin.Context) {
	isNew := false
	loc := models.Location{}
	locDb := models.Location{}

	c.BindJSON(&loc)
	loc.Id = c.Param("id")

	err := db.DB.Where("id =?", c.Param("id")).First(&locDb).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	lvn.Logger.Noticef("Creating location %s", loc.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isNew = true

		// fill locDb with data from loc
		locDb.Id = loc.Id
		locDb.Name = loc.Name
	}

	locDb.PipelineId = loc.PipelineId
	locDb.BookId = loc.BookId
	locDb.SalesId = loc.SalesId
	locDb.NoShowsId = loc.NoShowsId
	locDb.ShowNoSaleId = loc.ShowNoSaleId
	locDb.MemberId = loc.MemberId
	locDb.TrackNewLeads = loc.TrackNewLeads
	if loc.ZenotiApi != "***" {
		locDb.ZenotiApi = loc.ZenotiApi
	}
	locDb.ZenotiUrl = loc.ZenotiUrl
	locDb.ZenotiCenterId = loc.ZenotiCenterId
	locDb.ZenotiCenterName = loc.ZenotiCenterName
	locDb.ZenotiServiceId = loc.ZenotiServiceId
	locDb.ZenotiServiceName = loc.ZenotiServiceName
	locDb.ZenotiServicePrice = loc.ZenotiServicePrice
	locDb.SyncCalendars = loc.SyncCalendars
	locDb.SyncContacts = loc.SyncContacts
	locDb.AutoCreateContacts = loc.AutoCreateContacts

	err = db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&locDb).Error
	lvn.Logger.Noticef("Created location %s", locDb.Name)
	if err != nil {
		panic(err)
	}

	if isNew {
		// save tokens
		client, err := runway.FindRecentToken(c.Param("id"))
		if err != nil {
			panic(err)
		}

		err = locDb.SaveGhlTokens(client.GetTokens())
		lvn.GinErr(c, 500, err, "Error saving tokens")
	}
}

func getCenters(c *gin.Context) {
	api := c.Param("api")

	res := struct {
		Centers []zenotiv1.Center
		Url     string
	}{
		Url: zenoti.GetUsedUrlFromAPI(api),
	}

	centers, err := zenotiv1.CentersListAll(api)
	lvn.GinErr(c, 500, err, "Error getting centers")
	res.Centers = centers

	c.Data(lvn.Res(200, res, "OK"))
}

func getCenterServices(c *gin.Context) {
	api := c.Param("api")
	centerId := c.Param("centerid")

	cli, err := zenotiv1.NewClient("", centerId, api)
	lvn.GinErr(c, 500, err, "Error creating client")
	services, err := cli.CenterServicesGetAll(zenotiv1.CenterServicesFilter{})

	c.Data(lvn.Res(200, services, "OK"))
}
