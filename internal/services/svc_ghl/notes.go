package svc_ghl

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"strings"
	"time"
)

func UpdateLinkNote(contactId string, location models.Location, force bool) error {

	// getting client
	client, err := svc.NewClientFromId(location.Id)
	if err != nil {
		return err
	}

	notes, err := client.ContactsGetAllNotes(contactId)
	if err != nil {
		return err
	}

	// checking if the note is already updated
	updateLink := true

	for _, n := range notes {
		//checking link
		if strings.Contains(n.Body, "Please follow this link to create a guest in Zenoti") || strings.Contains(n.Body, "Please follow this link to view existing guest or create a new one in Zenoti") {
			if force {
				err := client.ContactsDeleteNote(n)
				if err != nil {
					return err
				}
			} else {
				updateLink = false
			}
		}
	}

	if updateLink {
		link := config.Confs.Settings.SrvDomain + "/contact/" + contactId + "/" + location.Id
		msg := "Please follow this link to view existing guest or create a new one in Zenoti: " + link
		_, err = client.ContactsCreateNote(
			runwayv2.Note{
				ContactId: contactId,
				Body:      msg,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterBooking(contactId, invoiceId string, date time.Time, l models.Location) (wasAlreadyRegistered bool, err error) {

	// getting client
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// getting notes
	notes, err := client.ContactsGetAllNotes(contactId)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	link := l.GetZenotiAppointmentLink(invoiceId)

	for _, n := range notes {
		if strings.Contains(n.Body, "Bookings for this contact") {
			if strings.Contains(n.Body, link) { //already registered
				return true, nil
			}

			_, err := client.ContactsUpdateNote(runwayv2.Note{ // note found but not registered
				ContactId: contactId,
				Id:        n.Id,
				Body:      fmt.Sprintf("%s\n%s: %s", n.Body, date.Format("01/02/2006"), link),
			})
			return false, err
		}
	}
	// note not found, creating new
	_, err = client.ContactsCreateNote(runwayv2.Note{
		ContactId: contactId,
		Body:      fmt.Sprintf("Bookings for this contact:\n%s: %s", date.Format("01/02/2006"), link),
	})
	return false, err
}

func RegisterCollection(contactId, invoiceId string, value float64, date time.Time, l models.Location) (wasAlreadyRegistered bool, err error) {
	// getting client
	client, err := svc.NewClientFromId(l.Id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// getting notes
	notes, err := client.ContactsGetAllNotes(contactId)
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
				ContactId: contactId,
				Id:        n.Id,
				Body:      fmt.Sprintf("%s\n%s %.2f: %s", n.Body, date.Format("01/02/2006"), value, link),
			})

			return false, err
		}
	}
	_, err = client.ContactsCreateNote(runwayv2.Note{
		ContactId: contactId,
		Body:      fmt.Sprintf("Collections for this contact:\n%s %.2f: %s", date.Format("01/02/2006"), value, link),
	})
	return false, err
}
