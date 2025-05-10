package runway

import (
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"time"
)

type (
	SurveyForm struct {
		Name    string
		Phone   string
		Email   string
		Answers []struct {
			Question string
			Answer   string
		}
	}
)

func SurveyPost(form SurveyForm, locationId, workflowId string) error {
	client, _ := svc.NewClientFromId(locationId)

	cfValues := []runwayv2.CustomFieldValue{}

	for _, a := range form.Answers {
		cf, err := client.CustomFieldsFirstOrCreate(runwayv2.CustomField{
			Name:     a.Question,
			DataType: "TEXT",
		})
		if err != nil {
			return err
		}

		cfValues = append(cfValues, runwayv2.CustomFieldValue{
			Id:          cf.Id,
			Field_value: a.Answer,
		})
	}

	contact, err := client.ContactsFirstOrCreate(runwayv2.Contact{
		FirstName: form.Name,
		Phone:     form.Phone,
		Email:     form.Email,
	})
	if err != nil {
		return err
	}

	contact.UpdateCustomFields(cfValues)
	contact, err = client.ContactsUpdate(contact)
	if err != nil {
		return err
	}
	if workflowId == "undefined" {
		return nil
	}
	err = client.ContactsAddToWorkflow(contact.Id, workflowId, time.Now())
	return err
}

func GetWorkflows(locationId string) ([]runwayv2.Workflow, error) {
	cli, _ := svc.NewClientFromId(locationId)

	return cli.WorkflowsGet()
}
