package main

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestUnknownMessage(t *testing.T) {
	message := "qwertyu"
	var chatID int64 = 1
	mock := &MockBotSender{}
	HandleUpdate(mock, message, chatID)
	if mock.text != UnknownCommandResponse {
		t.Errorf("Wrong response expected %s received %s", UnknownCommandResponse, mock.text)
	}
}

type MockBotSender struct {
	text string
}

func (t *MockBotSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if config, ok := c.(tgbotapi.MessageConfig); ok == true {
		t.text = config.Text
	}
	return tgbotapi.Message{}, fmt.Errorf("Test")
}