package svc_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetZenotiCenters(c *gin.Context) {
	zenotiApiId := c.Param("zenotiApiId")

	zenotiApi := models.ZenotiApi{}
	err := db.DB.First(&zenotiApi, "id = ?", zenotiApiId).Error
	lvn.GinErr(c, 400, err, "error while getting zenoti api")

	centers, err := zenotiv1.CentersListAll(zenotiApi.ApiKey)
	lvn.GinErr(c, 400, err, "error while getting zenoti centers")

	c.Data(lvn.Res(200, centers, ""))
}
