package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func TriggerSubscriptionsHandler(c *gin.Context) {
	var event runwayv2.TriggerSubscriptionData
	err := c.ShouldBindJSON(&event)
	lvn.GinErr(c, 400, err, "Error while binding trigger subscription data")

	// If it is a delete event, delete the trigger from the database
	if event.TriggerData.EventType == string(runwayv2.TriggerSubscriptionEventDelete) {
		err = db.DB.Where("id = ?", event.TriggerData.Id).Delete(&models.GhlTrigger{}).Error
		lvn.GinErr(c, 500, err, "Error while deleting trigger")
		c.Data(lvn.Res(200, nil, "Trigger deleted successfully"))
		return
	}

	// If it is a create or update event, upsert the trigger in the database
	trigger := triggerFromSubscriptionData(event)

	err = db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&trigger).Error
	lvn.GinErr(c, 500, err, "Error while upserting trigger")

	c.Data(lvn.Res(200, trigger, "Trigger upserted successfully"))
}

func triggerFromSubscriptionData(event runwayv2.TriggerSubscriptionData) models.GhlTrigger {
	res := models.GhlTrigger{
		Id:         event.TriggerData.Id,
		Key:        event.TriggerData.Key,
		Version:    event.Meta.Version,
		LocationId: event.Extras.LocationId,
		WorkflowId: event.Extras.WorkflowId,
		TargetUrl:  event.TriggerData.TargetUrl,
	}

	for _, filter := range event.TriggerData.Filters {
		switch filter.Field {
		// cerbo trigger
		case "practice_id":
			res.TextFilter1 = filter.Value[0]
		case "path":
			res.TextFilter2 = filter.Value[0]
		default:
		}
	}

	return res
}
