package svc_config

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetZenotiApis(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	zenotiApis := []models.ZenotiApi{}

	// do not select api_key
	err := db.DB.Select("id, api_name").Where("profile_id = ?", user.ProfileID).Find(&zenotiApis).Error
	lvn.GinErr(c, 400, err, "error while getting zenoti apis")

	c.Data(lvn.Res(200, zenotiApis, ""))
}

func CreateZenotiApi(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	payload := models.ZenotiApi{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	payload.ProfileId = user.ProfileID

	err = db.DB.Create(&payload).Error
	lvn.GinErr(c, 400, err, "error while creating zenoti api")

	c.Data(lvn.Res(200, payload, ""))
}

func UpdateZenotiApi(c *gin.Context) {
	zenotiApiId := c.Param("zenotiApiId")
	payload := models.ZenotiApi{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	var zenotiApi models.ZenotiApi
	err = db.DB.First(&zenotiApi, "id = ?", zenotiApiId).Error
	lvn.GinErr(c, 400, err, "error while getting zenoti api")

	err = db.DB.Model(&zenotiApi).Updates(payload).Error
	lvn.GinErr(c, 400, err, "error while updating zenoti api")

	c.Data(lvn.Res(200, zenotiApi, ""))
}

func DeleteZenotiApi(c *gin.Context) {
	zenotiApiId := c.Param("zenotiApiId")

	err := db.DB.Delete(&models.ZenotiApi{}, zenotiApiId).Error
	lvn.GinErr(c, 400, err, "error while deleting zenoti api")

	c.Data(lvn.Res(200, "", ""))
}
