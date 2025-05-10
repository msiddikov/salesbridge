package runaway

import (
	cmn "client-runaway-zenoti/internal/common"
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/types"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti/export"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ChangeOpportunity(op types.Opportunity, l config.Location) error {

	params := types.OpportunityChangeParams{
		Title:         op.Name,
		StageId:       op.PipelineStageId,
		Status:        op.Status,
		MonetaryValue: op.MonetaryValue,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "PUT",
		Endpoint: "/pipelines/" + l.PipelineId + "/opportunities/" + op.Id,
		Body:     string(body),
		Api:      l.Api,
	})
	if err != nil {
		return err
	}
	fmt.Printf("[STAGE CHANGE] %s moved to %s\n", op.Name, op.PipelineStageId)
	return nil
}

func SaveSalesToDB(c []types.Collection, l config.Location) {

	for _, v := range c {

		//fmt.Println(v.Invoice_id)
		if (v.Guest.Personal_info.Email == "" && v.Guest.Personal_info.Mobile_phone.Number == "") || v.Guest.Personal_info.First_name == "" {
			continue
		}
		contact, err := GetContact(types.Contact{
			Phone:     v.Guest.Personal_info.Mobile_phone.Number,
			Email:     v.Guest.Personal_info.Email,
			FirstName: v.Guest.Personal_info.First_name,
			LastName:  v.Guest.Personal_info.Last_name,
		}, l, false)
		if err != nil && err.Error() == "Contact not found" {

			continue
		}
		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		// saving to db
		sale := models.Sale{
			SaleId:    v.Invoice_id,
			Date:      v.Created_Date.Time,
			Total:     v.Total_collection,
			ContactId: contact.Id,
		}
		db.DB.Where("sale_id=?", v.Invoice_id).FirstOrCreate(&sale)
	}

}

func UpdateSales(c []types.Collection, l config.Location) error {
	bar := progressbar.Default(int64(len(c)))
	for _, v := range c {
		bar.Add(1)
		//fmt.Println(v.Invoice_id)
		if (v.Guest.Personal_info.Email == "" && v.Guest.Personal_info.Mobile_phone.Number == "") || v.Guest.Personal_info.First_name == "" {
			continue
		}
		contact, err := GetContact(types.Contact{
			Phone:     v.Guest.Personal_info.Mobile_phone.Number,
			Email:     v.Guest.Personal_info.Email,
			FirstName: v.Guest.Personal_info.First_name,
			LastName:  v.Guest.Personal_info.Last_name,
		}, l, false)
		if err != nil && err.Error() == "Contact not found" {

			continue
		}
		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		op, err := GetOpportunity(types.Opportunity{
			Name:            contact.FirstName + " " + contact.LastName,
			PipelineStageId: l.SalesId,
			Contact:         contact,
			PipelineId:      l.PipelineId,
		}, l, true)
		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}
		op.PipelineStageId = l.SalesId

		// saving to db
		contactToCreate := models.Contact{
			LocationId:    l.Id,
			ContactId:     op.Contact.Id,
			OpportunityId: op.Id,
			FullName:      op.Contact.FirstName + " " + op.Contact.LastName,
			CreatedDate:   op.CreatedAt,
		}

		db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&contactToCreate)

		sale := models.Sale{
			SaleId:    v.Invoice_id,
			Date:      v.Created_Date.Time,
			Total:     v.Total_collection,
			ContactId: contact.Id,
		}
		db.DB.Where("sale_id=?", v.Invoice_id).FirstOrCreate(&sale)

		isRegistered := isCollectionRegistered(contact.Id, v.Invoice_id, l)
		if isRegistered && op.PipelineStageId == l.SalesId {
			continue
		}
		if isRegistered && op.PipelineStageId != l.SalesId {
			err = ChangeOpportunity(op, l)
			if err != nil {
				fmt.Printf("ERROR: %s continuing...", err)
				continue
			}
			continue
		}

		op.MonetaryValue += v.Total_collection
		err = ChangeOpportunity(op, l)
		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		err = registerCollection(contact.Id, v.Invoice_id, v.Total_collection, l)
		if err != nil {
			op.MonetaryValue -= v.Total_collection
			ChangeOpportunity(op, l)
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

	}

	return nil
}

