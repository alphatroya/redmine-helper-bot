package main

import (
	"fmt"
	"testing"
)

func TestTokenRequest(t *testing.T) {
	data := []struct {
		command  string
		chatID   int64
		expected string
	}{
		{"/token", 1, WrongTokenMessageResponse},
		{"/token test test", 1, WrongTokenMessageResponse},
		{"/token fdsjfdsj", 1, SuccessTokenMessageResponse},
		{"qwertyu", 1, UnknownCommandResponse},
		{"/host", 1, WrongHostMessageResponse},
		{"/host test test", 1, WrongHostMessageResponse},
		{"/host test", 1, "parse test: invalid URI for request"},
		{"/host https://www.google.com", 1, SuccessHostMessageResponse},
	}

	for _, message := range data {
		mock := &MockBotSender{}
		redisMock := NewRedisMock()
		HandleUpdate(mock, message.command, message.chatID, redisMock)
		if mock.text != message.expected {
			t.Errorf("Wrong response expected: %s received: %s", message.expected, mock.text)
		}
	}
}

func TestStorageTokenData(t *testing.T) {
	data := []struct {
		command string
		chatID  int64
	}{
		{"431", 44},
		{"23", 45},
	}

	for _, message := range data {
		botMock := &MockBotSender{}
		redisMock := NewRedisMock()
		HandleUpdate(botMock, "/token "+message.command, message.chatID, redisMock)
		tokenValue := redisMock.Get(fmt.Sprint(message.chatID) + "_token").Val()
		if tokenValue != message.command {
			t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, message.command)
		}
	}
}

func TestMultipleRequestStorageTokenData(t *testing.T) {
	command := "4433"
	var chatID int64 = 1

	command2 := "4433"
	var chatID2 int64 = 2

	redisMock := NewRedisMock()
	botMock := &MockBotSender{}

	HandleUpdate(botMock, "/token "+command, chatID, redisMock)
	HandleUpdate(botMock, "/token "+command2, chatID2, redisMock)

	tokenValue := redisMock.Get(fmt.Sprint(chatID2) + "_token").Val()
	if tokenValue != command2 {
		t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, command2)
	}
}

func TestHandleHostMessageWithCorrectCommand(t *testing.T) {
	data := []struct {
		url    string
		chatID int64
	}{
		{"https://www.google.com", 44},
		{"https://www.tt.com", 45},
		{"https://tt.com", 46},
	}

	for _, message := range data {
		botMock := &MockBotSender{}
		redisMock := NewRedisMock()
		HandleUpdate(botMock, "/host"+" "+message.url, message.chatID, redisMock)
		hostValue := redisMock.Get(fmt.Sprint(message.chatID) + "_host").Val()
		if hostValue != message.url {
			t.Errorf("Wrong saved host value %s is not %s", hostValue, message.url)
		}
	}
}
