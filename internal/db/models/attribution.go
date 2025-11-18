package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	Person struct {
		ProfileId  uint
		LocationId string

		FirstName string
		LastName  string
		Email     string
		Phone     string

		ContactId     string
		OpportunityId string

		EmrSystem      string
		EmrContactId   string
		EMRContactLink string

		gorm.Model
	}

	AttributionFlow struct {
		ProfileId uint
		FlowName  string
		Stages    datatypes.JSONSlice[string] `gorm:"type:jsonb"`
		gorm.Model
	}

	StageHit struct {
		ProfileId         uint
		PersonId          string
		AttributionFlowId uint
		Stage             string
		OccurredAt        time.Time
		Revenue           float64

		RefSystem string // e.g. "zenoti", "sales-bridge"
		RefId     string // e.g. zenoti appointment id, zenoti invoice id
		RefLink   string

		gorm.Model
	}
)
