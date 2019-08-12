package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
		message, err := t.handleFillMessage(message, chatID, t.storage, t.client)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			telegramMessage := tgbotapi.NewMessage(chatID, message)
			telegramMessage.ParseMode = "Markdown"
			t.bot.Send(telegramMessage)
		}
	default:
		t.bot.Send(tgbotapi.NewMessage(chatID, UnknownCommandResponse))
	}
}

func (t *UpdateHandler) handleFillMessage(message string, chatID int64, redisClient storage.Manager, client redmine.Client) (string, error) {

	token, err := redisClient.GetToken(chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	client.SetToken(token)

	host, err := redisClient.GetHost(chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	client.SetHost(host)

	splitted := strings.Split(message, " ")
	if len(splitted) < 3 {
		return "", fmt.Errorf(WrongFillHoursWrongNumberOfArgumentsResponse)
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	issueID := splitted[0]
	if !regex.MatchString(issueID) {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}

	_, conversionError := strconv.ParseFloat(splitted[1], 32)
	if conversionError != nil {
		return "", fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}

	requestBody, err := client.FillHoursRequest(issueID, splitted[1], strings.Join(splitted[2:], " "))
	if err != nil {
		return "", err
	}

	issue, _ := client.Issue(issueID)

	return SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, issue, requestBody.TimeEntry.Hours, host), nil
}
