package types

import (
	"client-runaway-zenoti/internal/ringcentral/sdk"
	"time"
)

type (
	ChatRes struct {
		LocationId  string    `json:"locationId"`
		ChatId      uint      `json:"chatId"`
		ContactId   string    `json:"contactId"`
		ContactName string    `json:"contactName"`
		LastMessage string    `json:"lastMessage"`
		Date        time.Time `json:"date"`
	}

	ChatMessagesRes struct {
		Date        time.Time `json:"date"`
		MessageId   uint      `json:"messageId"`
		Text        string    `json:"text"`
		ManagerName string    `json:"managerName"`
		Inbound     bool      `json:"inbound"`
	}

	MessageListResponse struct {
		Records []MessageListMessageResponse `json:"records"`
		Paging  sdk.PagingResource           `json:"paging"`
	}

	SendMessageBody struct {
		ContactId   string
		LocationId  string
		Text        string
		ManagerName string
	}

	MessageListMessageResponse struct {
		Id      int    `json:"id"`
		BatchId string `json:"batchId"`
		From    struct {
			PhoneNumber string `json:"phoneNumber"`
		} `json:"from"`
		To []struct {
			PhoneNumber string `json:"phoneNumber"`
		} `json:"to"`
		CreationTime     string  `json:"creationTime"`
		LastModifiedTime string  `json:"lastModifiedTime"`
		MessageStatus    string  `json:"messageStatus"`
		SegmentCount     int     `json:"segmentCount"`
		Subject          string  `json:"subject"`
		Cost             float64 `json:"cost"`
		Direction        string  `json:"direction"`
		ErrorCode        string  `json:"errorCode"`
	}

	RCSendMessageRes struct {
		CreationTIme time.Time
		Subject      string
	}
)
