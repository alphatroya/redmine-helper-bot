package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func HandleUpdate(bot BotSender, message string, chatID int64, redisClient redis.Cmdable, client HTTPClient) {
	if strings.HasPrefix(message, "/token") {
		bot.Send(tgbotapi.NewMessage(chatID, handleTokenMessage(message, redisClient, chatID)))
	} else if strings.HasPrefix(message, "/host") {
		message, err := handleHostMessage(message, redisClient, chatID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, message))
		}
	} else if strings.HasPrefix(message, "/fillhours") {
		message, err := HandleFillMessage(message, chatID, redisClient, client)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			telegramMessage := tgbotapi.NewMessage(chatID, message)
			telegramMessage.ParseMode = "Markdown"
			bot.Send(telegramMessage)
		}
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, UnknownCommandResponse))
	}
}

func handleTokenMessage(message string, redisClient redis.Cmdable, chatID int64) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return WrongTokenMessageResponse
	}
	redisClient.Set(fmt.Sprint(chatID)+"_token", splittedMessage[1], 0)
	return SuccessTokenMessageResponse
}

func handleHostMessage(message string, redisClient redis.Cmdable, chatID int64) (string, error) {
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
