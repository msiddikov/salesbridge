package svc_jpmreport

import (
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
)

func breakBySource(leads []runwayv2.Opportunity, appt []zenotiv1.Appointment, collections []zenotiv1.Collection) map[string]SourceBreakdown {
	matches := matchApptGuestsToOpportunities(appt, leads)

	res := make(map[string]SourceBreakdown)

	for _, lead := range leads {
		if lead.Source == "" {
			lead.Source = "Other"
		}
		sb, ok := res[lead.Source]
		if !ok {
			sb = SourceBreakdown{}
		}

		sb.Leads++
		sb.Source = lead.Source
		res[lead.Source] = sb
	}

	for _, c := range collections {
		source := getSourceByGuestId(c.Guest.Id, matches)
		sb, ok := res[source]
		if !ok {
			sb = SourceBreakdown{}
		}

		if c.Total_collection == 0 {
			sb.Consultations++
		}

		sb.Revenue += c.Total_collection
		sb.Source = source
		res[source] = sb
	}

	return res
}
