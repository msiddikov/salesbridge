package zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
)

func MustGetClientFromLocation(l models.Location) zenotiv1.Client {

	client, _ := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	return client
}

func GetUsedUrlFromAPI(api string) string {
	loc := models.Location{}
	db.DB.Where("zenoti_api = ?", api).First(&loc)
	return loc.ZenotiUrl
}
