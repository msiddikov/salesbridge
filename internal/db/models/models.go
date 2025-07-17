package models

import (
	"client-runaway-zenoti/internal/types"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"time"

	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Svc *runwayv2.Service
)

type (
	ChatMessage struct {
		Date        time.Time
		Content     string
		Inbound     bool
		ManagerId   uint
		ManagerName string
		ChatId      uint
		Phone       string
		gorm.Model
	}

	Chat struct {
		LocationId  string
		ContactId   string
		ContactName string
		PhoneNo     string
		Messages    []ChatMessage `gorm:"foreignKey:ChatId"`
		RcId        string
		gorm.Model
	}

	Contact struct {
		LocationId    string
		ContactId     string `gorm:"primaryKey"`
		OpportunityId string
		FullName      string
		Sales         []Sale        `gorm:"foreignKey:ContactId"`
		Appointments  []Appointment `gorm:"foreignKey:ContactId"`
		CreatedDate   time.Time
	}

	Sale struct {
		SaleId    string `gorm:"primaryKey"`
		Date      time.Time
		ContactId string
		Contact   Contact
		Total     float64
	}

	Appointment struct {
		AppointmentId string `gorm:"primaryKey"`
		Date          time.Time
		ContactId     string
		Contact       Contact
		Total         float64
		Status        types.ZenotiStatus //NoShow = -2, Cancelled = -1, New = 0, Closed = 1, Checkin = 2, Confirm = 4, Break = 10, NotSpecified = 11, Available = 20, and Voided = 21
	}

	LocationExpense struct {
		ExpenseId  string `gorm:"primaryKey"`
		LocationId string
		Date       time.Time
		Total      float64
	}

	Location struct {
		Name               string
		Id                 string `gorm:"primaryKey"`
		PipelineId         string
		NewId              string
		BookId             string
		SalesId            string
		NoShowsId          string
		ShowNoSaleId       string
		MemberId           string
		TrackNewLeads      bool
		ZenotiApi          string
		ZenotiUrl          string
		ZenotiCenterId     string
		ZenotiCenterName   string
		ZenotiServiceId    string
		ZenotiServiceName  string
		ZenotiServicePrice float32
		SyncCalendars      bool
		SyncContacts       bool
		AutoCreateContacts bool
		ForceCheck         bool
		SalesSyncDate      time.Time
		Rank               int64
	}

	GhlTrigger struct {
		Id          string `gorm:"primaryKey"`
		Key         string
		Version     string
		LocationId  string
		WorkflowId  string
		TargetUrl   string
		TextFilter1 string
		TextFilter2 string
	}

	GhlTokens struct {
		LocationId   string `gorm:"primaryKey"`
		AccessToken  string
		RefreshToken string
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	ZenotiIntegrations struct {
		Api        string `gorm:"primaryKey"`
		Url        string
		ReferralId string
	}

	LocationInfo struct {
		Pipelines []runwayv2.Pipeline
		Workflows []runwayv2.Workflow
		Location  Location
	}

	Calendar struct {
		TherapistId   string `gorm:"primaryKey"`
		TherapistName string
		LocationId    string
		CalendarId    string
	}

	BlockSlot struct {
		Id         string `gorm:"primaryKey"`
		LocationId string
		CalendarId string
		StartTime  time.Time
		EndTime    time.Time
		Title      string
		Notes      string
		ZenotiId   string
	}
)
