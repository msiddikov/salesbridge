package runway

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/tgbot"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/jinzhu/copier"
)

type (
	CollectionUpdateRes struct {
		Opportunity     runwayv2.Opportunity
		Guest           zenotiv1.Guest
		InvoiceId       string
		TotalCollection float64
		Updated         bool
		Registered      bool
		Reason          string
	}
)

func UpdateBookings(apts []zenotiv1.Appointment, l models.Location) (err error) {
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return err
	}

	for _, v := range apts {
		err := UpdateAppt(v, l, client)
		if err != nil {
			fmt.Printf("ERROR: %s continuing...\n", err)
			continue
		}
	}

	return nil
}

func UpdateAppt(appt zenotiv1.Appointment, l models.Location, client runwayv2.Client) error {

	// Findding an opportunity
	ops, err := client.OpportunitiesFindByEmailPhone(appt.Guest.Email, appt.Guest.Mobile.Number, l.PipelineId)
	if err != nil {
		return err
	}

	// skipping if no opportunity found
	if len(ops) == 0 {
		return nil
	}

	op := ops[0]

	wasAlreadyRegistered, err := registerBooking(op.Contact.Id, appt.Invoice_id, appt.Start_time.Time, l)
	if err != nil {
		return err
	}

	// if already registered, no need to update the stage
	if wasAlreadyRegistered {
		return nil
	}

	// continue if already in the right stage
	rightStage := ""
	switch appt.Status {
	case zenotiv1.NoShowed:
		rightStage = l.NoShowsId
	default:
		rightStage = l.BookId
	}

	if op.PipelineStageId == rightStage ||
		op.PipelineStageId == l.SalesId || op.PipelineStageId == l.MemberId {
		return nil
	}

	op.PipelineStageId = rightStage

	req := runwayv2.OpportunityUpdateReq{}
	copier.Copy(&req, &op)
	_, err = client.OpportunitiesUpdate(req, op.Id)

	if err != nil {
		return err
	}
	return nil
}

func UpdateApptGroup(apptGrp zenotiv1.AppointmentGroupWebhookData) error {

	l := models.Location{}
	err := db.DB.Where("zenoti_center_id = ? and sync_contacts=true", apptGrp.Center_Id).First(&l).Error
	if err != nil {
		return nil
	}

	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return err
	}

	zclient, err := zenotiv1.NewClient(l.Id, l.ZenotiCenterId, l.ZenotiApi)
	if err != nil {
		return err
	}

	guest, err := zclient.GuestsGetById(apptGrp.Guest.Id)
	if err != nil {
		return err
	}

	// Findding an opportunity
	ops, err := client.OpportunitiesFindByEmailPhone(guest.Personal_info.Email, guest.Personal_info.Mobile_phone.Number, l.PipelineId)
	if err != nil {
		return err
	}

	// skipping if no opportunity found
	if len(ops) == 0 {
		return nil
	}
	op := ops[0]

	hasUnregistered := false
	for _, appt := range apptGrp.Appointments {
		wasAlreadyRegistered, err := registerBooking(op.Contact.Id, apptGrp.Invoice_id, appt.Start_time.Time, l)
		if err != nil {
			return err
		}

		if !wasAlreadyRegistered {
			hasUnregistered = true
		}
	}

	// if no unregistered appointments, no need to update the stage
	if !hasUnregistered {
		return nil
	}

	rightStage := l.BookId

	if op.PipelineStageId == rightStage ||
		op.PipelineStageId == l.SalesId || op.PipelineStageId == l.MemberId {
		return nil
	}

	op.PipelineStageId = rightStage

	req := runwayv2.OpportunityUpdateReq{}
	copier.Copy(&req, &op)
	_, err = client.OpportunitiesUpdate(req, op.Id)

	if err != nil {
		return err
	}
	return nil
}

func UpdateSales(c []zenotiv1.Collection, l models.Location) error {
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %x sales for %s \n", len(c), l.Name)

	for _, v := range c {

		_, err := UpdateCollection(v, l, client)
		if err != nil {
			fmt.Printf("ERROR: %s continuing...\n", err)
			continue
		}

	}

	return nil
}

func UpdateCollection(c zenotiv1.Collection, l models.Location, client runwayv2.Client) (CollectionUpdateRes, error) {
	res := CollectionUpdateRes{
		Guest:           c.Guest,
		InvoiceId:       c.Invoice_id,
		TotalCollection: c.Total_collection,
	}

	// Findding an opportunity
	ops, err := client.OpportunitiesFindByEmailPhone(c.Guest.Personal_info.Email, c.Guest.Personal_info.Mobile_phone.Number, l.PipelineId)
	if err != nil {
		res.Reason = "Error finding opportunity"
		return res, err
	}

	// skipping if no opportunity found
	if len(ops) == 0 {
		res.Reason = "No opportunity found"
		return res, nil
	}

	op := ops[0]

	res.Opportunity = op
	op.PipelineStageId = l.SalesId

	wasRegistered, err := registerCollection(op.Contact.Id, c.Invoice_id, c.Total_collection, l)
	if err != nil {
		res.Reason = "Error registering collection"
		return res, err
	}

	if !wasRegistered {
		op.MonetaryValue += c.Total_collection
		res.Registered = true
	}

	req := runwayv2.OpportunityUpdateReq{}
	copier.Copy(&req, &op)
	_, err = client.OpportunitiesUpdate(req, op.Id)

	if err != nil {
		res.Reason = "Error api updating opportunity"
		return res, err
	}

	res.Updated = true
	return res, nil
}

func UpdateCollectionFromWebhook(data zenotiv1.WebhookData) error {
	// check if there is a location with this center id
	l := models.Location{}
	err := db.DB.Where("zenoti_center_id = ? and sync_contacts=true", data.Data.Invoice.Center_id).First(&l).Error
	if err != nil {
		return nil
	}

	// create a collection from the webhook data
	c := zenotiv1.Collection{
		Invoice_id: data.Data.Invoice.Id,
		Guest_id:   data.Data.Invoice.Guest.Id,
		Guest: zenotiv1.Guest{
			Id:        data.Data.Invoice.Guest.Id,
			Center_id: data.Data.Invoice.Center_id,
			Personal_info: zenotiv1.Personal_info{
				First_name: data.Data.Invoice.Guest.First_name,
				Last_name:  data.Data.Invoice.Guest.Last_name,
				Email:      data.Data.Invoice.Guest.Email,
				Mobile_phone: zenotiv1.Phone_info{
					Number: data.Data.Invoice.Guest.Mobile_phone,
				},
			},
		},
		Total_collection: data.Data.Invoice.Total_price.Sum_total,
	}

	// create a client
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		return err
	}

	// update the collection
	res, err := UpdateCollection(c, l, client)

	if res.Updated {
		reg := lvn.Ternary(res.Registered, "!! Registered", "Updated")

		tgbot.Notify("Webhook", fmt.Sprintf("%s sales for %s %s from %s", reg, c.Guest.Personal_info.First_name, c.Guest.Personal_info.Last_name, l.Name), false)
	}

	return err
}

// Helper functions ________________________________________
