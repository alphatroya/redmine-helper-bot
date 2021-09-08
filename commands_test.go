package main

import (
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func setupSubTest(t *testing.T) func(t *testing.T) {
	setup()
	return func(t *testing.T) {
		tearDown()
	}
}

var botMock *MockBotSender
var redisMock *storage.Mock
var handler *UpdateHandler

func setup() {
	botMock = &MockBotSender{}
	redisMock = storage.NewStorageMock()
	handler = &UpdateHandler{botMock, redisMock}
}

func tearDown() {
	botMock = nil
	redisMock = nil
	handler = nil
}

func TestTokenRequest(t *testing.T) {
	const wrongHostArgumentsText = "Неправильное количество аргументов, введите адрес в формате `/host <адрес сервера>`(например, `/host https://google.ru`)"
	const wrongTokenArgumentsCountText = "Неправильное количество аргументов, введите токен доступа к АПИ в формате `/token <токен>`"
	data := []struct {
		command  string
		message  string
		chatID   int64
		expected string
	}{
		{"token", "", 1, wrongTokenArgumentsCountText},
		{"token", "test test", 1, wrongTokenArgumentsCountText},
		{"token", "fdsjfdsj", 1, "Токен успешно обновлен"},
		{"token", "  ", 1, wrongTokenArgumentsCountText},
		{"", "qwertyu", 1, "Введена неправильная команда"},
		{"host", "", 1, wrongHostArgumentsText},
		{"host", " ", 1, wrongHostArgumentsText},
		{"host", "test test", 1, wrongHostArgumentsText},
		{"host", "test", 1, "parse \"test\": invalid URI for request"},
		{"host", "https://www.google.com", 1, "Адрес сервера успешно обновлен"},
	}

	for _, message := range data {
		teardownSubTest := setupSubTest(t)

		handler.Handle(message.command, message.message, message.chatID)
		if botMock.text != message.expected {
			t.Errorf("Wrong response command: %s\narguments: %s\nexpected: %s\nreceived: %s", message.command, message.message, message.expected, botMock.text)
		}

		teardownSubTest(t)
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
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		handler.Handle("token", message.command, message.chatID)
		tokenValue := redisMock.StorageToken()[message.chatID]
		if tokenValue != message.command {
			t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, message.command)
		}
	}
}

func TestMultipleRequestStorageTokenData(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	command := "4433"
	var chatID int64 = 1

	command2 := "4433"
	var chatID2 int64 = 2

	handler.Handle("token", command, chatID)
	handler.Handle("token", command2, chatID2)

	tokenValue := redisMock.StorageToken()[chatID2]
	if tokenValue != command2 {
		t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, command2)
	}
}

func TestHandleHostMessageWithCorrectCommand(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	data := []struct {
		url    string
		chatID int64
	}{
		{"https://www.google.com", 44},
		{"https://www.tt.com", 45},
		{"https://tt.com", 46},
	}

	for _, message := range data {
		handler.Handle("host", message.url, message.chatID)
		hostValue := redisMock.StorageHost()[message.chatID]
		if hostValue != message.url {
			t.Errorf("Wrong saved host value %s is not %s", hostValue, message.url)
		}
	}
}

func TestHandleFillHoursNilTokenFailCommand(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	input := struct {
		command  string
		message  string
		chatID   int64
		expected string
	}{"fillhours", "43212 8 Test", 44, "Адрес сервера не задан! Пожалуйста задайте его с помощью команды /host <адрес сервера>"}

	handler.Handle(input.command, input.message, input.chatID)
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got: %s, expected: %s", botMock.text, input.expected)
	}
}

func TestHandleFillHoursNilHostFailCommand(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	input := struct {
		command  string
		message  string
		chatID   int64
		expected string
	}{"fillhours", "43212 8 Test", 44, "Адрес сервера не задан! Пожалуйста задайте его с помощью команды /host <адрес сервера>" +
		""}

	redisMock.SetToken("TestToken", input.chatID)
	handler.Handle(input.command, input.message, input.chatID)
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
	}
}
