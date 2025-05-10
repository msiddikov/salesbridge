package integrations_zenoti

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
)

func UpdateNotesV2(force bool) {
	locations := []models.Location{}
	db.DB.Where("sync_contacts = ?", true).Find(&locations)
	for _, l := range locations {
		err := UpdateNotesV2Location(l, force)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdateNotesV2Location(l models.Location, force bool) error {
	rwsvc := runway.GetSvc()
	client, err := rwsvc.NewClientFromId(l.Id)
	if err != nil {
		return err
	}

	opps, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		PipelineId: l.PipelineId,
	})

	if err != nil {
		return err
	}

	for _, o := range opps {
		runway.UpdateNote(o, l, force)
	}

	return nil
}

func OpportunitiesWithNotes(l models.Location, note string) ([]runwayv2.Opportunity, []string, error) {
	res := []runwayv2.Opportunity{}
	ress := []string{}
	rwsvc := runway.GetSvc()
	client, err := rwsvc.NewClientFromId(l.Id)
	if err != nil {
		return res, ress, err
	}

	opps, err := client.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		PipelineId: l.PipelineId,
	})

	for k, o := range opps {
		if k > 1000000 {
			fmt.Println("break")
		}
		if runway.HasNote(o, note, client) {
			res = append(res, o)
			ress = append(ress, o.Name)
		}
	}

	return res, ress, nil
}
