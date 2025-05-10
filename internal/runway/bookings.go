package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"time"
)

func CreateBooking(payload AppointmentCreateWebhookPayload) error {

	// getting location
	location := models.Location{}
	db.DB.Where("id = ?", payload.Location.Id).First(&location)

	// getting zenoti and runway clients
	zcli, err := zenotiv1.NewClient(location.Id, location.ZenotiCenterId, location.ZenotiApi)
	if err != nil {
		return err
	}

	// getting guest or creating one
	zguests, err := zcli.GuestsGetByPhoneEmail(payload.Phone, payload.Email)
	if err != nil && err.Error() != "guest not found" {
		return err
	}

	if len(zguests) == 0 {
		guest := zenotiv1.Guest{
			Center_id: location.ZenotiCenterId,
			Personal_info: zenotiv1.Personal_info{
				First_name: payload.First_name,
				Last_name:  payload.Last_name,
				Email:      payload.Email,
				Mobile_phone: zenotiv1.Phone_info{
					Country_code: 225,
					Number:       payload.Phone,
				},
				DateOfBirth: zenotiv1.ZenotiTime{Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
		}
		zguest, err := zcli.GuestsCreate(guest)
		if err != nil {
			return err
		}
		zguests = append(zguests, zguest)
	}

	// getting calendar
	calendar := models.Calendar{}
	db.DB.Where("location_id = ? and calendar_id=?", location.Id, payload.Calendar.Id).First(&calendar)

	// fixing time
	// timezone, err := time.LoadLocation(payload.Calendar.SelectedTimeZone)
	// if err != nil {
	// 	return err
	// }
	// payload.Calendar.StartTime.Time = payload.Calendar.StartTime.Time.In(timezone)

	// booking
	req := zenotiv1.BookingReq{
		Date:     zenotiv1.ZenotiDate{Time: payload.Calendar.StartTime.Time},
		CenterId: location.ZenotiCenterId,
		Guests: []zenotiv1.BookingReqGuest{
			{
				Id: zguests[0].Id,
				Items: []zenotiv1.BookingReqGuestsItems{
					{
						Item: zenotiv1.BookingReqItem{
							Id: location.ZenotiServiceId,
						},
						Therapist: zenotiv1.BookingReqTherapist{
							Id: calendar.TherapistId,
						},
					},
				},
			},
		},
	}

	_, err = zcli.BookWithConfirm(req)
	if err != nil {
		return err
	}

	return nil
}
