package svc_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"errors"

	"gorm.io/gorm"
)

func AGLinkRecord(appointmentGroup zenotiv1.AppointmentGroup, centerId string) error {
	link := models.ZenotiAppointmentGroupIdCenterIdLink{
		AppointmentGroupId: appointmentGroup.Appointment_Group_Id,
		CenterId:           centerId,
		InvoiceId:          appointmentGroup.Invoice_id,
	}

	for _, s := range appointmentGroup.Appointment_Services {
		link.OkToDelete = s.End_time.Time.AddDate(0, 1, 0) // 1 month after service end time
		break
	}

	return db.DB.Create(&link).Error
}

func GetCenterIdByAppointmentGroupId(appointmentGroupId string) (string, error) {

	link := models.ZenotiAppointmentGroupIdCenterIdLink{}
	err := db.DB.Where("appointment_group_id = ?", appointmentGroupId).First(&link).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "e2ff99ac-3be8-45aa-b8f1-c40f7fca6406", nil // default center id
	}

	if err != nil {
		return "", err
	}
	return link.CenterId, nil

}

func DeleteOldAGLinks() error {
	return db.DB.Where("ok_to_delete < NOW()").Delete(&models.ZenotiAppointmentGroupIdCenterIdLink{}).Error
}
