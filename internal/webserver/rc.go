package webServer

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/ringcentral"
	"client-runaway-zenoti/internal/types"
	"client-runaway-zenoti/internal/zenotiLegacy/runaway"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setChatRoutes(router *gin.Engine) {
	router.GET("/rc/contacts", contactsLookUp)
	router.GET("/rc/contacts/:contactId/:locationId", contactInfo)
	router.GET("/rc/chats/:locationId", chats)
	router.GET("/rc/messages/:chatId", messages)
	router.POST("/rc/message", sendMessage)
	router.GET("/rc/update", update)
}

func contactsLookUp(c *gin.Context) {
	query, _ := c.GetQuery("query")
	locationId, _ := c.GetQuery("location")

	contacts := runaway.QueryContacts(query, config.GetLocationByRWID(locationId))

	c.JSON(200, contacts)
}

func contactInfo(c *gin.Context) {
	contactId, _ := c.Params.Get("contactId")
	locationId, _ := c.Params.Get("locationId")

	contact, _ := runaway.GetContactById(contactId, config.GetLocationByRWID(locationId))

	c.JSON(200, contact)
}

func chats(c *gin.Context) {
	locationId, _ := c.Params.Get("locationId")

	chats := ringcentral.GetChatsWithLastMessage(locationId)

	c.JSON(200, chats)
}

func messages(c *gin.Context) {
	chatIdText, _ := c.Params.Get("chatId")
	msg := []types.ChatMessagesRes{}
	chatId, err := strconv.ParseUint(chatIdText, 0, 16)
	if err != nil {
		c.JSON(200, msg)
	}
	msg = ringcentral.GetChatMessages(uint(chatId))

	c.JSON(200, msg)
}

func sendMessage(c *gin.Context) {

	bodyBytes, _ := io.ReadAll(c.Request.Body)

	body := types.SendMessageBody{}
	json.Unmarshal(bodyBytes, &body)

	l := config.GetLocationByRWID(body.LocationId)
	contact, err := runaway.GetContactById(body.ContactId, l)
	if err != nil {
		c.Writer.WriteHeader(500)
		return
	}
	chat, _ := ringcentral.GetChat(models.Chat{
		LocationId:  body.LocationId,
		ContactId:   contact.Id,
		ContactName: contact.FirstName + " " + contact.LastName,
		PhoneNo:     contact.Phone,
		RcId:        config.Confs.RC.ClientId,
	}, true)
	err = ringcentral.SendMessage(chat, body.Text, body.ManagerName, []string{contact.Phone})

	if err != nil {

		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	c.JSON(200, chat)
}

func update(c *gin.Context) {
	ringcentral.SyncMessages()
	c.Writer.WriteHeader(200)
}
