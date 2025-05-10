package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"

	"gorm.io/gorm/clause"
)

func CalendarForTherapistFirsOrCreate(locationId string, therapist zenotiv1.Therapist) (models.Calendar, error) {
	res := models.Calendar{}
	db.DB.Where("therapist_id = ? AND location_id = ?", therapist.Id, locationId).Find(&res)
	if res.CalendarId != "" {
		return res, nil
	}

	return CreateCalendarForTherapist(locationId, therapist)
}

func CreateCalendarForTherapist(locationId string, therapist zenotiv1.Therapist) (models.Calendar, error) {
	res := models.Calendar{}
	// getting client
	client, err := svc.NewClientFromId(locationId)
	if err != nil {
		return res, err
	}

	therapistFirstNameWithInitial := therapist.First_name
	if therapist.Last_name != "" {
		therapistFirstNameWithInitial += " " + string(therapist.Last_name[0]) + "."
	}

	calendar, err := client.CalendarCreate(runwayv2.CalendarCreateReq{
		Name:       fmt.Sprintf("Zenoti (%s)", therapistFirstNameWithInitial),
		LocationId: locationId,
	})
	if err != nil {
		return res, err
	}

	res = models.Calendar{
		CalendarId:    calendar.Id,
		TherapistId:   therapist.Id,
		TherapistName: fmt.Sprintf("%s %s", therapist.First_name, therapist.Last_name),
		LocationId:    locationId,
	}

	err = db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&res).Error
	if err != nil {
		client.CalendarsDelete(calendar.Id)
	}

	return res, err
}
