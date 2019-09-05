package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alphatroya/redmine-helper-bot/storage"
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

	log.Printf("Authorized on account %s", bot.Self.UserName)

	handler := UpdateHandler{bot, redisClient}

	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
		configureLongPolling(handler, bot)
	} else {
		configureWebHookObserving(handler, bot)
	}
}

func configureLongPolling(handler UpdateHandler, bot *tgbotapi.BotAPI) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panicf("Failed to obtain updates channel, error: %s", err)
	}
	handleUpdates(updates, handler)
}

func configureWebHookObserving(updateHandler UpdateHandler, bot *tgbotapi.BotAPI) {
	port := os.Getenv("PORT")
	log.Printf("Port value %s", port)
	if _, err := bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("SERVER_URL") + ":443/" + bot.Token)); err != nil {
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
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	handleUpdates(updates, updateHandler)
}

func handleUpdates(updates tgbotapi.UpdatesChannel, handler UpdateHandler) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			go handler.Handle(update.Message.Command(), update.Message.CommandArguments(), update.Message.Chat.ID)
		} else {
			go handler.HandleMessage(update.Message.Text, update.Message.Chat.ID)
		}
	}
}
