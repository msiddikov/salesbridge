package integrations

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/zenotiLegacy/runaway"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/schollz/progressbar/v3"
)

func UpdateStatusesDb() {
	fmt.Println("Starting status updates...")
	for _, l := range config.GetLocations() {
		UpdateStatusesDbForLocation(l)
	}
}

func UpdateStatusesDbForLocation(l config.Location) {
	// Updating bookings
	fmt.Println("Getting bookings for " + l.Name)
	aptsAll, err := zenoti.GetAllAppointments(l)
	// apts := zenoti.GetNewAppts(aptsAll, l)
	apts := aptsAll
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updating bookings for " + l.Name)
	err = runaway.UpdateBookings(apts, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Booking successfully updated for " + l.Name)

	// Updating NoShows

	fmt.Println("Getting bookings for " + l.Name)
	aptsAll, err = zenoti.GetPastAppointments(l)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Updating bookings for " + l.Name)
	err = runaway.UpdateBookings(apts, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Booking successfully updated for " + l.Name)

	// Updating Sales
	start := config.GetSalesSyncDate(l)
	if start.IsZero() || start.Before(time.Now().Add(-3*24*time.Hour)) {
		start = time.Now().Add(-3 * 24 * time.Hour)
	}
	end := time.Now().Add(24 * time.Hour)

	fmt.Println("Getting sales for " + l.Name)
	s, err := zenoti.GetCollections(l, start, end)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Updating sales for " + l.Name)
	err = runaway.UpdateSales(s, l)
	if err != nil {
		fmt.Println(err)
		return
	}

	config.SetSalesSyncDate(l, end.Add(0-24*time.Hour))
	fmt.Println("Sales successfully updated for " + l.Name)
}

func UpdateNotes(force bool) {
	for _, l := range config.GetLocations() {
		println("Updating Notes for " + l.Name)
		// getting contacts
		opps := runaway.GetAllOpportunities(l)
		bar := progressbar.Default(int64(len(opps)))

		for _, c := range opps {
			runaway.UpdateNote(c, l, force)
			bar.Add(1)
		}
	}

	fmt.Println("All notes have been updated")
}

func StartScheduledJobs() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(4).Hours().Do(runway.ForceCheckAll)
	s.Every(2).Hours().Do(func() {
		UpdateStatusesDb()
	})
	s.StartBlocking()
}

func UpdateHistoricSales(from, to time.Time) {

	for _, l := range config.GetLocations() {

		bar := progressbar.Default(int64(to.Sub(from).Hours() / 24))
		fmt.Println("Fetching sales and bookings for " + l.Name)
		start := from
		end := start.Add(24 * time.Hour)
		for start.Before(to) {
			bar.Add(1)

			// Updating Sales
			s, err := zenoti.GetCollections(l, start, end)
			if err != nil {
				fmt.Println(err)
				continue
			}

			runaway.SaveSalesToDB(s, l)

			start = end
			end = end.Add(24 * time.Hour)
		}

		fmt.Println("Sales successfully fetched for " + l.Name)
	}
}

func UpdateHistoricBookings(from, to time.Time) {
	for _, l := range config.GetLocations() {

		bar := progressbar.Default(int64(to.Sub(from).Hours() / 24))
		fmt.Println("Fetching sales and bookings for " + l.Name)
		start := from
		end := start.Add(24 * time.Hour)
		for start.Before(to) {
			bar.Add(1)
			//Updating bookings
			a, err := zenoti.GetAppointments(l, start, end)
			if err != nil {
				fmt.Println(err)
				continue
			}

			runaway.SaveBookingsToDB(a, l)
			start = end
			end = end.Add(24 * time.Hour)
		}

		fmt.Println("Sales successfully fetched for " + l.Name)
	}
}
