package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// MCPApiKey represents an API key for MCP server authentication
type MCPApiKey struct {
	Name             string `json:"name"`
	PlainKey         string `json:"-"` // Stored for internal use; external keys are not persisted in plain text.
	KeyHash          string `gorm:"uniqueIndex" json:"-"`
	KeyPrefix        string `json:"keyPrefix"`
	ProfileID        uint   `json:"profileId"`
	Profile          Profile
	AllowedLocations []MCPApiKeyLocation `gorm:"foreignKey:MCPApiKeyID" json:"allowedLocations"`
	LastUsedAt       *time.Time          `json:"lastUsedAt"`
	ExpiresAt        *time.Time          `json:"expiresAt"`
	IsActive         bool                `gorm:"default:true" json:"isActive"`
	IsInternal       bool                `gorm:"default:false" json:"isInternal"` // Internal keys see only internal tools
	gorm.Model
}

// MCPApiKeyLocation is the join table for API key to location mapping
type MCPApiKeyLocation struct {
	MCPApiKeyID uint     `gorm:"index" json:"mcpApiKeyId"`
	LocationID  string   `gorm:"index" json:"locationId"`
	Location    Location `gorm:"foreignKey:LocationID;references:Id" json:"-"`
	gorm.Model
}

// CanAccessLocation checks if the API key can access a specific location
func (k *MCPApiKey) CanAccessLocation(locationID string) bool {
	for _, loc := range k.AllowedLocations {
		if loc.LocationID == locationID {
			return true
		}
	}
	return false
}

// GetLocationIDs returns a slice of all allowed location IDs
func (k *MCPApiKey) GetLocationIDs() []string {
	ids := make([]string, len(k.AllowedLocations))
	for i, loc := range k.AllowedLocations {
		ids[i] = loc.LocationID
	}
	return ids
}

// GenerateMCPApiKey generates a new API key and returns the plain key, hash, and prefix
func GenerateMCPApiKey() (plainKey string, keyHash string, keyPrefix string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", err
	}
	plainKey = "sbmcp_" + hex.EncodeToString(bytes)
	keyHash = HashMCPApiKey(plainKey)
	keyPrefix = plainKey[:14] // "sbmcp_" + first 8 hex chars
	return plainKey, keyHash, keyPrefix, nil
}

// HashMCPApiKey returns the SHA256 hash of an API key
func HashMCPApiKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
