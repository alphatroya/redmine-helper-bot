package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Panicf("Connection to Redis instance is broken: %s", err)
	}

	apiKey := os.Getenv("TELEGRAM_BOT_KEY")
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panicf("Connection to telegram bot is broken: %s", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		HandleUpdate(bot, update.Message.Text, update.Message.Chat.ID, redisClient)
	}
}

func HandleUpdate(bot BotSender, message string, chatID int64, redisClient redis.Cmdable) {
	if strings.HasPrefix(message, "/token") {
		bot.Send(tgbotapi.NewMessage(chatID, HandleTokenMessage(message, redisClient, chatID)))
	} else if strings.HasPrefix(message, "/host") {
		message, err := HandleHostMessage(message, redisClient, chatID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(chatID, message))
	} else if strings.HasPrefix(message, "/fillhours") {
		message, err := HandleFillMessage(message, chatID, redisClient, &http.Client{})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		}
		telegramMessage := tgbotapi.NewMessage(chatID, message)
		telegramMessage.ParseMode = "Markdown"
		bot.Send(telegramMessage)
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, UnknownCommandResponse))
	}
}

type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
