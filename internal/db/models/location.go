package models

import (
	cmn "client-runaway-zenoti/internal/common"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (l *Location) Get(id string) error {
	err := DB.Where("id=?", id).First(&l).Error
	return err
}

func (l *Location) Save() error {
	err := DB.Save(&l).Error
	return err
}

func (l *Location) GetZenotiAppointmentLink(id string) string {
	url := l.ZenotiApiObj.Url
	if url == "" {
		url = l.ZenotiUrl
	}
	return url + "/Appointment/DlgAppointment1.aspx?invoiceid=" + id
}

func (l *Location) GetGhlTokensInstance() (GhlTokens, error) {
	tokensInstance := GhlTokens{}
	err := DB.Where("location_id=?", l.Id).First(&tokensInstance).Error
	return tokensInstance, err
}

func (l *Location) GetGhlTokens() (accessToken string, refreshToken string, err error) {
	tokensInstance, err := l.GetGhlTokensInstance()
	return tokensInstance.AccessToken, tokensInstance.RefreshToken, err
}

func (l *Location) SaveGhlTokens(accessToken, refreshToken string) error {
	tokensInstance, _ := l.GetGhlTokensInstance()
	tokensInstance.AccessToken = accessToken
	tokensInstance.RefreshToken = refreshToken
	tokensInstance.LocationId = l.Id

	err := DB.Clauses(clause.OnConflict{UpdateAll: true}).Save(&tokensInstance).Error

	return err
}

func (l *Location) BeforeUpdate(tx *gorm.DB) (err error) {
	// notify if autoCreateContacts is changed
	if tx.Statement.Changed("auto_create_contacts") {
		// notify slack
		notifyMsg := fmt.Sprintf("Auto create contacts for %s has been changed to %s", l.Name, lvn.Ternary(l.AutoCreateContacts, "true", "false"))
		cmn.NotifySlack("", notifyMsg)
	}

	return
}
