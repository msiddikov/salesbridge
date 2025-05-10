package webServer

import (
	"client-runaway-zenoti/internal/runway"

	"github.com/Lavina-Tech-LLC/lavinagopackage/v2/conf"
	"github.com/gin-gonic/gin"
)

func setSurveyRoutes(router *gin.Engine) {
	router.POST("/survey/:locationId/:workflowId", surveyPost)
	router.StaticFS("/apps/survey", gin.Dir(conf.GetPath()+"internal/webserver/survey-app/build/", false))
	router.StaticFS("/apps/survey2", gin.Dir(conf.GetPath()+"survey/build/", false))
}

func surveyPost(c *gin.Context) {
	body := runway.SurveyForm{}
	c.BindJSON(&body)
	err := runway.SurveyPost(body, c.Param("locationId"), c.Param("workflowId"))
	if err != nil {
		panic(err)
	}
}
