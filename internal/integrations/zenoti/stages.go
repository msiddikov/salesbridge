package integrations_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/tgbot"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func UpdateStages() {

	locs := []models.Location{}

	err := db.DB.Where("sync_contacts = ?", true).Find(&locs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, l := range locs {
		UpdateStagesForLocation(l)
	}
}

func UpdateStagesForLocation(l models.Location) {

	tgbot.Notify("Stages", "Starting Update stages for "+l.Name+"...", false)

	// check if the runway client is available
	rwsvc := runway.GetSvc()
	client, err := rwsvc.NewClientFromId(l.Id)

	if err != nil {
		tgbot.Notify("Stages-errors", "Error in UpdateStages for "+l.Name+" "+err.Error(), false)
		return
	}

	_, err = client.LocationsInfo()
	if err != nil {
		tgbot.Notify("Stages-errors", "Error in UpdateStages for "+l.Name+" "+err.Error(), false)
		return
	}

	// update Bookings
	err = UpdateBookings(l)
	if err != nil {
		tgbot.Notify("Stages-errors", "Error in UpdateBookings for "+l.Name+" "+err.Error(), false)
	}

	// update Sales
	err = UpdateSalesChunks(l)
	if err != nil {
		tgbot.Notify("Stages-errors", "Error in UpdateSales for "+l.Name+" "+err.Error(), false)
	}

	tgbot.Notify("Stages", "Finished Update stages for "+l.Name+"...", false)
}

func UpdateBookings(l models.Location) error {

	apts, err := getAppointments(l)
	if err != nil {
		return err
	}

	fmt.Println("Updating bookings for " + l.Name)
	err = runway.UpdateBookings(apts, l)
	if err != nil {
		return err
	}

	fmt.Println("Booking successfully updated for " + l.Name)
	return nil
}

func UpdateSales(l models.Location) error {

	s, err := getSales(l)
	if err != nil {
		return err
	}

	fmt.Println("Updating sales for " + l.Name)
	err = runway.UpdateSales(s, l)
	if err != nil {
		return err
	}

	l.SalesSyncDate = time.Now()
	db.DB.Save(&l)
	fmt.Println("Sales successfully updated for " + l.Name)
	return nil
}

func UpdateSalesChunks(l models.Location) error {

	to := time.Now()
	curTo := l.SalesSyncDate.Add(24 * time.Hour)
	curTo = lvn.Ternary(curTo.After(to), to, curTo)

	for curTo.After(l.SalesSyncDate) {
		s, err := GetSalesChunks(l, l.SalesSyncDate, curTo)
		if err != nil {
			return err
		}

		fmt.Println("Updating sales for " + l.Name)
		err = runway.UpdateSales(s, l)
		if err != nil {
			return err
		}

		l.SalesSyncDate = curTo
		db.DB.Save(&l)

		l.SalesSyncDate = curTo
		db.DB.Save(&l)
		curTo = l.SalesSyncDate.Add(5 * 24 * time.Hour)
		curTo = lvn.Ternary(curTo.After(to), to, curTo)
	}
	return nil

}

func GetSalesChunks(l models.Location, start time.Time, end time.Time) ([]zenotiv1.Collection, error) {

	client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		return nil, err
	}

	fmt.Println("Getting sales for " + l.Name)
	s, err := client.ReportsAllCollections(start, end)
	if err != nil {
		return nil, err
	}

	res := []zenotiv1.Collection{}
	// filling in guest info and recalculating total collection
	for i, sale := range s {
		guest, err := client.GuestsGetById(sale.Guest_id)
		if err != nil {
			continue
		}
		s[i].Guest = guest

		s[i].Total_collection = 0
		for _, item := range sale.Items {
			s[i].Total_collection += item.Final_sale_price
		}
		res = append(res, s[i])
	}

	return res, nil
}

// we cannot ask for more than 10 days in API, so we need to make multiple requests
func getAppointments(l models.Location) ([]zenotiv1.Appointment, error) {

	fmt.Println("Getting bookings for " + l.Name)

	client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	apts, err := client.AppointmentsListAllAppointments(zenotiv1.AppointmentFilter{
		StartDate:           now.Add(appointmentsFetchFromDays * 24 * time.Hour),
		EndDate:             now.Add(appointmentsFetchToDays*24*time.Hour + 1*time.Second),
		IncludeNoShowCancel: false,
	})

	return apts, err
}

func getSales(l models.Location) ([]zenotiv1.Collection, error) {

	client, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		return nil, err
	}

	start := l.SalesSyncDate
	if start.IsZero() {
		start = time.Now().Add(-3 * 24 * time.Hour)
	}
	end := time.Now().Add(24 * time.Hour)

	fmt.Println("Getting sales for " + l.Name)
	s, err := client.ReportsAllCollections(start, end)
	if err != nil {
		return nil, err
	}

	// filling in guest info and recalculating total collection
	for i, sale := range s {
		guest, err := client.GuestsGetById(sale.Guest_id)
		if err != nil {
			return nil, err
		}
		s[i].Guest = guest

		s[i].Total_collection = 0
		for _, item := range sale.Items {
			s[i].Total_collection += item.Final_sale_price
		}
	}

	return s, nil
}
