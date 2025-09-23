package db

import (
	"client-runaway-zenoti/internal/db/models"
	"fmt"
)

func Migrate() {
	err := DB.AutoMigrate(&models.Chat{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.ChatMessage{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Contact{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Sale{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Appointment{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.LocationExpense{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Location{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Calendar{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.BlockSlot{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.GhlTokens{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.GhlTrigger{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Setting{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.JpmReportInvoice{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.JpmReportNewLead{})
	if err != nil {
		panic(err)
	}

	defContact := models.Contact{
		LocationId:    "0",
		ContactId:     "0",
		OpportunityId: "0",
		FullName:      "Default contact",
	}

	DB.FirstOrCreate(&defContact)
	fmt.Println("Success")
}
