package svc_jpmreport

func breakBySource(leadsData []LeadRawData) map[string]SourceBreakdown {

	res := make(map[string]SourceBreakdown)

	for _, lead := range leadsData {
		if lead.Source == "" {
			lead.Source = "Other"
		}
		sb, ok := res[lead.Source]
		if !ok {
			sb = SourceBreakdown{}
		}

		sb.Leads++
		sb.Source = lead.Source
		if lead.HadConsultation {
			sb.Consultations++
		}
		sb.Revenue += lead.Revenue
		res[lead.Source] = sb
	}

	return res
}
