package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type UpdateHandler struct {
	bot         BotSender
	redisClient redis.Cmdable
	client      redmine.HTTPClient
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

func (t *UpdateHandler) handleTokenMessage(message string, redisClient redis.Cmdable, chatID int64) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return WrongTokenMessageResponse
	}
	redisClient.Set(fmt.Sprint(chatID)+"_token", splittedMessage[1], 0)
	return SuccessTokenMessageResponse
}

func (t *UpdateHandler) handleHostMessage(message string, redisClient redis.Cmdable, chatID int64) (string, error) {
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

func (t *UpdateHandler) handleFillMessage(message string, chatID int64, redisClient redis.Cmdable, client redmine.HTTPClient) (string, error) {
	chatIDString := fmt.Sprint(chatID)

	token, err := redisClient.Get(chatIDString + "_token").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}

	host, err := redisClient.Get(chatIDString + "_host").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf(WrongFillHoursWrongNumberOfArgumentsResponse)
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	if !regex.MatchString(splitted[1]) {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}

	_, conversionError := strconv.ParseFloat(splitted[2], 32)
	if conversionError != nil {
		return "", fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}

	requestBody, err := redmine.FillHoursRequest(token, host, splitted, client)
	if err != nil {
		return "", err
	}
	resultMessage := SuccessFillHoursMessageResponse(requestBody.TimeEntry.IssueID, requestBody.TimeEntry.Hours, host)
	return resultMessage, nil
}
