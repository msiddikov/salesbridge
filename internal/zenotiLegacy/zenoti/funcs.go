package zenoti

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	testdata "client-runaway-zenoti/internal/tests/data"
	"client-runaway-zenoti/internal/types"
	"time"
)

// gets appointments for 30 days from now
func GetAllAppointments(l config.Location) ([]types.Appointment, error) {
	start := time.Now()
	end := start.Add(7 * 24 * time.Hour)
	return GetAppointmentsPeriod(l, start, end)
}

// gets appointments for 30 days before now
func GetPastAppointments(l config.Location) ([]types.Appointment, error) {
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	return GetAppointmentsPeriod(l, start, end)
}

// get all appointments for a location for period
func GetAppointmentsPeriod(l config.Location, start, end time.Time) ([]types.Appointment, error) {
	if config.IsTesting() {
		return testdata.Appointments(), nil
	}
	empty := []types.Appointment{}
	appointments := []types.Appointment{}
	a, err := GetAppointments(l, start, end)
	if err != nil {
		return empty, err
	}
	appointments = append(appointments, a...)

	return appointments, nil
}

func CompareAppointments(data []types.Appointment, log []string) []types.Appointment {
	newApts := []types.Appointment{}

	// recognize new ones
	for _, v := range data {
		if contains(v.Id, log) {
			continue
		}
		newApts = append(newApts, v)
	}
	return newApts
}

func contains(v string, a []string) bool {
	for _, s := range a {
		if s == v {
			return true
		}
	}
	return false
}

func containsAppt(id string, a []models.Appointment) bool {
	for _, s := range a {
		if s.AppointmentId == id {
			return true
		}
	}
	return false
}

func GetNewAppts(appts []types.Appointment, l config.Location) []types.Appointment {
	newApts := []types.Appointment{}
	allApts := []models.Appointment{}

	db.DB.Joins("Contact", db.DB.Where("location_id=?", l.Id)).Where("date>?", l.Zenoti.AppointmentsSyncDay).Find(&allApts)

	for _, v := range appts {
		if containsAppt(v.Id, allApts) {
			continue
		}
		newApts = append(newApts, v)
	}
	return newApts
}
