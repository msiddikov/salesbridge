package runway

import (
	cmn "client-runaway-zenoti/internal/common"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"time"
)

func LastNewLeadDate(locationId string) (time.Time, error) {
	res := time.Now()
	client, err := svc.NewClientFromId(locationId)
	if err != nil {
		return res, err
	}

	opp, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		Limit: 1,
		Order: "createdAt_desc",
	})

	if err != nil {
		return res, err
	}

	if len(opp) == 0 {
		return res, fmt.Errorf("there are no opportunities at all in location: %s", locationId)
	}

	return opp[0].CreatedAt, nil
}

func CheckForNewLeads() {
	locations := []models.Location{}
	db.DB.Find(&locations)

	for _, l := range locations {
		if l.TrackNewLeads {
			date, err := LastNewLeadDate(l.Id)
			if err != nil {
				msg := fmt.Sprintf("Error while checking last opportunity date in location %s: %s", l.Name, err.Error())
				cmn.NotifySlack("urgent", msg)
			}
			if time.Since(date) > 24*time.Hour {
				msg := fmt.Sprintf("%s: NO LEADS since %s", l.Name, date.Format("Jan 02, 2006"))
				cmn.NotifySlack("urgent", msg)
			}
		}
	}
}
