package svc_config

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LocationItem represents a location with just ID and Name
type LocationItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LocationWithIntegrations represents a location with integration status
type LocationWithIntegrations struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HasZenoti        bool   `json:"hasZenoti"`
	ZenotiCenterName string `json:"zenotiCenterName,omitempty"`
	HasCerbo         bool   `json:"hasCerbo"`
	HasGoogleAds     bool   `json:"hasGoogleAds"`
}

// IntegrationInfo represents an integration without credentials
type IntegrationInfo struct {
	Type   string `json:"type"`
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// GetLocationsForProfile returns all locations for a profile with integration status
func GetLocationsForProfile(profileID uint) ([]LocationWithIntegrations, error) {
	var locations []models.Location
	err := db.DB.
		Where("profile_id = ?", profileID).
		Preload("ZenotiApiObj").
		Preload("CerboApiObj").
		Find(&locations).Error
	if err != nil {
		return nil, err
	}

	// Check for Google Ads connections
	var gaSettings []models.GoogleAdsLocationSetting
	db.DB.Where("profile_id = ?", profileID).Find(&gaSettings)
	gaByLocation := make(map[string]bool)
	for _, s := range gaSettings {
		gaByLocation[s.LocationId] = true
	}

	result := make([]LocationWithIntegrations, len(locations))
	for i, loc := range locations {
		result[i] = LocationWithIntegrations{
			ID:               loc.Id,
			Name:             loc.Name,
			HasZenoti:        loc.ZenotiCenterId != "",
			ZenotiCenterName: loc.ZenotiCenterName,
			HasCerbo:         loc.CerboApiObjId != nil && *loc.CerboApiObjId > 0,
			HasGoogleAds:     gaByLocation[loc.Id],
		}
	}
	return result, nil
}

// GetLocationForProfile returns a single location for a profile with preloaded relations
func GetLocationForProfile(profileID uint, locationID string) (*models.Location, error) {
	var location models.Location
	err := db.DB.
		Where("id = ? AND profile_id = ?", locationID, profileID).
		Preload("ZenotiApiObj").
		Preload("CerboApiObj").
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// GetIntegrationsForProfile returns all integrations for a profile
func GetIntegrationsForProfile(profileID uint) []IntegrationInfo {
	var integrations []IntegrationInfo

	// Zenoti APIs
	var zenotiApis []models.ZenotiApi
	db.DB.Select("id, api_name").Where("profile_id = ?", profileID).Find(&zenotiApis)
	for _, api := range zenotiApis {
		integrations = append(integrations, IntegrationInfo{
			Type:   "zenoti",
			ID:     api.ID,
			Name:   api.ApiName,
			Status: "configured",
		})
	}

	// Cerbo APIs
	var cerboApis []models.CerboApi
	db.DB.Select("id, api_name").Where("profile_id = ?", profileID).Find(&cerboApis)
	for _, api := range cerboApis {
		integrations = append(integrations, IntegrationInfo{
			Type:   "cerbo",
			ID:     api.ID,
			Name:   api.ApiName,
			Status: "configured",
		})
	}

	// Google Ads connections
	var gaConnections []models.GoogleAdsConnection
	db.DB.Where("profile_id = ?", profileID).Find(&gaConnections)
	for _, conn := range gaConnections {
		integrations = append(integrations, IntegrationInfo{
			Type:   "google_ads",
			ID:     conn.ID,
			Name:   conn.DisplayName,
			Status: "connected",
		})
	}

	return integrations
}

// ListLocations returns all locations available to the user's profile
func ListLocations(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var locations []models.Location
	err := db.DB.
		Select("id, name").
		Where("profile_id = ?", user.ProfileID).
		Find(&locations).Error

	if err != nil {
		lvn.GinErr(c, 500, err, "Failed to list locations")
		return
	}

	response := make([]LocationItem, len(locations))
	for i, loc := range locations {
		response[i] = LocationItem{
			ID:   loc.Id,
			Name: loc.Name,
		}
	}

	c.Data(lvn.Res(200, response, ""))
}

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

// DeleteLocation deletes a location and all connected objects
func DeleteLocation(c *gin.Context) {
	locationId := c.Param("locationId")
	user := c.MustGet("user").(models.User)

	// Verify location exists and belongs to user's profile
	var location models.Location
	err := db.DB.Where("id = ? AND profile_id = ?", locationId, user.ProfileID).First(&location).Error
	if err != nil {
		lvn.GinErr(c, 404, err, "Location not found")
		return
	}

	// Start transaction
	tx := db.DB.Begin()

	// Delete related records in order (respecting foreign key constraints)
	deletions := []struct {
		model interface{}
		query string
	}{
		// Automation-related (delete runs first due to FK)
		{&models.AutomationRunNode{}, "run_id IN (SELECT id FROM automation_runs WHERE location_id = ?)"},
		{&models.AutomationRun{}, "location_id = ?"},
		{&models.AutomationBatchRun{}, "location_id = ?"},
		{&models.Edge{}, "automation_id IN (SELECT id FROM automations WHERE location_id = ?)"},
		{&models.Node{}, "automation_id IN (SELECT id FROM automations WHERE location_id = ?)"},
		{&models.Automation{}, "location_id = ?"},

		// MCP API key locations
		{&models.MCPApiKeyLocation{}, "location_id = ?"},

		// Google Ads settings
		{&models.GoogleAdsLocationSetting{}, "location_id = ?"},

		// GHL related
		{&models.GhlTrigger{}, "location_id = ?"},
		{&models.GhlTokens{}, "location_id = ?"},

		// Calendars and block slots
		{&models.BlockSlot{}, "location_id = ?"},
		{&models.Calendar{}, "location_id = ?"},

		// Reports
		{&models.JpmReportInvoice{}, "location_id = ?"},
		{&models.JpmReportNewLead{}, "location_id = ?"},

		// Contacts and related
		{&models.Appointment{}, "contact_id IN (SELECT contact_id FROM contacts WHERE location_id = ?)"},
		{&models.Sale{}, "contact_id IN (SELECT contact_id FROM contacts WHERE location_id = ?)"},
		{&models.Contact{}, "location_id = ?"},

		// Expenses
		{&models.LocationExpense{}, "location_id = ?"},

		// Chats
		{&models.ChatMessage{}, "chat_id IN (SELECT id FROM chats WHERE location_id = ?)"},
		{&models.Chat{}, "location_id = ?"},
	}

	for _, d := range deletions {
		if err := tx.Where(d.query, locationId).Delete(d.model).Error; err != nil {
			tx.Rollback()
			lvn.GinErr(c, 500, err, fmt.Sprintf("Failed to delete related records: %v", err))
			return
		}
	}

	// Check if ZenotiApi is used by other locations before deleting
	if location.ZenotiApiObjId != nil && *location.ZenotiApiObjId > 0 {
		var count int64
		tx.Model(&models.Location{}).Where("zenoti_api_obj_id = ? AND id != ?", *location.ZenotiApiObjId, locationId).Count(&count)
		if count == 0 {
			if err := tx.Delete(&models.ZenotiApi{}, *location.ZenotiApiObjId).Error; err != nil {
				tx.Rollback()
				lvn.GinErr(c, 500, err, "Failed to delete Zenoti API")
				return
			}
		}
	}

	// Check if CerboApi is used by other locations before deleting
	if location.CerboApiObjId != nil && *location.CerboApiObjId > 0 {
		var count int64
		tx.Model(&models.Location{}).Where("cerbo_api_obj_id = ? AND id != ?", *location.CerboApiObjId, locationId).Count(&count)
		if count == 0 {
			if err := tx.Delete(&models.CerboApi{}, *location.CerboApiObjId).Error; err != nil {
				tx.Rollback()
				lvn.GinErr(c, 500, err, "Failed to delete Cerbo API")
				return
			}
		}
	}

	// Finally delete the location itself
	if err := tx.Delete(&location).Error; err != nil {
		tx.Rollback()
		lvn.GinErr(c, 500, err, "Failed to delete location")
		return
	}

	tx.Commit()

	c.Data(lvn.Res(200, nil, "Location and all related data deleted successfully"))
}

// DeleteLocationDryRun returns a summary of what would be deleted without actually deleting
func DeleteLocationDryRun(c *gin.Context) {
	locationId := c.Param("locationId")
	user := c.MustGet("user").(models.User)

	// Verify location exists and belongs to user's profile
	var location models.Location
	err := db.DB.Where("id = ? AND profile_id = ?", locationId, user.ProfileID).First(&location).Error
	if err != nil {
		lvn.GinErr(c, 404, err, "Location not found")
		return
	}

	summary := make(map[string]int64)

	// Count related records
	counts := []struct {
		name  string
		model interface{}
		query func(*gorm.DB) *gorm.DB
	}{
		{"automations", &models.Automation{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"automation_runs", &models.AutomationRun{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"automation_batch_runs", &models.AutomationBatchRun{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"mcp_api_key_locations", &models.MCPApiKeyLocation{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"google_ads_settings", &models.GoogleAdsLocationSetting{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"ghl_triggers", &models.GhlTrigger{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"calendars", &models.Calendar{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"block_slots", &models.BlockSlot{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"contacts", &models.Contact{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"chats", &models.Chat{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"jpm_report_leads", &models.JpmReportNewLead{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
		{"jpm_report_invoices", &models.JpmReportInvoice{}, func(d *gorm.DB) *gorm.DB { return d.Where("location_id = ?", locationId) }},
	}

	for _, cnt := range counts {
		var count int64
		cnt.query(db.DB.Model(cnt.model)).Count(&count)
		if count > 0 {
			summary[cnt.name] = count
		}
	}

	// Check API objects
	if location.ZenotiApiObjId != nil && *location.ZenotiApiObjId > 0 {
		var count int64
		db.DB.Model(&models.Location{}).Where("zenoti_api_obj_id = ? AND id != ?", *location.ZenotiApiObjId, locationId).Count(&count)
		if count == 0 {
			summary["zenoti_api"] = 1
		}
	}

	if location.CerboApiObjId != nil && *location.CerboApiObjId > 0 {
		var count int64
		db.DB.Model(&models.Location{}).Where("cerbo_api_obj_id = ? AND id != ?", *location.CerboApiObjId, locationId).Count(&count)
		if count == 0 {
			summary["cerbo_api"] = 1
		}
	}

	c.Data(lvn.Res(200, map[string]interface{}{
		"location_id":   locationId,
		"location_name": location.Name,
		"to_delete":     summary,
	}, "Dry run complete - no data was deleted"))
}
