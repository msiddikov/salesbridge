package runwayv2

import (
	"encoding/json"
)

type (
	NewMessageParams struct {
		Type                   string `json:"type"`
		Message                string `json:"message"`
		ConversationId         string `json:"conversationId"`
		ConversationProviderId string `json:"conversationProviderId"`
	}

	NewMessageRes struct {
		ConverstionId string
		MessageId     string
		MessageIds    []string
		Msg           string
	}
)

func (a *Client) ConversationsAddInboundMessage(p NewMessageParams) (NewMessageRes, error) {
	body, err := json.Marshal(p)
	res := NewMessageRes{}
	if err != nil {
		return res, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/conversations/messages/inbound",
		Body:     string(body),
	}, &res)

	return res, err
}

func (a *Client) ConversationsFind(contactId string) ([]Conversation, error) {
	res := struct {
		Conversations []Conversation
	}{}

	_, _, err := a.fetch(reqParams{
		Method:   "GET",
		Endpoint: "/conversations/search",
		QParams: []queryParam{
			{
				Key:   "locationId",
				Value: a.cfg.locationId,
			},
			{
				Key:   "contactId",
				Value: contactId,
			},
		},
	}, &res)
	return res.Conversations, err
}

func (a *Client) ConversationsCreate(contactId string) (Conversation, error) {
	bd := struct {
		LocationId string `json:"locationId"`
		ContactId  string `json:"contactId"`
	}{
		LocationId: a.cfg.locationId,
		ContactId:  contactId,
	}
	body, err := json.Marshal(bd)
	res := struct {
		Conversation Conversation
	}{}
	if err != nil {
		return res.Conversation, err
	}
	_, _, err = a.fetch(reqParams{
		Method:   "POST",
		Endpoint: "/conversations/",
		Body:     string(body),
	}, &res)

	return res.Conversation, err
}
