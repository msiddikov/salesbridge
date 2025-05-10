package integrations_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func SyncCalendars() {
	fmt.Println("Starting calendar sync...")

	locs := []models.Location{}

	err := db.DB.Where("sync_calendars = ?", true).Find(&locs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, l := range locs {
		SyncCalendarsForLocation(l)
	}
}

func SyncCalendarsForLocation(l models.Location) {
	slotsToBe, err := getBlockSlotsToBe(l)
	if err != nil {
		fmt.Println(err)
		return
	}

	for c, slots := range slotsToBe {
		SyncCalendar(c, slots)
	}

}

func SyncCalendar(calendar models.Calendar, slots []models.BlockSlot) {
	toCreate := []models.BlockSlot{}
	toUpdate := []models.BlockSlot{}
	toDelete := []models.BlockSlot{}

	asIs := []models.BlockSlot{}
	db.DB.Where("calendar_id = ?", calendar.CalendarId).Find(&asIs)

	for _, s := range slots {
		found := false
		for _, ais := range asIs {
			if ais.ZenotiId == s.ZenotiId {
				found = true
				if !ais.StartTime.Equal(s.StartTime) || !ais.EndTime.Equal(s.EndTime) {
					s.Id = ais.Id
					toUpdate = append(toUpdate, s)
				}
			}
		}
		if !found {
			toCreate = append(toCreate, s)
		}
	}

	for _, ais := range asIs {
		found := false
		for _, s := range slots {
			if ais.ZenotiId == s.ZenotiId {
				found = true
				break
			}
		}
		if !found {
			toDelete = append(toDelete, ais)
		}
	}

	// saving changes to db
	for _, s := range toCreate {
		err := db.DB.Create(&s).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, s := range toUpdate {
		err := db.DB.Save(&s).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, s := range toDelete {
		err := db.DB.Delete(&s).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Helper functions _____________________________________________

func getBlockSlotsToBe(l models.Location) (map[models.Calendar][]models.BlockSlot, error) {
	res := map[models.Calendar][]models.BlockSlot{}
	appts, err := getAppointments(l)
	if err != nil {
		return res, err
	}

	blocksFromAppts := mapApptsToBlockSlots(appts, l)
	res = blocksFromAppts

	for c, _ := range res {

		client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
		if err != nil {
			fmt.Println(err)
			continue
		}
		blockOuts, err := client.EmployeesListBlockOutTimesAll(
			c.TherapistId,
			time.Now().Add(appointmentsFetchFromDays*24*time.Hour),
			time.Now().Add(appointmentsFetchToDays*24*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, bo := range blockOuts {
			slot := models.BlockSlot{
				LocationId: l.Id,
				CalendarId: c.CalendarId,
				StartTime:  bo.Start_time.Time,
				EndTime:    bo.End_time.Time,
				Title:      lvn.Ternary(bo.Notes == "", bo.Block_out_time_type.Name, bo.Notes),
				ZenotiId:   bo.Block_out_time_id,
			}
			res[c] = append(res[c], slot)
		}
	}

	return res, nil
}

func mapApptsToBlockSlots(appts []zenotiv1.Appointment, l models.Location) map[models.Calendar][]models.BlockSlot {
	res := map[models.Calendar][]models.BlockSlot{}
	calendars := []models.Calendar{}
	db.DB.Where("location_id = ?", l.Id).Find(&calendars)

	for _, a := range appts {
		if a.Start_time_utc == a.End_time_utc { // skip merged appointments
			continue
		}

		if a.Status == zenotiv1.Canceled { // skip canceled appointments
			continue
		}

		c, err := calendarFirsOrCreate(&calendars, l.Id, a.Therapist)
		if err != nil {
			fmt.Println(err)
			continue
		}

		res[c] = append(res[c], models.BlockSlot{
			LocationId: l.Id,
			CalendarId: c.CalendarId,
			StartTime:  a.Start_time_utc.Time,
			EndTime:    a.End_time_utc.Time,
			Title:      fmt.Sprintf("%s %s", a.Guest.First_name, a.Guest.Last_name),
			Notes:      fmt.Sprintf("Appointment: %s", l.GetZenotiAppointmentLink(a.Invoice_id)),
			ZenotiId:   a.Id,
		})
	}
	return res
}

func calendarFirsOrCreate(calendars *([]models.Calendar), locationId string, therapist zenotiv1.Therapist) (models.Calendar, error) {
	for _, c := range *calendars {
		if c.TherapistId == therapist.Id {
			return c, nil
		}
	}

	calendar, err := runway.CreateCalendarForTherapist(locationId, therapist)
	if err != nil {
		return models.Calendar{}, err
	}

	*calendars = append(*calendars, calendar)
	return calendar, nil
}
