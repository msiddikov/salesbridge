package svc_cerbo

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetEncounterTypesForLocation(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locationID := c.Param("locationId")

	var location models.Location
	err := db.DB.Where("id = ? AND profile_id = ?", locationID, user.ProfileID).First(&location).Error
	lvn.GinErr(c, 400, err, "error while getting location")
	if err != nil {
		return
	}

	cli, err := clientForLocation(location)
	lvn.GinErr(c, 400, err, "error while creating cerbo client")
	if err != nil {
		return
	}

	types, err := cli.GetEncounterTypes()
	lvn.GinErr(c, 400, err, "error while getting encounter types")
	if err != nil {
		return
	}

	c.Data(lvn.Res(200, types, ""))
}
