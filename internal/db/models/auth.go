package models

import "gorm.io/gorm"

type (
	Profile struct {
		Name      string
		OwnerID   uint
		Users     []User
		Locations []Location
		gorm.Model
	}

	User struct {
		Email     string
		Password  string
		Profile   Profile
		ProfileID uint
		Profiles  []Profile `gorm:"foreignKey:OwnerID"`
		gorm.Model
	}
)
