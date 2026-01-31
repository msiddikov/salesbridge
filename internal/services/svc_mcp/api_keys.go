package svc_mcp

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// CreateAPIKeyRequest is the request body for creating a new MCP API key
type CreateAPIKeyRequest struct {
	Name        string   `json:"name" binding:"required"`
	LocationIDs []string `json:"locationIds" binding:"required"`
}

// CreateAPIKeyResponse includes the plain key (shown only once)
type CreateAPIKeyResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	KeyPrefix string `json:"keyPrefix"`
	PlainKey  string `json:"plainKey"` // Only returned on creation
}

// APIKeyListItem is a single item in the list response
type APIKeyListItem struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	KeyPrefix   string   `json:"keyPrefix"`
	IsActive    bool     `json:"isActive"`
	LocationIDs []string `json:"locationIds"`
}

// UpdateAPIKeyRequest is the request body for updating an API key
type UpdateAPIKeyRequest struct {
	Name        *string  `json:"name"`
	LocationIDs []string `json:"locationIds"`
}

// CreateAPIKey creates a new MCP API key with location access
func CreateAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req CreateAPIKeyRequest
	if err := c.BindJSON(&req); err != nil {
		lvn.GinErr(c, 400, err, "Invalid request body")
		return
	}

	// Generate API key
	plainKey, keyHash, keyPrefix, err := models.GenerateMCPApiKey()
	if err != nil {
		lvn.GinErr(c, 500, err, "Failed to generate API key")
		return
	}

	// Create key record
	apiKey := models.MCPApiKey{
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		ProfileID: user.ProfileID,
		IsActive:  true,
	}

	// Start transaction
	tx := db.DB.Begin()

	if err := tx.Create(&apiKey).Error; err != nil {
		tx.Rollback()
		lvn.GinErr(c, 500, err, "Failed to create API key")
		return
	}

	// Add allowed locations
	for _, locID := range req.LocationIDs {
		// Verify location belongs to user's profile
		var loc models.Location
		if err := tx.Where("id = ? AND profile_id = ?", locID, user.ProfileID).First(&loc).Error; err != nil {
			tx.Rollback()
			lvn.GinErr(c, 400, err, "Invalid location ID: "+locID)
			return
		}

		keyLoc := models.MCPApiKeyLocation{
			MCPApiKeyID: apiKey.ID,
			LocationID:  locID,
		}
		if err := tx.Create(&keyLoc).Error; err != nil {
			tx.Rollback()
			lvn.GinErr(c, 500, err, "Failed to associate location")
			return
		}
	}

	tx.Commit()

	// Return response with plain key (only time it's shown)
	c.Data(lvn.Res(200, CreateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		KeyPrefix: apiKey.KeyPrefix,
		PlainKey:  plainKey,
	}, "API key created. Save this key securely - it won't be shown again."))
}

// ListAPIKeys returns all MCP API keys for the user's profile
func ListAPIKeys(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var keys []models.MCPApiKey
	err := db.DB.
		Where("profile_id = ?", user.ProfileID).
		Preload("AllowedLocations").
		Find(&keys).Error

	if err != nil {
		lvn.GinErr(c, 500, err, "Failed to list API keys")
		return
	}

	// Map to response format
	response := make([]APIKeyListItem, len(keys))
	for i, key := range keys {
		response[i] = APIKeyListItem{
			ID:          key.ID,
			Name:        key.Name,
			KeyPrefix:   key.KeyPrefix,
			IsActive:    key.IsActive,
			LocationIDs: key.GetLocationIDs(),
		}
	}

	c.Data(lvn.Res(200, response, ""))
}

// RevokeAPIKey deactivates an MCP API key
func RevokeAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	keyID := c.Param("keyId")

	result := db.DB.
		Model(&models.MCPApiKey{}).
		Where("id = ? AND profile_id = ?", keyID, user.ProfileID).
		Update("is_active", false)

	if result.Error != nil {
		lvn.GinErr(c, 500, result.Error, "Failed to revoke API key")
		return
	}

	if result.RowsAffected == 0 {
		lvn.GinErr(c, 404, nil, "API key not found")
		return
	}

	c.Data(lvn.Res(200, nil, "API key revoked"))
}