func isCollectionRegistered(id, invoiceId string, l config.Location) bool {
	notes, err := getNotes(id, l)
	if err != nil {
		return true
	}

	for _, n := range notes {
		if strings.Contains(n.Body, "Collections for this contact") && strings.Contains(n.Body, invoiceId) {
			return true
		}
	}
	return false
}

func registerCollection(id, invoiceId string, value float64, l config.Location) error {
	notes, err := getNotes(id, l)
	if err != nil {
		return err
	}

	link := l.Zenoti.Url + "/Appointment/DlgAppointment1.aspx?invoiceid=" + invoiceId
	for _, n := range notes {
		if strings.Contains(n.Body, "Collections for this contact") {
			collections := strings.Split(n.Body, "\n")
			for _, c := range collections {
				if strings.Contains(c, link) {
					return nil
				}
			}
			err := updateNote(id, n.Id, fmt.Sprintf("%s\n%.2f: %s", n.Body, value, link), l)
			return err
		}
	}

	return setNote(id, fmt.Sprintf("Collections for this contact:\n%.2f: %s", value, link), l)
}

func registerBooking(id, invoiceId string, date time.Time, l config.Location) error {

	notes, err := getNotes(id, l)
	if err != nil {
		return err
	}

	link := l.Zenoti.Url + "/Appointment/DlgAppointment1.aspx?invoiceid=" + invoiceId
	for _, n := range notes {
		if strings.Contains(n.Body, "Bookings for this contact") {
			if strings.Contains(n.Body, invoiceId) { //already registered
				return nil
			}

			collections := strings.Split(n.Body, "\n")
			for _, c := range collections {
				if strings.Contains(c, link) {
					return nil
				}
			}
			err := updateNote(id, n.Id, fmt.Sprintf("%s\n%s: %s", n.Body, date.Format("01/02/2006"), link), l)
			return err
		}
	}
	return setNote(id, fmt.Sprintf("Bookings for this contact:\n%s: %s", date.Format("01/02/2006"), link), l)
}

func SaveBookingsToDB(apts []types.Appointment, l config.Location) {

	for _, v := range apts {

		//fmt.Println(v.Id)
		if (v.Guest.Mobile.Number == "" && v.Guest.Email == "") || v.Guest.First_name == "" {
			continue
		}
		contact, err := GetContact(types.Contact{
			Phone:     v.Guest.Mobile.Number,
			Email:     v.Guest.Email,
			FirstName: v.Guest.First_name,
			LastName:  v.Guest.Last_name,
		}, l, false)
		if err != nil && err.Error() == "Contact not found" {
			continue
		}

		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		appt := models.Appointment{
			AppointmentId: v.Id,
			ContactId:     contact.Id,
			Date:          v.Start_time.Time,
			Status:        v.Status,
			Total:         float64(v.Price.Sales),
		}

		err = db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&appt).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdateBookings(apts []types.Appointment, l config.Location) (err error) {

	bar := progressbar.Default(int64(len(apts)))

	for _, v := range apts {
		bar.Add(1)
		// continue if no contact info
		if (v.Guest.Mobile.Number == "" && v.Guest.Email == "") || v.Guest.First_name == "" {
			continue
		}

		contact, err := GetContact(types.Contact{
			Phone:     v.Guest.Mobile.Number,
			Email:     v.Guest.Email,
			FirstName: v.Guest.First_name,
			LastName:  v.Guest.Last_name,
		}, l, false)

		// continue if contact not found
		if err != nil && err.Error() == "Contact not found" {
			appt := models.Appointment{
				AppointmentId: v.Id,
				ContactId:     l.Id,
				Status:        v.Status,
				Total:         float64(v.Price.Sales),
				Date:          v.Start_time.Time,
			}
			db.DB.Create(&appt)
			continue
		}

		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		op, err := GetOpportunity(types.Opportunity{
			Name:            v.Guest.First_name + " " + v.Guest.Last_name,
			PipelineStageId: l.BookId,
			Contact:         contact,
			PipelineId:      l.PipelineId,
		}, l, true)
		if err != nil {
			fmt.Printf("ERROR: %s continuing...", err)
			continue
		}

		contactToCreate := models.Contact{
			LocationId:    l.Id,
			ContactId:     op.Contact.Id,
			OpportunityId: op.Id,
			FullName:      op.Contact.FirstName + " " + op.Contact.LastName,
			CreatedDate:   op.CreatedAt,
		}

		err1 := db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&contactToCreate).Error

		appt := models.Appointment{
			AppointmentId: v.Id,
			ContactId:     op.Contact.Id,
			Date:          v.Start_time.Time,
			Status:        v.Status,
			Total:         float64(v.Price.Sales),
		}

		err = db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&appt).Error
		if err != nil && err1 != nil {
			UpdateNote(op, l, false)
			err = db.DB.Create(&appt).Error
			if err != nil {
				fmt.Println(err)
			}
		}

		err = registerBooking(contact.Id, v.Invoice_id, v.Start_time.Time, l)
		if err != nil {
			fmt.Printf("ERROR while registering booking: %s continuing...\n", err)
			continue
		}

		// continue if already in the right stage
		rightStage := ""
		switch v.Status {
		case types.NoShowed:
			rightStage = l.NoShowsId
		default:
			rightStage = l.BookId
		}

		if op.PipelineStageId == rightStage ||
			op.PipelineStageId == l.SalesId {
			continue
		}

		op.PipelineStageId = rightStage

		err = ChangeOpportunity(op, l)

		if err != nil {
			fmt.Printf("ERROR while changing opportunity: %s continuing...\n", err)
			continue
		}

	}

	return nil
}

