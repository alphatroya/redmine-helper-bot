package commands

import (
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type Builder interface {
	Build(command string, message string, chatID int64) Command
}

type BotCommandsBuilder struct {
	storage       storage.Manager
	redmineClient redmine.Client
}

func NewBotCommandsBuilder(storage storage.Manager, redmineClient redmine.Client) *BotCommandsBuilder {
	return &BotCommandsBuilder{storage: storage, redmineClient: redmineClient}
}

func (b BotCommandsBuilder) Build(command string, message string, chatID int64) Command {
	switch command {
	case "token":
		return newSetTokenCommand(b.storage, chatID)
	case "host":
		return newSetHostCommand(b.storage, chatID)
	case "fillhours":
		return newPartlyFillHoursCommand(b.redmineClient, b.storage, chatID)
	case "activities":
		return newActivitiesCommand(b.redmineClient, b.storage, chatID)
	case "start":
		return newIntroCommand()
	case "stop":
		return newStopCommand()
	default:
		return NewUnknownCommand()
	}
}
