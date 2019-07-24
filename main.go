package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
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
		HandleUpdate(bot, update)
	}
}

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	messageString := update.Message.Text
	if strings.HasPrefix(messageString, "/token") {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HandleTokenMessage(messageString, redisClient, update.Message.Chat.ID)))
	} else if strings.HasPrefix(messageString, "/host") {
		message, err := HandleHostMessage(messageString, redisClient, update.Message.Chat.ID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	} else if strings.HasPrefix(messageString, "/fillhours") {
		message, err := HandleFillMessage(messageString, update, redisClient)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	} else {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введена неправильная команда"))
	}
}