// creates contact without checking for existence
func CreateContact(c types.Contact, l config.Location) (types.Contact, error) {
	empty := types.Contact{}
	//return empty, fmt.Errorf("Creating new contacts is disabled...")
	body, err := json.Marshal(c)
	if err != nil {
		return types.Contact{}, err
	}
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "POST",
		Endpoint: "/contacts",
		Body:     string(body),
		Api:      l.Api,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Contact types.Contact
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Contact, nil
}

func GetPipelines(api string) ([]types.Pipeline, error) {
	empty := []types.Pipeline{}
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/pipelines",
		Api:      api,
	})

	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Pipelines []types.Pipeline
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Pipelines, nil
}

func GetPipeline(l config.Location) (types.Pipeline, error) {
	pplns, err := GetPipelines(l.Api)
	if err != nil {
		return types.Pipeline{}, err
	}
	for _, p := range pplns {
		if p.Id == l.PipelineId {
			return p, nil
		}
	}
	return types.Pipeline{}, errors.New("pipeline not found")
}

func GetContact(c types.Contact, l config.Location, createIfNotExist bool) (types.Contact, error) {

	params := []cmn.QueryParams{
		{Key: "email", Value: c.Email},
		{Key: "phone", Value: c.Phone},
	}

	//res, err := runawayReq("GET", "contacts/lookup", "", api, params)
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/contacts/lookup",
		Api:      l.Api,
		QParams:  params,
	})
	if err != nil && res.StatusCode != 404 && res.StatusCode != 422 {
		return types.Contact{}, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if (res.StatusCode == 404 || res.StatusCode == 422) && !createIfNotExist {
		return types.Contact{}, fmt.Errorf("contact not found")
	}

	if (res.StatusCode == 404 || res.StatusCode == 422) && createIfNotExist {
		return CreateContact(c, l)
	}

	if err != nil {
		return types.Contact{}, err
	}

	contacts := struct {
		Contacts []types.Contact
	}{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return types.Contact{}, err
	}

	return contacts.Contacts[0], nil
}

func GetContactById(id string, l config.Location) (types.Contact, error) {

	//res, err := runawayReq("GET", "contacts/lookup", "", api, params)
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/contacts/" + id,
		Api:      l.Api,
	})

	if err != nil {
		return types.Contact{}, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return types.Contact{}, err
	}

	contacts := struct {
		Contact types.Contact
	}{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return types.Contact{}, err
	}

	return contacts.Contact, nil
}

