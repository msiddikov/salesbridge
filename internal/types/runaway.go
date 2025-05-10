package types

import (
	"time"
)

type (
	Contact struct {
		Id        string `json:"id"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	Meta struct {
		StartAfterId string
		StartAfter   int
		Total        int
	}

	Note struct {
		Id   string
		Body string
	}

	Pipeline struct {
		Id     string
		Name   string
		Stages []struct {
			Id   string
			Name string
		}
	}

	Opportunity struct {
		Id              string
		Name            string
		MonetaryValue   float64
		Status          string
		Contact         Contact
		PipelineStageId string
		PipelineId      string
		CreatedAt       time.Time
	}

	OpportunityChangeParams struct {
		Title         string  `json:"title"`
		StageId       string  `json:"stageId"`
		Status        string  `json:"status"`
		MonetaryValue float64 `json:"monetaryValue"`
	}

	Location struct {
		Id           string
		Name         string
		IsIntegrated bool
	}
)
