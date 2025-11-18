package svc_attribution

import (
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func CreateAttributionFlow(c *gin.Context) {
	payload := models.AttributionFlow{}
	user := c.MustGet("user").(models.User)
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "Unable to bind JSON")

	payload.ProfileId = user.ProfileID

	err = models.DB.Create(&payload).Error
	lvn.GinErr(c, 400, err, "Unable to create attribution flow")

	c.Data(lvn.Res(200, payload, ""))
}

func DeleteAttributionFlow(c *gin.Context) {
	flowId := c.Param("flowId")

	err := models.DB.Delete(&models.AttributionFlow{}, "id = ?", flowId).Error
	lvn.GinErr(c, 400, err, "Unable to delete attribution flow")

	c.Data(lvn.Res(200, "", ""))
}

func GetAttributionFlows(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	flows := []models.AttributionFlow{}

	err := models.DB.Where("profile_id = ?", user.ProfileID).Find(&flows).Error
	lvn.GinErr(c, 400, err, "Unable to get attribution flows")

	c.Data(lvn.Res(200, flows, ""))
}