func GetOpportunity(op types.Opportunity, l config.Location, createIfNotExist bool) (types.Opportunity, error) {
	empty := types.Opportunity{}
	if l.PipelineId == "" {
		return empty, fmt.Errorf("pipeline id is not set")
	}
	if op.Contact.Email != "" {
		opt, err := getOpportunityByQuery(op.Contact.Email, l)
		if err == nil {
			return opt, nil
		}
		if err.Error() != "Not found" {
			return empty, err
		}
	}

	if op.Contact.Phone != "" {
		opt, err := getOpportunityByQuery(op.Contact.Phone, l)
		if err == nil {
			return opt, nil
		}
		if err.Error() != "Not found" {
			return empty, err
		}
	}

	if !createIfNotExist {
		return empty, fmt.Errorf("Not found")
	}

	param := map[string]string{
		"title":     op.Name,
		"stageId":   op.PipelineStageId,
		"contactId": op.Contact.Id,
		"status":    "open",
	}
	body, _ := json.Marshal(param)

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "POST",
		Endpoint: "/pipelines/" + l.PipelineId + "/opportunities",
		Body:     string(body),
		Api:      l.Api,
	})
	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := types.Opportunity{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result, nil
}

func getOpportunityByQuery(query string, l config.Location) (types.Opportunity, error) {

	params := []cmn.QueryParams{
		{Key: "query", Value: query},
	}

	//res, err := runawayReq("GET", "pipelines/"+pipelineId+"/opportunities", "", api, params)
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/pipelines/" + l.PipelineId + "/opportunities",
		Api:      l.Api,
		QParams:  params,
	})
	if err != nil && res.StatusCode != 404 {
		return types.Opportunity{}, err
	}

	if res.StatusCode == 404 {
		return types.Opportunity{}, fmt.Errorf("Not found")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return types.Opportunity{}, err
	}

	ops := struct {
		Opportunities []types.Opportunity
	}{}
	err = json.Unmarshal(data, &ops)
	if err != nil {
		return types.Opportunity{}, err
	}
	if len(ops.Opportunities) == 0 {
		return types.Opportunity{}, fmt.Errorf("Not found")
	}

	return ops.Opportunities[0], nil
}

func GetAllOpportunities(l config.Location) (ids []types.Opportunity) {
	ids = []types.Opportunity{}
	meta := types.Meta{}
	fetchedAll := false

	for !fetchedAll {
		opps, meta1, err := getOpportunitiesByMeta(l, meta, "open", "")
		if err != nil {
			fmt.Println(err)
			continue
		}
		meta = meta1
		if meta.StartAfter == 0 {
			fetchedAll = true
		}

		ids = append(ids, opps...)
	}

	meta = types.Meta{}
	fetchedAll = false

	for !fetchedAll {
		opps, meta1, err := getOpportunitiesByMeta(l, meta, "won", "")
		if err != nil {
			fmt.Println(err)
			continue
		}
		meta = meta1
		if meta.StartAfter == 0 {
			fetchedAll = true
		}

		ids = append(ids, opps...)

	}
	return
}

