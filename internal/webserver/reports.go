package webServer

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/reports"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/types"
	"fmt"
	"net/http"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func setReportRoutes(router *gin.Engine) {
	router.POST("/reports/stats/:reportType", stats)
	router.POST("/reports/getDetails", getDetails)
	router.POST("/reports/setExpense", setExpense)
	router.GET("/reports/settings", getSettings)
}

func stats(c *gin.Context) {
	body := struct {
		From      time.Time
		To        time.Time
		Locations []string
		Tags      []string
	}{}
	c.Bind(&body)

	res := struct{ Data interface{} }{}
	var err error
	res.Data, err = runway.GetStats(c.Param("reportType"), body.From, body.To, body.Locations, body.Tags)

	if err != nil {
		panic(err)
	}

	c.JSON(200, res)
}

func setExpense(c *gin.Context) {
	body := struct {
		Locations []string
		From      time.Time
		To        time.Time
		Total     float64
	}{}

	c.Bind(&body)

	reports.SetExpenses(body.From, body.To, body.Total, body.Locations)
}

func getSettings(c *gin.Context) {
	locations := []models.Location{}
	res := struct {
		Locations []types.Location
		Tags      []string
	}{
		Locations: []types.Location{},
	}

	err := db.DB.Find(&locations).Error
	if err != nil {
		panic(err)
	}

	for _, v := range locations {
		zenotiLocation := config.GetLocationByRWID(v.Id)
		res.Locations = append(res.Locations, types.Location{
			Name:         v.Name,
			Id:           v.Id,
			IsIntegrated: zenotiLocation.Name != "",
		})
	}

	res.Tags = runway.GetTags()

	c.Data(lvn.Res(200, res, "OK"))
}

func getDetails(c *gin.Context) {
	body := struct {
		From      time.Time
		To        time.Time
		Locations []string
		Tags      []string
	}{}
	c.Bind(&body)

	res, err := runway.GetOpportunitiesForLocations(body.From, body.To, body.Locations, body.Tags, "")
	lvn.GinErr(c, 500, err, "Unable to get opportunities for locations")

	bytes, err := runway.GetDetails(res)
	lvn.GinErr(c, 500, err, "Unable to build xls file")

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s (%s).xlsx", body.From.Format("2006-01-02"), body.To.Format("2006-01-02"), body.Locations[0]))
	c.Data(http.StatusOK, "application/octet-stream", bytes)
}
