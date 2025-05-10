package tgbot

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func (svc *Service) setHandlers(dispatcher *ext.Dispatcher) {
	dispatcher.AddHandler(handlers.NewCommand("id", svc.id))
	dispatcher.AddHandler(handlers.NewCommand("turnOffUpd", svc.stopUpdates))

	dispatcher.AddHandler(handlers.NewMessage(messageIdFilter, svc.NewMessage))
}

func (svc *Service) id(b *gotgbot.Bot, ctx *ext.Context) error {
	err := svc.SendString([]int64{ctx.EffectiveChat.Id}, fmt.Sprintf("Your chat id is %d", ctx.EffectiveChat.Id), 0)
	return err
}

func messageIdFilter(msg *gotgbot.Message) bool {
	return true
}

func (svc *Service) log(b *gotgbot.Bot, ctx *ext.Context) error {
	fmt.Println(ctx.EffectiveMessage)
	return nil
}

func (svc *Service) NewMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

func (svc *Service) stopUpdates(b *gotgbot.Bot, ctx *ext.Context) error {
	svc.updater.Stop()
	return nil
}
