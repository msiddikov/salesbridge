package tests

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
)

func getClients(locationName string) (models.Location, zenotiv1.Client, runwayv2.Client, error) {
	l := models.Location{}
	db.DB.Where("name = ?", locationName).First(&l)

	if l.Id == "" {
		return l, zenotiv1.Client{}, runwayv2.Client{}, fmt.Errorf("location not found")
	}

	svc := runway.GetSvc()
	rclient, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return l, zenotiv1.Client{}, runwayv2.Client{}, err
	}

	zClient, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		return l, zenotiv1.Client{}, runwayv2.Client{}, err
	}

	return l, zClient, rclient, nil
}
