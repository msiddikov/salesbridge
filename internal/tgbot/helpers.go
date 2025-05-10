package tgbot

import (
	"client-runaway-zenoti/internal/config"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	topics = make(map[string]int64)
)

var (
	autorizedChats = []int64{
		144917817,
	}
	authorizedCommands = []string{
		"/servers",
	}
)

func authorizeChat(opts TgTaskOpts) bool {
	for _, id := range autorizedChats {
		if opts.Id == id {
			return true
		}
	}

	return false
}

func isAuthorizedCommand(simpleOpts TgTaskOpts) bool {
	for _, cmd := range authorizedCommands {
		if simpleOpts.Words[0] == cmd {
			return true
		}
	}
	return false
}

func (svc *Service) CreateTopic(name string) (int64, error) {
	topic, err := svc.bot.CreateForumTopic(svc.Opts.MainGroupId, name, &gotgbot.CreateForumTopicOpts{})
	if err != nil {
		return 0, err
	}
	topics[name] = topic.MessageThreadId
	svc.saveTopics()

	return topic.MessageThreadId, nil
}

func (svc *Service) ChangeTopicName(threadId int64, name string) error {
	_, err := svc.bot.EditForumTopic(svc.Opts.MainGroupId, threadId, &gotgbot.EditForumTopicOpts{
		Name: name,
	})
	return err
}

func (svc *Service) Notify(category string, msg string, urgent bool) {

	// if no topic found, create new one
	if topics[category] == 0 {
		_, err := svc.CreateTopic(category)
		if err != nil {
			fmt.Println(err)
		}
	}

	// if urgent, assign to m_siddikov
	if urgent {
		msg = "@m_siddikov " + msg
	}

	_, err := svc.bot.SendMessage(svc.Opts.MainGroupId, msg, &gotgbot.SendMessageOpts{
		MessageThreadId: topics[category],
	})
	if err != nil {
		fmt.Println(err)
	}
}

func Notify(category string, msg string, urgent bool) {

	category = config.Confs.DB.User + "_" + category
	svc, err := NewTestService()
	if err != nil {
		fmt.Println(err)
	}

	svc.Notify(category, msg, urgent)
}

func NewTestService() (*Service, error) {
	return NewService(TgBotOpts{
		Token:       "8038813737:AAEbpYFV4LiakqOJlcx2WwyCsUHkWo6Dxa8",
		MainGroupId: -1002445791126,
	})
}
