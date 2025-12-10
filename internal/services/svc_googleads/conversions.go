package svc_googleads

import (
	"client-runaway-zenoti/internal/db/models"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// GetLocationConversionActions lists conversion actions for the account configured on a location.
func GetLocationConversionActions(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locID := c.Param("locationId")
	if locID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("location id required"), "invalid location id")
	}

	cli, err := cliForLocation(locID, user.ProfileID)
	lvn.GinErr(c, 400, err, "unable to create google ads client")
	defer cli.Close()

	actions, err := cli.ListConversionActions()
	lvn.GinErr(c, 400, err, "unable to list conversion actions")

	c.Data(lvn.Res(200, actions, "success"))
}
