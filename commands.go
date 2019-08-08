package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/storage"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type UpdateHandler struct {
	bot         BotSender
	redisClient storage.Manager
	client      redmine.Client
}

func (t *UpdateHandler) Handle(message string, chatID int64) {
	if strings.HasPrefix(message, "/token") {
		t.bot.Send(tgbotapi.NewMessage(chatID, t.handleTokenMessage(message, t.redisClient, chatID)))
	} else if strings.HasPrefix(message, "/host") {
		message, err := t.handleHostMessage(message, t.redisClient, chatID)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			t.bot.Send(tgbotapi.NewMessage(chatID, message))
		}
	} else if strings.HasPrefix(message, "/fillhours") {
		message, err := t.handleFillMessage(message, chatID, t.redisClient, t.client)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			telegramMessage := tgbotapi.NewMessage(chatID, message)
			telegramMessage.ParseMode = "Markdown"
			t.bot.Send(telegramMessage)
		}
	} else {
		t.bot.Send(tgbotapi.NewMessage(chatID, UnknownCommandResponse))
	}
}

func (t *UpdateHandler) handleTokenMessage(message string, redisClient storage.Manager, chatID int64) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return WrongTokenMessageResponse
	}
	redisClient.Set(fmt.Sprint(chatID)+"_token", splittedMessage[1], 0)
	return SuccessTokenMessageResponse
}

func (t *UpdateHandler) handleHostMessage(message string, redisClient storage.Manager, chatID int64) (string, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return "", fmt.Errorf(WrongHostMessageResponse)
	}
	_, err := url.ParseRequestURI(splittedMessage[1])
	if err != nil {
		return "", err
	}
	redisClient.Set(fmt.Sprint(chatID)+"_host", splittedMessage[1], 0)
	return SuccessHostMessageResponse, nil
}

func (t *UpdateHandler) handleFillMessage(message string, chatID int64, redisClient storage.Manager, client redmine.Client) (string, error) {
	chatIDString := fmt.Sprint(chatID)

	token, err := redisClient.Get(chatIDString + "_token").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	client.SetToken(token)

	host, err := redisClient.Get(chatIDString + "_host").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	client.SetHost(host)

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf(WrongFillHoursWrongNumberOfArgumentsResponse)
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	issueID := splitted[1]
	if !regex.MatchString(issueID) {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}

	_, conversionError := strconv.ParseFloat(splitted[2], 32)
	if conversionError != nil {
		return "", fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}

	requestBody, err := client.FillHoursRequest(issueID, splitted[2], strings.Join(splitted[3:], " "))
	if err != nil {
		return "", err
	}

	issue, _ := client.Issue(issueID)

	return SuccessFillHoursMessageResponse(requestBody.TimeEntry.ID, issue, requestBody.TimeEntry.Hours, host), nil
}
