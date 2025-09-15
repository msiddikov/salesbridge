package meta

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
)

func init() {

}

func getTestClient() Client {
	setting := models.Setting{
		Key: "meta_token",
	}

	err := db.DB.First(&setting, "key = ?", "meta_token").Error
	if err != nil {
		return Client{}
	}

	s := Service{}
	cli, _ := s.NewClient(setting.Value)
	return cli
}
