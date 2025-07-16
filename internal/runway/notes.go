package runway

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"fmt"
	"strings"
	"time"
)

func UpdateNote(o runwayv2.Opportunity, location models.Location, force bool) {
	// contact := models.Contact{}

	// if the contact is in the db return
	// err := db.DB.Where("contact_id=?", o.Contact.Id).First(&contact).Error
	// if !force && err != gorm.ErrRecordNotFound {
	// 	return
	// }

	// getting client
	client, err := svc.NewClientFromId(location.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	notes, err := client.ContactsGetAllNotes(o.Contact.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// checking if the note is already updated
	updateGuest := true
	updateLink := true
	for _, n := range notes {
		//checking link
		if strings.Contains(n.Body, "Please follow this link to create a guest in Zenoti") || strings.Contains(n.Body, "Please follow this link to view existing guest or create a new one in Zenoti") {
			if force {
				err := client.ContactsDeleteNote(n)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				updateLink = false
			}
		}

		//checking updates in zenoti
		if strings.Contains(n.Body, "Zenoti update:") {
			// force update for the time being
			if force || true {
				err := client.ContactsDeleteNote(n)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				updateGuest = false
			}
		}
	}

	if updateLink {
		link := config.Confs.Settings.SrvDomain + "/contact/" + o.Contact.Id + "/" + location.Id
		msg := "Please follow this link to view existing guest or create a new one in Zenoti: " + link
		_, err = client.ContactsCreateNote(
			runwayv2.Note{
				ContactId: o.Contact.Id,
				Body:      msg,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	if updateGuest {
		msg := "Zenoti update:"

		_, err := CheckAutoCreateOpportunity(o, location)
		if err != nil {
			fmt.Println(err)
			msg += " " + err.Error()
		}

		c, err := zenotiv1.NewClient(location.Id, location.ZenotiCenterId, location.ZenotiApi)
		if err != nil {
			fmt.Println(err)
			msg += " " + err.Error()
		}

		err = c.GuestsUpdate(o.Contact.Email, o.Contact.Phone)
		if err != nil {
			msg += " " + err.Error()
		} else {
			msg += " Guest updated"
		}

		note := runwayv2.Note{
			ContactId: o.Contact.Id,
			Body:      msg,
		}
		_, err = client.ContactsCreateNote(note)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Updating db
	opp := models.Contact{
		ContactId:     o.Contact.Id,
		LocationId:    o.Id,
		OpportunityId: o.Id,
		FullName:      o.Name,
		CreatedDate:   o.CreatedAt,
	}

	db.DB.Where("Contact_id=?", o.Contact.Id).FirstOrCreate(&opp)
}

func HasNote(o runwayv2.Opportunity, note string, client runwayv2.Client) bool {

	notes, err := client.ContactsGetAllNotes(o.Contact.Id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// checking if the note is already updated
	for _, n := range notes {
		//checking link
		if strings.Contains(n.Body, note) {
			return true
		}
	}

	return false
}

func registerBooking(id, invoiceId string, date time.Time, l models.Location) (wasAlreadyRegistered bool, err error) {

	// getting client
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// getting notes
	notes, err := client.ContactsGetAllNotes(id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	link := l.GetZenotiAppointmentLink(invoiceId)

	for _, n := range notes {
		if strings.Contains(n.Body, "Bookings for this contact") {
			if strings.Contains(n.Body, invoiceId) { //already registered
				return true, nil
			}

			collections := strings.Split(n.Body, "\n")
			for _, c := range collections {
				if strings.Contains(c, link) {
					return false, nil
				}
			}

			_, err := client.ContactsUpdateNote(runwayv2.Note{
				ContactId: id,
				Id:        n.Id,
				Body:      fmt.Sprintf("%s\n%s: %s", n.Body, date.Format("01/02/2006"), link),
			})
			return false, err
		}
	}
	_, err = client.ContactsCreateNote(runwayv2.Note{
		ContactId: id,
		Body:      fmt.Sprintf("Bookings for this contact:\n%s: %s", date.Format("01/02/2006"), link),
	})
	return false, err
}

func registerCollection(id, invoiceId string, value float64, l models.Location) (wasRegistered bool, err error) {
	// getting client
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// getting notes
	notes, err := client.ContactsGetAllNotes(id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	link := l.GetZenotiAppointmentLink(invoiceId)

	for _, n := range notes {
		if strings.Contains(n.Body, "Collections for this contact") {
			collections := strings.Split(n.Body, "\n")
			for _, c := range collections {
				if strings.Contains(c, link) {
					return true, nil
				}
			}
			_, err := client.ContactsUpdateNote(runwayv2.Note{
				ContactId: id,
				Id:        n.Id,
				Body:      fmt.Sprintf("%s\n%.2f: %s", n.Body, value, link),
			})

			return false, err
		}
	}
	_, err = client.ContactsCreateNote(runwayv2.Note{
		ContactId: id,
		Body:      fmt.Sprintf("Collections for this contact:\n%.2f: %s", value, link),
	})
	return false, err
}
