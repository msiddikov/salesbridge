package models

import (
	"time"
)

// GoogleAdsConnection stores OAuth tokens and metadata for a Google Ads account per profile.
// jpmReport still uses the legacy service-account path; this model is for the new OAuth flow.
type GoogleAdsConnection struct {
	ID          uint `gorm:"primaryKey"`
	ProfileID   uint
	DisplayName string // Friendly label for UI selection
	Email       string // Account email (if available)

	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// GoogleAdsLocationSetting stores per-location account configuration for Google Ads.
// A location can either point directly to a customer account (CustomerID) or,
// when using a manager account, specify ManagerID + ClientCustomerID.
type GoogleAdsLocationSetting struct {
	ID               uint   `gorm:"primaryKey"`
	LocationId       string `gorm:"index:idx_ga_loc_profile,unique"`
	ProfileID        uint   `gorm:"index:idx_ga_loc_profile,unique"`
	ConnectionID     uint   `gorm:"index"`
	ConnectionName   string
	CustomerID       string
	CustomerName     string
	ManagerID        string
	ManagerName      string
	ClientCustomerID string
	ClientName       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
}
