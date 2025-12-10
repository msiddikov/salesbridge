package svc_googleads

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/googleads"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	svc *googleads.Service
)

func init() {
	svcObj := googleads.NewService(
		config.Confs.Settings.GoogleAdsClientId,
		config.Confs.Settings.GoogleAdsClientSecret,
		config.Confs.Settings.GoogleAdsRedirectUrl,
		saveConnection,
		updateTokens,
		getConnection,
	)
	svc = &svcObj
}

func saveConnection(conn googleads.Connection) (googleads.Connection, error) {
	model := toModel(conn)
	if model.DisplayName == "" {
		nameSource := model.Email
		if nameSource == "" {
			nameSource = "google ads connection"
		}
		model.DisplayName = fmt.Sprintf("%s's google ads", nameSource)
	}
	if err := db.DB.Save(&model).Error; err != nil {
		return conn, err
	}
	return fromModel(model), nil
}

func updateTokens(conn googleads.Connection) error {
	if conn.ID == 0 {
		return fmt.Errorf("googleads oauth: connection id is required to update tokens")
	}
	return db.DB.Model(&models.GoogleAdsConnection{}).Where("id = ?", conn.ID).
		Updates(map[string]any{
			"access_token":  conn.AccessToken,
			"refresh_token": conn.RefreshToken,
			"token_expiry":  conn.TokenExpiry,
		}).Error
}

func getConnection(id uint) (googleads.Connection, error) {
	m := models.GoogleAdsConnection{}
	if err := db.DB.First(&m, id).Error; err != nil {
		return googleads.Connection{}, err
	}
	return fromModel(m), nil
}

func toModel(conn googleads.Connection) models.GoogleAdsConnection {
	return models.GoogleAdsConnection{
		ID:           conn.ID,
		ProfileID:    conn.ProfileID,
		DisplayName:  conn.DisplayName,
		Email:        conn.Email,
		AccessToken:  conn.AccessToken,
		RefreshToken: conn.RefreshToken,
		TokenExpiry:  conn.TokenExpiry,
	}
}

func fromModel(m models.GoogleAdsConnection) googleads.Connection {
	return googleads.Connection{
		ID:           m.ID,
		ProfileID:    m.ProfileID,
		DisplayName:  m.DisplayName,
		Email:        m.Email,
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		TokenExpiry:  m.TokenExpiry,
	}
}

func cliForLocation(locationId string, profileId uint) (*googleads.Client, error) {
	setting := models.GoogleAdsLocationSetting{}
	err := db.DB.Where("location_id = ? AND profile_id = ?", locationId, profileId).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("location google ads settings not found")
		}
		return nil, fmt.Errorf("unable to fetch setting: %w", err)
	}

	return svc.NewClient(setting.ConnectionID, googleads.CustomerInfo{
		CustomerID:       setting.CustomerID,
		ClientCustomerID: setting.ClientCustomerID,
		ManagerID:        setting.ManagerID,
	})
}
