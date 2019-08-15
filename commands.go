package main

import (
	"log"

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
	commandBuilder := commands.NewBotCommandsBuilder(t.storage, t.client, chatID)
	commandHandler := commandBuilder.Build(command, message, nil)
	message, err := commandHandler.Handle(message)
	if err != nil {
		_, err = t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		log.Printf("got error during send operation, got: %s", err)
		return
	}
	telegramMessage := tgbotapi.NewMessage(chatID, message)
	telegramMessage.ParseMode = "Markdown"
	_, err = t.bot.Send(telegramMessage)
	log.Printf("got error during send operation, got: %s", err)
}
