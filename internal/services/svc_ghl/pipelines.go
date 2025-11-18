package svc_ghl

import (
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetPipelines(c *gin.Context) {
	locationId := c.Param("locationId")

	cli, err := svc.NewClientFromId(locationId)
	lvn.GinErr(c, 500, err, "Error creating client")

	pipelines, err := cli.OpportunitiesGetPipelines()
	lvn.GinErr(c, 500, err, "Error getting pipelines")

	c.Data(lvn.Res(200, pipelines, ""))
}
