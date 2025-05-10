package models

import "gorm.io/gorm"

func (c *Contact) HasSalesOrAppointments(db *gorm.DB) (bool, error) {
	err := db.Where(c).Preload("Sales").Preload("Appointments").First(c).Error
	return len(c.Sales) > 0 || len(c.Appointments) > 0, err
}