// UpdateAPIKey updates the name and/or allowed locations for an API key
func UpdateAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	keyID := c.Param("keyId")

	var req UpdateAPIKeyRequest
	if err := c.BindJSON(&req); err != nil {
		lvn.GinErr(c, 400, err, "Invalid request body")
		return
	}

	// Verify key belongs to user's profile
	var apiKey models.MCPApiKey
	if err := db.DB.Where("id = ? AND profile_id = ?", keyID, user.ProfileID).First(&apiKey).Error; err != nil {
		lvn.GinErr(c, 404, err, "API key not found")
		return
	}

	// Start transaction
	tx := db.DB.Begin()

	// Update name if provided
	if req.Name != nil {
		apiKey.Name = *req.Name
		if err := tx.Save(&apiKey).Error; err != nil {
			tx.Rollback()
			lvn.GinErr(c, 500, err, "Failed to update API key")
			return
		}
	}

	// Update locations if provided
	if req.LocationIDs != nil {
		// Delete existing location associations
		if err := tx.Where("mcp_api_key_id = ?", apiKey.ID).Delete(&models.MCPApiKeyLocation{}).Error; err != nil {
			tx.Rollback()
			lvn.GinErr(c, 500, err, "Failed to update locations")
			return
		}

		// Add new location associations
		for _, locID := range req.LocationIDs {
			// Verify location belongs to user's profile
			var loc models.Location
			if err := tx.Where("id = ? AND profile_id = ?", locID, user.ProfileID).First(&loc).Error; err != nil {
				tx.Rollback()
				lvn.GinErr(c, 400, err, "Invalid location ID: "+locID)
				return
			}

			keyLoc := models.MCPApiKeyLocation{
				MCPApiKeyID: apiKey.ID,
				LocationID:  locID,
			}
			if err := tx.Create(&keyLoc).Error; err != nil {
				tx.Rollback()
				lvn.GinErr(c, 500, err, "Failed to associate location")
				return
			}
		}
	}

	tx.Commit()

	c.Data(lvn.Res(200, nil, "API key updated"))
}

// RegenerateAPIKey generates a new key for a revoked API key
func RegenerateAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	keyID := c.Param("keyId")

	// Find the key
	var apiKey models.MCPApiKey
	if err := db.DB.Where("id = ? AND profile_id = ?", keyID, user.ProfileID).First(&apiKey).Error; err != nil {
		lvn.GinErr(c, 404, err, "API key not found")
		return
	}

	// Only allow regeneration of revoked keys
	if apiKey.IsActive {
		lvn.GinErr(c, 400, nil, "Cannot regenerate an active key. Revoke it first.")
		return
	}

	// Generate new API key
	plainKey, keyHash, keyPrefix, err := models.GenerateMCPApiKey()
	if err != nil {
		lvn.GinErr(c, 500, err, "Failed to generate API key")
		return
	}

	// Update the key
	apiKey.KeyHash = keyHash
	apiKey.KeyPrefix = keyPrefix
	apiKey.IsActive = true

	if err := db.DB.Save(&apiKey).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to regenerate API key")
		return
	}

	c.Data(lvn.Res(200, CreateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		KeyPrefix: apiKey.KeyPrefix,
		PlainKey:  plainKey,
	}, "API key regenerated. Save this key securely - it won't be shown again."))
}

// DeleteAPIKey permanently deletes an MCP API key
func DeleteAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	keyID := c.Param("keyId")

	// Start transaction
	tx := db.DB.Begin()

	// Delete location associations first
	if err := tx.Where("mcp_api_key_id = ?", keyID).Delete(&models.MCPApiKeyLocation{}).Error; err != nil {
		tx.Rollback()
		lvn.GinErr(c, 500, err, "Failed to delete API key")
		return
	}

	// Delete the key
	result := tx.Where("id = ? AND profile_id = ?", keyID, user.ProfileID).Delete(&models.MCPApiKey{})
	if result.Error != nil {
		tx.Rollback()
		lvn.GinErr(c, 500, result.Error, "Failed to delete API key")
		return
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		lvn.GinErr(c, 404, nil, "API key not found")
		return
	}

	tx.Commit()

	c.Data(lvn.Res(200, nil, "API key deleted"))
}
