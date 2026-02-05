package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"errors"
	"testing"

	"gorm.io/gorm"
)

// TestUpdateInternalMcpKeys updates (or creates if missing) internal MCP keys for all profiles.
func TestUpdateInternalMcpKeys(t *testing.T) {
	var profiles []models.Profile
	if err := db.DB.Preload("MCPApiKeys").Find(&profiles).Error; err != nil {
		t.Fatalf("failed to load profiles: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("no profiles found")
	}

	for _, profile := range profiles {
		plainKey, keyHash, keyPrefix, err := models.GenerateMCPApiKey()
		if err != nil {
			t.Fatalf("profile %d: failed to generate api key: %v", profile.ID, err)
		}

		tx := db.DB.Begin()
		if tx.Error != nil {
			t.Fatalf("profile %d: failed to start transaction: %v", profile.ID, tx.Error)
		}

		var key models.MCPApiKey
		err = tx.Where("profile_id = ? AND is_internal = ?", profile.ID, true).Order("id").First(&key).Error
		switch {
		case err == nil:
			if err := tx.Model(&key).Updates(map[string]interface{}{
				"key_hash":   keyHash,
				"plain_key":  plainKey,
				"key_prefix": keyPrefix,
				"is_active":  true,
			}).Error; err != nil {
				tx.Rollback()
				t.Fatalf("profile %d: failed to update internal key: %v", profile.ID, err)
			}
		case errors.Is(err, gorm.ErrRecordNotFound):
			key = models.MCPApiKey{
				Name:       "internal-profile-" + keyPrefix,
				PlainKey:   plainKey,
				KeyHash:    keyHash,
				KeyPrefix:  keyPrefix,
				ProfileID:  profile.ID,
				IsActive:   true,
				IsInternal: true,
			}
			if err := tx.Create(&key).Error; err != nil {
				tx.Rollback()
				t.Fatalf("profile %d: failed to create internal key: %v", profile.ID, err)
			}
		default:
			tx.Rollback()
			t.Fatalf("profile %d: failed to load internal key: %v", profile.ID, err)
		}

		if err := tx.Commit().Error; err != nil {
			t.Fatalf("profile %d: failed to commit: %v", profile.ID, err)
		}

		t.Logf("profile %d internal MCP key updated: %s", profile.ID, plainKey)
	}
}
