package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"testing"
	"time"
)

var (
	testLocId         = "EgSKJqMXN6w7t2svxveY"
	testTemplateLocId = "qMMRPtVblxzpx9KfyoZE"
	fairlawnLocId     = "TpwQvq1uDohQXHFebMQj"
	testCalendarId    = "i7qAafOINi5PLTrtSg7j"
	testEventId       = "NoTFidNDQVM6bKe8nIea"
	avonLocId         = "jEbtFz8wbaNI1WWDDFIb"
	casaLocId         = "kS6owm5hnT0mCVoDDPwO"
)

func TestGetCalendars(t *testing.T) {
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)
	calendars, err := client.CalendarsGet()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(calendars)
}

func TestGetCalendarEvents(t *testing.T) {
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)

	events, err := client.CalendarGetEvents("ZheBEkvQ9gQdHOq2G7e1")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(events)
}

func TestCreateCalendar(t *testing.T) {
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)

	calendar, err := client.CalendarCreate(runwayv2.CalendarCreateReq{
		Name:       "Test Calendar",
		LocationId: testLocId,
	})

	if err != nil {
		t.Error(err)
	}

	fmt.Println(calendar)
}

func TestCreateBlockSlot(t *testing.T) {

	slot := runwayv2.BlockSlot{
		LocationId:    testLocId,
		CalendarId:    testCalendarId,
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(time.Hour),
		Title:         "Test Block Slot",
		CalendarNotes: "Test Notes and stuff",
	}

	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)

	blockSlot, err := client.CalendarCreateBlockSlot(slot)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(blockSlot)
}

func TestEditBlocSlot(t *testing.T) {
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)

	blockSlot, err := client.CalendarEditBlockSlot(testEventId, runwayv2.BlockSlot{
		LocationId:    testLocId,
		CalendarId:    testCalendarId,
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(time.Hour),
		Title:         "Test Block Slot 2",
		CalendarNotes: "Test Notes",
	})

	if err != nil {
		t.Error(err)
	}

	fmt.Println(blockSlot)
}

func TestDeleteEvent(t *testing.T) {
	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(testLocId)

	err := client.CalendarDeleteEvent(testEventId)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateToken(t *testing.T) {
	l := models.Location{}
	locationName := "Hand & Stone Massage and Facial Spa - Chesapeake"
	db.DB.Where("name = ?", locationName).First(&l)

	svc := runway.GetSvc()
	client, _ := svc.NewClientFromId(l.Id)

	err := client.UpdateToken()

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateAllTokens(t *testing.T) {
	runway.UpdateAllTokens()
}
