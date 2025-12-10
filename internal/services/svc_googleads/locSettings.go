package svc_googleads

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"errors"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SaveLocationSetting creates or updates a Google Ads location setting for the authenticated profile.
func SaveLocationSetting(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locID := c.Param("locationId")
	if locID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("location id required"), "invalid location id")
		return
	}

	var payload struct {
		ConnectionID     uint   `json:"connectionId" binding:"required"`
		CustomerID       string `json:"customerId"`
		CustomerName     string `json:"customerName"`
		ManagerID        string `json:"managerId"`
		ManagerName      string `json:"managerName"`
		ClientCustomerID string `json:"clientCustomerId"`
		ClientName       string `json:"clientName"`
	}
	err := c.ShouldBindJSON(&payload)
	lvn.GinErr(c, 400, err, "invalid payload")

	// Normalize: ensure mutually exclusive configs make sense.
	if payload.ManagerID != "" && payload.ClientCustomerID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("clientCustomerId required when managerId is set"), "invalid manager configuration")
		return
	}
	if payload.ManagerID == "" && payload.CustomerID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("customerId or (managerId + clientCustomerId) required"), "missing account ids")
		return
	}

	// Ensure connection belongs to the profile.
	conn := models.GoogleAdsConnection{}
	if err := db.DB.Where("profile_id = ? AND id = ?", user.ProfileID, payload.ConnectionID).First(&conn).Error; err != nil {
		lvn.GinErr(c, 400, err, "connection not found")
		return
	}
	connectionName := conn.DisplayName
	if connectionName == "" {
		connectionName = conn.Email
	}
	if connectionName == "" {
		connectionName = fmt.Sprintf("connection %d", conn.ID)
	}

	setting := models.GoogleAdsLocationSetting{
		LocationId:       locID,
		ProfileID:        user.ProfileID,
		ConnectionID:     payload.ConnectionID,
		ConnectionName:   connectionName,
		CustomerID:       payload.CustomerID,
		CustomerName:     payload.CustomerName,
		ManagerID:        payload.ManagerID,
		ManagerName:      payload.ManagerName,
		ClientCustomerID: payload.ClientCustomerID,
		ClientName:       payload.ClientName,
	}

	// Upsert by (locationId, profileId)
	err = db.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "location_id"}, {Name: "profile_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"connection_id", "connection_name", "customer_id", "customer_name", "manager_id", "manager_name", "client_customer_id", "client_name", "updated_at"}),
		},
	).Create(&setting).Error
	lvn.GinErr(c, 400, err, "unable to save setting")
	if err != nil {
		return
	}

	c.Data(lvn.Res(200, setting, "saved"))
}

// GetLocationSetting returns the Google Ads setting for a given location (if any).
func GetLocationSetting(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locID := c.Param("locationId")
	if locID == "" {
		lvn.GinErr(c, 400, fmt.Errorf("location id required"), "invalid location id")
		return
	}

	setting := models.GoogleAdsLocationSetting{}
	err := db.DB.Where("location_id = ? AND profile_id = ?", locID, user.ProfileID).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Data(lvn.Res(200, nil, "not found"))
			return
		}
		lvn.GinErr(c, 400, err, "unable to fetch setting")
		return
	}

	c.Data(lvn.Res(200, setting, "success"))
}
