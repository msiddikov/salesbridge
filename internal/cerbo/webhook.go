package cerbo

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/packages/cerbo"
	"encoding/json"
	"fmt"
	"strings"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

var (
	subdomain = "medmatrixemr"
	cerboUser = "pk_hello123"
	cerboPass = "sk_hWhAqg_HaCesGMt_VfrxSyOt_FaGX"
)

func WebhookHandler(c *gin.Context) {
	// Extract the secret from the URL parameter
	secret := c.Param("secret")

	// Extract the body of the request
	webhookData := cerbo.WebhookData{}
	err := c.ShouldBindJSON(&webhookData)
	lvn.GinErr(c, 400, err, "Invalid webhook data")
	webhookData.Path = secret

	// switch based on the webhook type
	switch webhookData.EventType {
	// schedule.created or schedule.modified
	case "schedule.created", "schedule.modified":
		err = ScheduleUpsertHandler(webhookData)
		lvn.GinErr(c, 500, err, "Error while handling schedule upsert")
	}

	c.JSON(200, gin.H{"message": "Webhook received"})
}

func ScheduleUpsertHandler(data cerbo.WebhookData) error {
	// find triggers to be fired
	triggers := []models.GhlTrigger{}

	err := db.DB.Where("text_filter1 = ? AND text_filter2 = ?", data.PracticeId, data.Path).Find(&triggers).Error
	if err != nil {
		return err
	}

	dataStruct := cerbo.Schedule{}
	err = json.Unmarshal(data.Data, &dataStruct)
	if err != nil {
		return err
	}

	for _, provider := range dataStruct.AssignedProviders {
		dataStruct.Patient.Provider += provider.First
	}

	cli, err := cerbo.NewClient(subdomain, cerboUser, cerboPass)
	if err != nil {
		return err
	}
	patient, err := cli.GetPatient(dataStruct.Patient.Id)
	if err != nil {
		return err
	}

	user, err := cli.GetUser(fmt.Sprint(patient.PrimaryProviderId))
	if err != nil {
		return err
	}

	dataStruct.Patient.Provider = user.FirstName + " " + user.LastName
	tags := []string{}
	for _, tag := range patient.Tags {
		tags = append(tags, tag.Name)
	}

	dataStruct.Patient.Tags = strings.Join(tags, ", ")

	// update the data field
	dataBytes, err := lvn.Marshal(dataStruct)
	if err != nil {
		return err
	}

	data.Data = dataBytes

	// marshal the data to be sent to the triggers
	body, err := lvn.Marshal(data)
	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		svc := runway.GetSvc()
		cli, err := svc.NewClientFromId(trigger.LocationId)
		if err != nil {
			return err
		}

		err = cli.TriggerFire(trigger.TargetUrl, string(body))
		if err != nil {
			return err
		}
	}
	return nil
}