func getOpportunitiesByMeta(l config.Location, meta types.Meta, status, stageId string) ([]types.Opportunity, types.Meta, error) {
	empty := []types.Opportunity{}
	emptyMeta := types.Meta{}

	params := []cmn.QueryParams{
		{Key: "limit", Value: "100"},
		{Key: "status", Value: status},
	}
	if stageId != "" {
		params = append(params, cmn.QueryParams{
			Key: "stageId", Value: stageId,
		})
	}

	if meta != emptyMeta {

		params = append(params, cmn.QueryParams{Key: "startAfterId", Value: meta.StartAfterId},
			cmn.QueryParams{Key: "startAfter", Value: fmt.Sprintf("%v", meta.StartAfter)})
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/pipelines/" + l.PipelineId + "/opportunities",
		Api:      l.Api,
		QParams:  params,
	})

	if err != nil {
		return empty, emptyMeta, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, emptyMeta, err
	}

	opps := struct {
		Opportunities []types.Opportunity
		Meta          types.Meta
	}{}
	err = json.Unmarshal(data, &opps)
	if err != nil {
		return empty, emptyMeta, err
	}

	return opps.Opportunities, opps.Meta, nil
}
func getOpportunitiesByMetaPeriod(from, to time.Time, l config.Location, meta types.Meta, status, stageId string) ([]types.Opportunity, types.Meta, error) {
	empty := []types.Opportunity{}
	emptyMeta := types.Meta{}

	fromS := fmt.Sprint(from)
	toS := fmt.Sprint(to)

	params := []cmn.QueryParams{
		{Key: "limit", Value: "100"},
		{Key: "startDate", Value: fromS},
		{Key: "endDate", Value: toS},
	}

	if stageId != "" {
		params = append(params, cmn.QueryParams{
			Key: "stageId", Value: stageId,
		})
	}
	if status != "" {
		params = append(params, cmn.QueryParams{
			Key: "status", Value: status,
		})
	}

	if meta != emptyMeta {

		params = append(params, cmn.QueryParams{Key: "startAfterId", Value: meta.StartAfterId},
			cmn.QueryParams{Key: "startAfter", Value: fmt.Sprintf("%v", meta.StartAfter)})
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/pipelines/" + l.PipelineId + "/opportunities",
		Api:      l.Api,
		QParams:  params,
	})

	if err != nil {
		return empty, emptyMeta, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, emptyMeta, err
	}

	opps := struct {
		Opportunities []types.Opportunity
		Meta          types.Meta
	}{}
	err = json.Unmarshal(data, &opps)
	if err != nil {
		return empty, emptyMeta, err
	}

	return opps.Opportunities, opps.Meta, nil
}

func GetAllGuests(l config.Location, meta types.Meta) ([]types.Contact, types.Meta, error) {
	empty := []types.Contact{}
	emptyMeta := types.Meta{}

	params := []cmn.QueryParams{
		{Key: "limit", Value: "100"},
	}

	if meta != emptyMeta {

		params = []cmn.QueryParams{
			{Key: "startAfterId", Value: meta.StartAfterId},
			{Key: "startAfter", Value: fmt.Sprintf("%v", meta.StartAfter)},
			{Key: "limit", Value: "100"},
		}
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/contacts",
		Api:      l.Api,
		QParams:  params,
	})

	if err != nil {
		return empty, emptyMeta, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, emptyMeta, err
	}

	conts := struct {
		Contacts []types.Contact
		Meta     types.Meta
	}{}
	err = json.Unmarshal(data, &conts)
	if err != nil {
		return empty, emptyMeta, err
	}

	return conts.Contacts, conts.Meta, nil
}

func UpdateLocations() {
	l := config.GetLocations()
	for k, v := range l {
		ppln, err := GetPipelines(v.Api)
		if err != nil {
			panic(fmt.Sprintf("Error updating %s: %s", v.Id, err))
		}

		pipeline := types.Pipeline{}
		for _, pipe := range ppln {
			if pipe.Name == v.Pipeline {
				l[k].PipelineId = pipe.Id
				pipeline = pipe
				break
			}
		}
		if pipeline.Id == "" {
			panic("For " + v.Id + " pipeline not found: " + v.Pipeline)
		}

		l[k].SalesId = stageId(v.Id, v.Sales, pipeline)
		l[k].BookId = stageId(v.Id, v.Book, pipeline)
		l[k].NoShowsId = stageId(v.Id, v.NoShows, pipeline)
		l[k].LeadsId = stageId(v.Id, v.Leads, pipeline)

		defContact := models.Contact{
			ContactId:     v.Id,
			LocationId:    v.Id,
			OpportunityId: v.Id,
			FullName:      "Default contact",
		}
		db.DB.Where("contact_id=?", v.Id).FirstOrCreate(&defContact)

	}
	config.UpdateLocations(l)
	fmt.Println("Locations successfully updated")
}

func stageId(locName, stageName string, pipeline types.Pipeline) string {
	if stageName == "" {
		return ""
	}
	for _, s := range pipeline.Stages {
		if s.Name == stageName {
			return s.Id
		}
	}
	panic(fmt.Sprintf("For location %s pipline %s stage %s not found", locName, pipeline.Name, stageName))
}

