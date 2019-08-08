package main

import (
	"github.com/alphatroya/redmine-helper-bot/storage"
	"log"
	"net/http"
	"os"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	storageURLEnvKey     = "REDIS_URL"
	telegramBotKeyEnvKey = "TELEGRAM_BOT_KEY"
)

func main() {
	redisClient, err := storage.NewStorageInstance(storageURLEnvKey)
	if err != nil {
		log.Panicf("Storage configuration failed with error: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv(telegramBotKeyEnvKey))
	if err != nil {
		log.Panicf("Connection to telegram bot is broken, error: %s", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	port := os.Getenv("PORT")
	log.Printf("Port value %s", port)
	if _, err = bot.SetWebhook(tgbotapi.NewWebhook("https://alphatroya-telegram-bot.herokuapp.com:443/" + bot.Token));
		err != nil {
		log.Panicf("Webhook setup failed with error; %s", err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %v", info)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+port, nil)

	clientManager := redmine.NewClientManager(&http.Client{})
	handler := UpdateHandler{bot, redisClient, clientManager}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		handler.Handle(update.Message.Command(), update.Message.CommandArguments(), update.Message.Chat.ID)
	}
}
