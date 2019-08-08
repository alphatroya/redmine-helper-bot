package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(opt)
	_, err = redisClient.Ping().Result()
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

	clientManager := redmine.NewClientManager(&http.Client{})

	handler := UpdateHandler{bot, redisClient, clientManager}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		handler.Handle(update.Message.Text, update.Message.Chat.ID)
	}
}