func UpdateNote(c types.Opportunity, l config.Location, force bool) {
	contact := models.Contact{}

	err := db.DB.Where("contact_id=?", c.Contact.Id).First(&contact).Error
	if !force && err != gorm.ErrRecordNotFound {
		return
	}
	updateGuest := true
	updateLink := true
	notes, err := getNotes(c.Contact.Id, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range notes {
		//checking link
		if strings.Contains(n.Body, "Please follow this link to create a guest in Zenoti") {
			if force {
				err := deleteNote(c.Contact.Id, n.Id, l)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				updateLink = false
			}
		}

		//checking updates in zenoti
		if strings.Contains(n.Body, "Zenoti update:") {
			if force {
				err := deleteNote(c.Contact.Id, n.Id, l)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				updateGuest = false
			}
		}
	}
	if updateLink {
		link := config.Confs.Settings.SrvDomain + "/contact/" + c.Contact.Id + "/" + l.Id
		msg := "Please follow this link to create a guest in Zenoti: " + link
		setNote(c.Contact.Id, msg, l)
	}

	if updateGuest {
		err := export.UpdateGuest(c.Contact.Email, c.Contact.Email, l)
		note := "Zenoti update: updated"
		if err != nil {
			note = "Zenoti update: " + err.Error()
		}
		setNote(c.Contact.Id, note, l)
	}

	// Updating db
	opp := models.Contact{
		ContactId:     c.Contact.Id,
		LocationId:    l.Id,
		OpportunityId: c.Id,
		FullName:      c.Name,
		CreatedDate:   c.CreatedAt,
	}

	db.DB.Where("Contact_id=?", c.Contact.Id).FirstOrCreate(&opp)

}

func getNotes(id string, l config.Location) ([]types.Note, error) {
	empty := []types.Note{}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/contacts/" + id + "/notes",
		Api:      l.Api,
	})

	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Notes []types.Note
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Notes, nil
}

func deleteNote(id, noteId string, l config.Location) error {

	_, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "DELETE",
		Endpoint: "/contacts/" + id + "/notes/" + noteId,
		Api:      l.Api,
	})

	if err != nil {
		return err
	}

	return nil
}

func setNote(id, note string, l config.Location) error {

	body := struct {
		Body string `json:"body"`
	}{
		Body: note,
	}
	bodyBytes, _ := json.Marshal(body)

	_, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "POST",
		Endpoint: "/contacts/" + id + "/notes",
		Api:      l.Api,
		Body:     string(bodyBytes),
	})

	if err != nil {
		return err
	}
	return nil
}

func updateNote(id, noteId, noteBody string, l config.Location) error {

	body := struct {
		Body string `json:"body"`
	}{
		Body: noteBody,
	}
	bodyBytes, _ := json.Marshal(body)

	_, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "PUT",
		Endpoint: "/contacts/" + id + "/notes/" + noteId,
		Api:      l.Api,
		Body:     string(bodyBytes),
	})

	if err != nil {
		return err
	}

	return nil
}

func QueryContacts(s string, l config.Location) []types.Contact {
	empty := []types.Contact{}

	if s == "" {
		return empty
	}

	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/contacts",
		Api:      l.Api,
		QParams: []cmn.QueryParams{
			{
				Key:   "query",
				Value: s,
			},
			{
				Key:   "limit",
				Value: "20",
			},
		},
	})

	if err != nil {
		return empty
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty
	}

	result := struct {
		Contacts []types.Contact
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty
	}
	c := result.Contacts
	for k := range c {
		c[k].FirstName = strings.Title(c[k].FirstName)
		c[k].LastName = strings.Title(c[k].LastName)
	}

	return result.Contacts
}

func GetAllLocations() ([]types.Location, error) {
	empty := []types.Location{}
	res, err := cmn.Req(cmn.ReqParams{
		Platform: "R",
		Method:   "GET",
		Endpoint: "/locations",
		Api:      config.Confs.Settings.CRAgencyAPI,
	})

	if err != nil {
		return empty, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	result := struct {
		Locations []types.Location
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return empty, err
	}

	return result.Locations, nil
}
