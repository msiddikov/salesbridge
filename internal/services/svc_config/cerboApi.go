package svc_config

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetCerboApis(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	cerboApis := []models.CerboApi{}

	err := db.DB.Select("id, api_name, subdomain, username").Where("profile_id = ?", user.ProfileID).Find(&cerboApis).Error
	lvn.GinErr(c, 400, err, "error while getting cerbo apis")

	c.Data(lvn.Res(200, cerboApis, ""))
}

func CreateCerboApi(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	payload := models.CerboApi{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	payload.ProfileId = user.ProfileID

	err = db.DB.Create(&payload).Error
	lvn.GinErr(c, 400, err, "error while creating cerbo api")

	c.Data(lvn.Res(200, payload, ""))
}

func UpdateCerboApi(c *gin.Context) {
	cerboApiId := c.Param("cerboApiId")
	payload := models.CerboApi{}
	err := c.BindJSON(&payload)
	lvn.GinErr(c, 400, err, "error while binding json")

	var cerboApi models.CerboApi
	err = db.DB.First(&cerboApi, "id = ?", cerboApiId).Error
	lvn.GinErr(c, 400, err, "error while getting cerbo api")

	err = db.DB.Model(&cerboApi).Updates(payload).Error
	lvn.GinErr(c, 400, err, "error while updating cerbo api")

	c.Data(lvn.Res(200, cerboApi, ""))
}

func DeleteCerboApi(c *gin.Context) {
	cerboApiId := c.Param("cerboApiId")

	err := db.DB.Delete(&models.CerboApi{}, cerboApiId).Error
	lvn.GinErr(c, 400, err, "error while deleting cerbo api")

	c.Data(lvn.Res(200, "", ""))
}
