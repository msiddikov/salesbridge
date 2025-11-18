package svc_config

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func UpdateLocation(c *gin.Context) {
	locationId := c.Param("locationId")
	payload := models.Location{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	var location models.Location
	err = db.DB.First(&location, "id = ?", locationId).Error
	lvn.GinErr(c, 400, err, "error while getting location")

	// check if profile id matches
	user := c.MustGet("user").(models.User)
	if location.ProfileID != user.ProfileID {
		lvn.GinErr(c, 403, nil, "forbidden")
	}

	err = db.DB.Model(&location).Updates(payload).Error
	lvn.GinErr(c, 400, err, "error while updating location")

	c.Data(lvn.Res(200, location, ""))
}
