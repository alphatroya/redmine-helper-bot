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

var commandHandlers map[int64]commands.Command

func init() {
	commandHandlers = make(map[int64]commands.Command)
}

func (t *UpdateHandler) Handle(command string, message string, chatID int64) {
	commandBuilder := commands.NewBotCommandsBuilder(t.storage, t.client)
	commandHandler := commandBuilder.Build(command, message, chatID)
	result, err := commandHandler.Handle(message)
	var newMessage tgbotapi.MessageConfig
	if err != nil {
		newMessage = tgbotapi.NewMessage(chatID, err.Error())
	} else {
		newMessage = tgbotapi.NewMessage(chatID, result.Message())
	}
	newMessage.ParseMode = tgbotapi.ModeMarkdown
	commandHandlers[chatID] = commandHandler

	_, err = t.bot.Send(newMessage)
	if err != nil {
		log.Printf("error during send operation, got: %s", err)
	}
}

func (t *UpdateHandler) HandleMessage(message string, chatID int64) {
	var result *commands.CommandResult
	var err error
	commandHandler, ok := commandHandlers[chatID]
	if !ok || commandHandler == nil || commandHandler.IsCompleted() {
		result, err = commands.NewUnknownCommand().Handle(message)
	} else {
		result, err = commandHandler.Handle(message)
	}
	var newMessage tgbotapi.MessageConfig
	if err != nil {
		newMessage = tgbotapi.NewMessage(chatID, err.Error())
	} else {
		newMessage = tgbotapi.NewMessage(chatID, result.Message())
	}
	newMessage.ParseMode = tgbotapi.ModeMarkdown

	_, err = t.bot.Send(newMessage)
	if err != nil {
		log.Printf("error during send operation, got: %s", err)
	}
}
