package models

import "gorm.io/gorm"

type ZenotiApi struct {
	ApiName   string
	ProfileId uint
	ApiKey    string
	Url       string
	gorm.Model
}
