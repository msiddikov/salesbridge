package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/tgbot"
	"fmt"
)

func UpdateAllTokens() {
	locations := []models.Location{}

	db.DB.Find(&locations)

	for _, l := range locations {
		client, err := svc.NewClientFromId(l.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = client.UpdateToken()
		if err != nil {
			tgbot.Notify("Tokens failures", fmt.Sprintf("Location: %s, Error: %s", l.Name, err.Error()), false)
		}
	}
}
