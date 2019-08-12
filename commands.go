package main

import (
	"github.com/alphatroya/redmine-helper-bot/commands"

	"github.com/alphatroya/redmine-helper-bot/storage"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type UpdateHandler struct {
	bot     BotSender
	storage storage.Manager
	client  redmine.Client
}

func (t *UpdateHandler) Handle(command string, message string, chatID int64) {
	switch command {
	case "token":
		command := commands.NewSetTokenCommand(t.storage, chatID)
		message, err := command.Handle(message)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
			return
		}
		t.bot.Send(tgbotapi.NewMessage(chatID, message))
	case "host":
		command := commands.NewSetHostCommand(t.storage, chatID)
		message, err := command.Handle(message)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
			return
		}
		t.bot.Send(tgbotapi.NewMessage(chatID, message))
	case "fillhours":
		command := commands.NewFillHoursCommand(t.storage, chatID, t.client)
		message, err := command.Handle(message)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
			return
		}
		telegramMessage := tgbotapi.NewMessage(chatID, message)
		telegramMessage.ParseMode = "Markdown"
		t.bot.Send(telegramMessage)
	default:
		t.bot.Send(tgbotapi.NewMessage(chatID, UnknownCommandResponse))
	}
}
