package models

import (
	"time"

	"gorm.io/gorm"
)

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

// We need this link because the Zenoti Appointment Group status change do not provide CenterId
type ZenotiAppointmentGroupIdCenterIdLink struct {
	AppointmentGroupId string
	InvoiceId          string
	CenterId           string
	OkToDelete         time.Time
	gorm.Model
}
