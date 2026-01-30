package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"testing"
)

func GetMedMatrixClient() *runwayv2.Client {
	l := models.Location{}
	locationName := "Med Matrix"
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		panic("Location not found")
	}

	svc := runway.GetSvc()
	client, err := svc.NewClientFromId(l.Id)

	if err != nil {
		panic(err)
	}

	return &client
}

func TestExportMessages(t *testing.T) {
	client := GetMedMatrixClient()

	messages, err := client.MessagesExport(runwayv2.MessagesFilter{
		Channel:   "Call",
		ContactId: "vIzeURrBH5kRI79ujxXi",
	})

	fmt.Println(len(messages.Messages))
	if err != nil {
		t.Error(err)
	}

}

func TestGetTranscribe(t *testing.T) {
	client := GetMedMatrixClient()

	transcription, err := client.MessagesGetTranscription("zpZQPf5jyttRFa3mAaNk")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(transcription)
}
