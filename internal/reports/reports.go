package reports

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"fmt"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"gorm.io/gorm/clause"
)

func GetExpenses(from, to time.Time, l []string) (float64, error) {
	expense := models.LocationExpense{}
	res := []struct {
		Total float64
		Num   int
	}{}

	from = lvn.Time(from).StartOfTheDay()
	to = lvn.Time(to).StartOfTheDay().Add(time.Second)

	err := db.DB.Model(&expense).Select("sum(total) as total, count(expense_id) as num").Where("date between ? and ? and location_id in ?", from, to, l).Find(&res).Error
	return res[0].Total, err
}

func GetSales(from, to time.Time, l []string) (float64, error) {
	sale := models.Sale{}
	res := []struct {
		Total float64
		Num   int
	}{}
	err := db.DB.Model(&sale).Joins("inner join contacts on contacts.contact_id=sales.contact_id and contacts.location_id in ?", l).Select("sum(total) as total, count(sale_id) as num, count(contacts.contact_id)").Where("date between ? and ? and not sales.contact_id in ?", from, to, l).Find(&res).Error
	return res[0].Total, err
}

func GetLeads(from, to time.Time, l []string) (float64, error) {
	contact := models.Contact{}
	var res float64
	err := db.DB.Model(&contact).Select("count(contact_id)").Where("created_date between ? and ? and location_id in ?", from, to, l).Find(&res).Error

	return res, err
}

func GetBookings(from, to time.Time, l []string) (float64, error) {
	appts := models.Appointment{}

	res := []struct {
		Total float64
		Num   float64
	}{}
	err := db.DB.Model(&appts).Joins("inner join contacts on contacts.contact_id=appointments.contact_id and contacts.location_id in ? and contacts.created_date between ? and ?", l, from, to).Select("sum(total) as total, count(contacts.contact_id)").Where("date between ? and ? and not appointments.contact_id in ?", from, to, l).Group("contacts.contact_id").Find(&res).Error

	return res[0].Num, err
}

func GetNoShows(from, to time.Time, l []string) (float64, error) {
	appts := models.Appointment{}

	res := []struct {
		Total float64
		Num   float64
	}{}
	err := db.DB.Model(&appts).Joins("inner join contacts on contacts.contact_id=appointments.contact_id and contacts.location_id in ? and contacts.created_date between ? and ?", l, from, to).Select("sum(total) as total, count(appointment_id) as num, count(contacts.contact_id)").Where("date between ? and ? and not appointments.contact_id in ? and status=-2", from, to, l).Find(&res).Error

	return res[0].Num, err
}

func SetExpenses(from, to time.Time, total float64, locations []string) {
	from = lvn.Time(from).StartOfTheDay()
	to = lvn.Time(to).StartOfTheDay()
	days := to.Sub(from).Hours() / 24
	daily := total / days / float64(len(locations))
	for _, l := range locations {
		for i := 0; i < int(days); i++ {
			date := from.Add(time.Duration(i) * 24 * time.Hour)
			expense := models.LocationExpense{
				ExpenseId:  fmt.Sprintf("%s%s", l, date.Format("2006-01-02")),
				Date:       date,
				Total:      daily,
				LocationId: l,
			}
			db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&expense)
		}
	}
}
