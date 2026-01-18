package models

import "gorm.io/gorm"

type ZenotiApi struct {
	ApiName   string
	ProfileId uint
	ApiKey    string
	Url       string
	gorm.Model
}

type CerboApi struct {
	ApiName   string
	ProfileId uint
	Subdomain string
	Username  string
	ApiKey    string
	gorm.Model
}
