package main

import (
	"log"
	"sync"

	"github.com/alphatroya/redmine-helper-bot/commands"

	"github.com/alphatroya/redmine-helper-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type UpdateHandler struct {
	bot     BotSender
	storage storage.Manager
}

var commandHandlers = struct {
	sync.RWMutex
	handlers map[int64]commands.Command
}{handlers: make(map[int64]commands.Command)}

func (t *UpdateHandler) Handle(command string, message string, chatID int64) {
	commandBuilder := commands.NewBotCommandsBuilder(t.storage)
	commandHandler := commandBuilder.Build(command, message, chatID)
	result, err := commandHandler.Handle(message)
	commandHandlers.RLock()
	commandHandlers.handlers[chatID] = commandHandler
	commandHandlers.RUnlock()
	t.sendMessage(chatID, result, err)
}

func (t *UpdateHandler) HandleMessage(message string, chatID int64) {
	var result *commands.CommandResult
	var err error
	commandHandlers.RLock()
	commandHandler, ok := commandHandlers.handlers[chatID]
	commandHandlers.RUnlock()
	if !ok || commandHandler == nil || commandHandler.IsCompleted() {
		result, err = commands.NewUnknownCommand().Handle(message)
	} else {
		result, err = commandHandler.Handle(message)
	}
	t.sendMessage(chatID, result, err)
}

func (t *UpdateHandler) sendMessage(chatID int64, result *commands.CommandResult, err error) {
	if err != nil {
		newMessage := tgbotapi.NewMessage(chatID, err.Error())
		newMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		newMessage.ParseMode = tgbotapi.ModeMarkdown
		_, err = t.bot.Send(newMessage)
		if err != nil {
			log.Printf("error during send operation, got: %s", err)
		}
	} else {
		for i, message := range result.Messages() {
			newMessage := tgbotapi.NewMessage(chatID, message)
			buttons := result.Buttons()
			if i == len(result.Messages())-1 && len(buttons) != 0 {
				var rows [][]tgbotapi.KeyboardButton
				var keyboards []tgbotapi.KeyboardButton
				for _, button := range buttons {
					keyboards = append(keyboards, tgbotapi.NewKeyboardButton(button))
					if len(keyboards) == 2 {
						rows = append(rows, tgbotapi.NewKeyboardButtonRow(keyboards...))
						keyboards = []tgbotapi.KeyboardButton{}
					}
				}
				newMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)
			} else {
				newMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
			newMessage.ParseMode = tgbotapi.ModeMarkdown
			_, err = t.bot.Send(newMessage)
			if err != nil {
				log.Printf("error during send operation, got: %s", err)
			}
		}
	}
}
