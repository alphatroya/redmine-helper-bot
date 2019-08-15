package commands

import (
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type Builder interface {
	Build(command string, message string, previousCommand Command) Command
}

type BotCommandsBuilder struct {
	storage       storage.Manager
	redmineClient redmine.Client
	chatID        int64
}

func NewBotCommandsBuilder(storage storage.Manager, redmineClient redmine.Client, chatID int64) *BotCommandsBuilder {
	return &BotCommandsBuilder{storage: storage, redmineClient: redmineClient, chatID: chatID}
}

func (b BotCommandsBuilder) Build(command string, message string, previousCommand Command) Command {
	switch command {
	case "token":
		return newSetTokenCommand(b.storage, b.chatID)
	case "host":
		return newSetHostCommand(b.storage, b.chatID)
	case "fillhours":
		return newFillHoursCommand(b.storage, b.chatID, b.redmineClient)
	default:
		return newUnknownCommand()
	}
}
