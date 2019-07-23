package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var tokens = make(map[int64]string)
var hosts = make(map[int64]string)
var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	pong, err := redisClient.Ping().Result()
	log.Println(pong, err)

	apiKey := os.Getenv("TELEGRAM_BOT_KEY")
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
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
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HandleTokenMessage(messageString, tokens, update.Message.Chat.ID)))
	} else if strings.HasPrefix(messageString, "/host") {
		message, err := HandleHostMessage(messageString, hosts, update.Message.Chat.ID)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	} else if strings.HasPrefix(messageString, "/fillhours") {
		message, err := HandleFillMessage(messageString, update, tokens, hosts)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	} else {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введена неправильная команда"))
	}
}
