package zenotiv1

import "time"

func (ag *AppointmentGroup) GetDate() time.Time {
	for _, a := range ag.Appointment_Services {
		return a.Start_time.Time
	}

	for _, p := range ag.Appointment_Packages {
		if !p.Service.Start_time.Time.IsZero() {
			return p.Service.Start_time.Time
		}

		for _, sa := range p.Appointment_services {
			return sa.Start_time.Time
		}
	}
	return time.Time{}
}
