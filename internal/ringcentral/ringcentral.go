package ringcentral

import (
	"bytes"
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/ringcentral/sdk"
	"client-runaway-zenoti/internal/types"
	"client-runaway-zenoti/internal/zenotiLegacy/runaway"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	rc sdk.RestClient
)

func GetMessages(from, to time.Time) types.MessageListResponse {

	query := fmt.Sprintf(
		"/restapi/v1.0/account/~/extension/~/message-store?page=1&perPage=100&dateFrom=%s&dateTo=%s",
		from.UTC().Format("2006-01-02T15:04:05.999Z"),
		to.UTC().Format("2006-01-02T15:04:05.999Z"))
	resStr := rc.Get(query)
	res := types.MessageListResponse{}
	json.Unmarshal(resStr, &res)
	return res
}

func SendMessage(chat models.Chat, msg, managerName string, to []string) error {

	body := sdk.CreateSMSMessage{}
	body.From = sdk.MessageStoreCallerInfoRequest{
		PhoneNumber: config.Confs.RC.Username,
	}

	for _, v := range to {
		body.To = append(body.To, sdk.MessageStoreCallerInfoRequest{
			PhoneNumber: v,
		})
	}

	body.Text = msg

	bodyBytes, _ := json.Marshal(body)
	bytes := rc.Post("/restapi/v1.0/account/~/extension/~/sms", bytes.NewReader(bodyBytes))
	if strings.Contains(string(bytes), "errorCode") {
		return fmt.Errorf("%s", bytes)
	}

	res := types.RCSendMessageRes{}

	json.Unmarshal(bytes, &res)

	for _, v := range to {
		AddMessage(models.ChatMessage{
			Date:        res.CreationTIme,
			Content:     res.Subject,
			ManagerName: managerName,
			ChatId:      chat.ID,
			Phone:       v,
		})
	}

	return nil
}

func auth() {

	rc1 := sdk.RestClient{
		ClientID:     config.Confs.RC.ClientId,
		ClientSecret: config.Confs.RC.ClientSecret,
		Server:       config.Confs.RC.ServerURL,
	}

	rc1.Authorize(sdk.GetTokenRequest{
		GrantType: "password",
		Username:  config.Confs.RC.Username,
		Extension: config.Confs.RC.Extension,
		Password:  config.Confs.RC.Password,
	})

	rc = rc1

}

func Auth() {
	auth()

	go func() {
		for {
			time.Sleep(300 * time.Second)
			auth()
		}
	}()

}

func SyncMessages() {
	now := time.Now()
	messages := GetMessages(config.Confs.RC.SyncDate, now)

	for _, v := range messages.Records {
		messageDate, _ := time.Parse("2006-01-02T15:04:05.999Z", v.CreationTime)

		// getting all the phones that were used for batch sending
		phones := []string{}
		if v.Direction == "Inbound" {
			phones = append(phones, v.From.PhoneNumber)
		} else {
			for _, pn := range v.To {
				phones = append(phones, pn.PhoneNumber)
			}
		}

		for _, p := range phones {
			if p == "" {
				continue
			}
			for _, l := range config.Confs.RC.Locations {
				RWID := config.GetLocationById(l).RWID
				contact, err := runaway.GetContact(types.Contact{
					Phone: p,
				}, config.GetLocationByRWID(RWID), false)
				if err != nil {
					continue
				}

				chat, err := GetChat(models.Chat{
					LocationId:  RWID,
					ContactId:   contact.Id,
					ContactName: contact.FirstName + " " + contact.LastName,
					RcId:        config.Confs.RC.ClientId,
				}, true)

				AddMessage(models.ChatMessage{
					Date:    messageDate,
					Content: v.Subject,
					Inbound: v.Direction == "Inbound",
					ChatId:  chat.ID,
					Phone:   p,
				})
			}
		}
	}

	config.Confs.RC.SyncDate = now
}

// DB stuff

func GetChat(chat models.Chat, createIfNotExist bool) (models.Chat, error) {
	res := models.Chat{}

	search := models.Chat{
		LocationId: chat.LocationId,
		ContactId:  chat.ContactId,
	}
	err := db.DB.Where(&search).First(&res).Error

	if errors.Is(err, gorm.ErrRecordNotFound) && createIfNotExist {
		err = db.DB.Create(&chat).Error
		if err == nil {
			return chat, nil
		}
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func AddMessage(message models.ChatMessage) error {
	res := models.ChatMessage{}
	search := models.ChatMessage{
		Content: message.Content,
		Date:    message.Date,
	}

	err := db.DB.Where(&search).First(&res).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.DB.Create(&message).Error
	}

	return err
}

func GetChatsWithLastMessage(locationId string) []types.ChatRes {
	chats := []types.ChatRes{}

	db.DB.Raw(`select c.id as chat_id,c.contact_name, msg.content as last_message, msg.date, c.contact_id, c.location_id
	from (select r.id, r.contact_name, r.contact_id, r.location_id
		 from chats as r
		 where r.location_id=?) as c 
	left join lateral
		(select r.content, r.date, r.chat_id
		from chat_messages as r
		where r.chat_id = c.id
		order by r.date desc
		limit 1) as msg on c.id=msg.chat_id`, locationId).Scan(&chats)

	return chats
}

func GetChatMessages(chatId uint) []types.ChatMessagesRes {
	msgs := []types.ChatMessagesRes{}

	db.DB.Raw(`SELECT date, id as message_id, content as text,  manager_name,  inbound
	FROM chat_messages
    where chat_id=?`, chatId).Scan(&msgs)

	return msgs
}
