package svc_jpmreport

import (
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
)

func matchApptGuestsToOpportunities(appt []zenotiv1.Appointment, ops []runwayv2.Opportunity) map[string]runwayv2.Opportunity {
	// Match appointment guests to opportunities
	matches := make(map[string]runwayv2.Opportunity)
	for _, a := range appt {
		_, ok := matches[a.Guest.Id]
		if ok {
			continue
		}

		matchByPhone := runwayv2.Opportunity{}
		matchByEmail := runwayv2.Opportunity{}

		for _, o := range ops {
			// Email and Phone
			if a.Guest.Email != "" && o.Contact.Email != "" && a.Guest.Email == o.Contact.Email && a.Guest.Mobile.Number != "" && o.Contact.Phone != "" && a.Guest.Mobile.Number == o.Contact.Phone {
				matches[a.Guest.Id] = o
				continue
			}

			// Email
			if a.Guest.Email != "" && o.Contact.Email != "" && a.Guest.Email == o.Contact.Email {
				matchByEmail = o
			}

			// Phone
			if a.Guest.Mobile.Number != "" && o.Contact.Phone != "" && a.Guest.Mobile.Number == o.Contact.Phone {
				matchByPhone = o
			}
		}

		if matchByPhone.Id != "" {
			matches[a.Guest.Id] = matchByPhone
			continue
		}

		if matchByEmail.Id != "" {
			matches[a.Guest.Id] = matchByEmail
		}

	}

	return matches
}

func getSourceByGuestId(guestId string, matches map[string]runwayv2.Opportunity) string {
	op, ok := matches[guestId]

	if !ok {
		return "Other"
	}

	return op.Source
}
