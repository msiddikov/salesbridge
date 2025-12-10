package svc_googleads

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"fmt"
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// ListConnections returns all Google Ads OAuth connections for the authenticated user's profile.
func ListConnections(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	var conns []models.GoogleAdsConnection

	err := db.DB.Where("profile_id = ?", user.ProfileID).Find(&conns).Error
	lvn.GinErr(c, 400, err, "unable to list connections")
	if err != nil {
		return
	}

	resp := make([]gin.H, 0, len(conns))
	for _, cn := range conns {
		accountName := cn.DisplayName
		if accountName == "" {
			accountName = cn.Email
		}
		if accountName == "" {
			accountName = fmt.Sprintf("connection %d", cn.ID)
		}
		resp = append(resp, gin.H{
			"id":          cn.ID,
			"accountName": accountName,
			"email":       cn.Email,
			"connectedAt": cn.CreatedAt,
		})
	}

	c.Data(lvn.Res(200, resp, "OK"))
}

// DeleteConnection removes a Google Ads connection for the authenticated user's profile.
func DeleteConnection(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	id, err := strconv.ParseUint(c.Param("accountId"), 10, 64)
	lvn.GinErr(c, 400, err, "invalid account id")
	if err != nil {
		return
	}

	res := db.DB.Where("profile_id = ? AND id = ?", user.ProfileID, id).Delete(&models.GoogleAdsConnection{})
	if res.Error != nil {
		lvn.GinErr(c, 400, res.Error, "unable to delete connection")
		return
	}
	if res.RowsAffected == 0 {
		lvn.GinErr(c, 404, nil, "connection not found")
		return
	}

	c.Data(lvn.Res(200, gin.H{"deleted": id}, "OK"))
}
