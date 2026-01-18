package runwayv2

import "fmt"

type (
	MessagesFilter struct {
		LocationId     string `json:"locationId"`
		Channel        string `json:"channel"`
		Limit          int    `json:"limit"`
		Cursor         string `json:"cursor"`
		SortBy         string `json:"sortBy"`
		SortOrder      string `json:"sortOrder"`
		ConversationId string `json:"conversationId"`
		ContactId      string `json:"contactId"`
		StartDate      string `json:"startDate"`
		EndDate        string `json:"endDate"`
	}

	MessagesExportResponse struct {
		Messages   []Message `json:"messages"`
		NextCursor string    `json:"nextCursor"`
		Total      int       `json:"total"`
	}

	Transcription struct {
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		Transcript string `json:"transcript"`
		Confidence string `json:"confidence"`
	}
)

func (a *Client) MessagesExport(f MessagesFilter) (MessagesExportResponse, error) {
	res := MessagesExportResponse{}
	f.LocationId = a.cfg.locationId

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/conversations/messages/export",
		QParams:  f.toQueryParams(),
	}, &res)

	return res, err
}

func (a *Client) MessagesGetRecording(messageId string) ([]byte, error) {

	_, res, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/conversations/messages/" + messageId + "/locations/" + a.cfg.locationId + "/recording",
	}, nil)
	return res, err
}

func (a *Client) MessagesGetTranscription(messageId string) ([]Transcription, error) {

	res := []Transcription{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/conversations" + "/locations/" + a.cfg.locationId + "/messages/" + messageId + "/transcription",
	}, &res)

	return res, err
}

func (mf *MessagesFilter) toQueryParams() []queryParam {
	params := []queryParam{}

	if mf.LocationId != "" {
		params = append(params, queryParam{Key: "locationId", Value: mf.LocationId})
	}
	if mf.Channel != "" {
		params = append(params, queryParam{Key: "channel", Value: mf.Channel})
	}
	if mf.Limit != 0 {
		params = append(params, queryParam{Key: "limit", Value: fmt.Sprintf("%d", mf.Limit)})
	}
	if mf.Cursor != "" {
		params = append(params, queryParam{Key: "cursor", Value: mf.Cursor})
	}
	if mf.SortBy != "" {
		params = append(params, queryParam{Key: "sortBy", Value: mf.SortBy})
	}
	if mf.SortOrder != "" {
		params = append(params, queryParam{Key: "sortOrder", Value: mf.SortOrder})
	}
	if mf.ConversationId != "" {
		params = append(params, queryParam{Key: "conversationId", Value: mf.ConversationId})
	}
	if mf.ContactId != "" {
		params = append(params, queryParam{Key: "contactId", Value: mf.ContactId})
	}
	if mf.StartDate != "" {
		params = append(params, queryParam{Key: "startDate", Value: mf.StartDate})
	}
	if mf.EndDate != "" {
		params = append(params, queryParam{Key: "endDate", Value: mf.EndDate})
	}

	return params
}
