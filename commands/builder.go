package commands

import (
	"net/http"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type Builder interface {
	Build(command string, message string, chatID int64) Command
}

type BotCommandsBuilder struct {
	storage storage.Manager
	printer Printer
}

func NewBotCommandsBuilder(storage storage.Manager) *BotCommandsBuilder {
	return &BotCommandsBuilder{storage: storage, printer: redmine.TablePrinter{}}
}

func (b BotCommandsBuilder) Build(command string, message string, chatID int64) Command {
	switch command {
	case "token":
		return newSetTokenCommand(b.storage, chatID)
	case "host":
		return newSetHostCommand(b.storage, chatID)
	case "activity":
		redmineClient := redmine.NewClientManager(&http.Client{}, b.storage, chatID)
		return newActivitiesCommand(redmineClient, b.storage, chatID)
	case "start":
		return newIntroCommand()
	case "stop":
		return newStopCommand(b.storage, chatID)
	case "fhm":
		redmineClient := redmine.NewClientManager(&http.Client{}, b.storage, chatID)
		return NewFillHoursMany(redmineClient, b.storage, chatID)
	case "fh":
		redmineClient := redmine.NewClientManager(&http.Client{}, b.storage, chatID)
		command, err := NewFillHoursCommand(redmineClient, b.printer, b.storage, chatID, message)
		if err != nil {
			return NewUnknownCommandWithMessage(err.Error())
		}
		return command
	case "fstatus":
		redmineClient := redmine.NewClientManager(&http.Client{}, b.storage, chatID)
		return NewFillStatus(redmineClient, chatID)
	default:
		return NewUnknownCommand()
	}
}
