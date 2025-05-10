package tgbot

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type (
	TgTaskOpts struct {
		Id       int64
		ReplyId  int
		Words    []string
		StopCh   chan bool
		Username string
		Title    string
	}

	TgBotOpts struct {
		Token       string
		MainGroupId int64
	}

	Service struct {
		Opts    TgBotOpts
		bot     *gotgbot.Bot
		updater *ext.Updater
	}
)

func NewService(opts TgBotOpts) (*Service, error) {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(opts.Token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout, // Customise the default request timeout here
				APIURL:  gotgbot.DefaultAPIURL,  // As well as the Default API URL here (in case of using local bot API servers)
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}

	svc := &Service{
		Opts: opts,
		bot:  b,
	}
	svc.loadTopics()

	return svc, nil
}

func (svc *Service) StartPolling() {
	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	// Add echo handler to reply to all text messages.
	svc.setHandlers(dispatcher)

	svc.updater = updater
	// Start receiving updates.
	err := updater.StartPolling(svc.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", svc.bot.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func (svc *Service) SendString(ids []int64, s string, replyId int64) error {
	if svc.bot == nil {
		lvn.Logger.Errorf("bot is offline: %s", s)
		return nil
	}
	for _, id := range ids {
		opts := &gotgbot.SendMessageOpts{}
		if replyId != 0 {
			opts.ReplyParameters.MessageId = replyId
		}
		opts.ParseMode = "MarkdownV2"

		_, err := svc.bot.SendMessage(id, s, opts)
		if err != nil {
			if strings.Contains(err.Error(), "must be escaped") {
				s = E(s)
				_, err = svc.bot.SendMessage(id, s, opts)
				if err != nil {
					fmt.Println(err)
					return err
				}
				continue
			}
			return err
		}
	}

	return nil
}

func E(si interface{}) (r string) {
	s := fmt.Sprint(si)
	var escapeChars = []string{
		"\\", "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
	}

	// escape the string
	for _, v := range escapeChars {
		s = strings.Replace(s, v, "\\"+v, -1)
	}

	return s
}
