package main

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/commands"
	"github.com/alphatroya/redmine-helper-bot/mocks"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func setupSubTest(t *testing.T) func(t *testing.T) {
	setup()
	return func(t *testing.T) {
		tearDown()
	}
}

var botMock *MockBotSender
var redisMock *mocks.StorageMock
var redmineMock *RedmineClientMock
var handler *UpdateHandler

func setup() {
	botMock = &MockBotSender{}
	redisMock = mocks.NewStorageMock()
	redmineMock = &RedmineClientMock{"", "", nil, nil}
	handler = &UpdateHandler{botMock, redisMock, redmineMock}
}

func tearDown() {
	botMock = nil
	redisMock = nil
	handler = nil
	redmineMock = nil
}

func TestTokenRequest(t *testing.T) {
	data := []struct {
		command  string
		message  string
		chatID   int64
		expected string
	}{
		{"token", "", 1, commands.WrongTokenMessageResponse},
		{"token", "test test", 1, commands.WrongTokenMessageResponse},
		{"token", "fdsjfdsj", 1, commands.SuccessTokenMessageResponse},
		{"token", "  ", 1, commands.WrongTokenMessageResponse},
		{"", "qwertyu", 1, "Введена неправильная команда"},
		{"host", "", 1, "Неправильное количество аргументов"},
		{"host", " ", 1, "Неправильное количество аргументов"},
		{"host", "test test", 1, "Неправильное количество аргументов"},
		{"host", "test", 1, "parse test: invalid URI for request"},
		{"host", "https://www.google.com", 1, "Адрес сервера успешно обновлен"},
	}

	for _, message := range data {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		handler.Handle(message.command, message.message, message.chatID)
		if botMock.text != message.expected {
			t.Errorf("Wrong response command: %s, arguments: %s, expected: %s, received: %s", message.command, message.message, message.expected, botMock.text)
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

func newTimeEntryResponse(comments string, hours float32, ID int) *redmine.TimeEntryResponse {
	return &redmine.TimeEntryResponse{Comments: comments, Hours: hours, Issue: struct {
		ID int "json:\"id\""
	}{ID: ID}}
}

func TestHandleFillHoursSuccessCommand(t *testing.T) {
	host := "https://test_host.com"
	tables := []struct {
		issueID  int
		hours    float32
		comments string
		chatID   int64
		expected string
	}{
		{43212, 8, "Test", 44, commands.SuccessFillHoursMessageResponse(43212, nil, 8, host)},
		{51293, 8.0, "Test", 44, commands.SuccessFillHoursMessageResponse(51293, nil, 8.0, host)},
		{51293, 9.6, "Test", 44, commands.SuccessFillHoursMessageResponse(51293, nil, 9.6, host)},
	}

	for _, message := range tables {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)
		redmineBody := newTimeEntryResponse(message.comments, message.hours, message.issueID)
		redmineMock.SetFillHoursResponse(&redmine.TimeEntryBodyResponse{TimeEntry: *redmineBody}, nil)
		redisMock.SetToken("Test_TOKEN", message.chatID)
		redisMock.SetHost(host, message.chatID)
		handler.Handle("fillhours", fmt.Sprintf("%d %f %s", message.issueID, message.hours, message.comments), message.chatID)
		if botMock.text != message.expected {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, message.expected)
		}
		if redmineMock.host != host {
			t.Errorf("Command should set host parameter, received %s", redmineMock.host)
		}
		if redmineMock.token != "Test_TOKEN" {
			t.Errorf("Command should set token parameter, received %s", redmineMock.token)
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
	}{"fillhours", "43212 8 Test", 44, commands.WrongFillHoursTokenNilResponse}

	handler.Handle(input.command, input.message, input.chatID)
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
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
	}{"fillhours", "43212 8 Test", 44, commands.WrongFillHoursHostNilResponse}

	redisMock.SetToken("TestToken", input.chatID)
	handler.Handle(input.command, input.message, input.chatID)
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
	}
}

func TestFillHoursWrongInput(t *testing.T) {
	inputs := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"aaaa 8 Test", 44, commands.WrongFillHoursWrongIssueIDResponse},
		{"<51293 8 Test", 44, commands.WrongFillHoursWrongIssueIDResponse},
		{"51293 8a Test", 44, commands.WrongFillHoursWrongHoursCountResponse},
		{"51293 ff Test", 44, commands.WrongFillHoursWrongHoursCountResponse},
		{"51293 9,6 Test", 44, commands.WrongFillHoursWrongHoursCountResponse},
		{"51293", 44, commands.WrongFillHoursWrongNumberOfArgumentsResponse},
		{"51293 6", 44, commands.WrongFillHoursWrongNumberOfArgumentsResponse},
	}

	for _, input := range inputs {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		redisMock.SetToken("Test_Token", input.chatID)
		redisMock.SetHost("https://test_host.com", input.chatID)
		handler.Handle("fillhours", input.message, input.chatID)
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}

func TestFillHoursWrongResponse(t *testing.T) {
	inputs := []struct {
		command  string
		message  string
		chatID   int64
		expected string
	}{
		{"fillhours", "51293 8 Test", 44, redmine.WrongStatusCodeError(400, "Bad Request").Error()},
	}

	for _, input := range inputs {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		redisMock.SetToken("Test_Token", input.chatID)
		redisMock.SetHost("https://test_host.com", input.chatID)
		redmineMock.SetFillHoursResponse(&redmine.TimeEntryBodyResponse{}, redmine.WrongStatusCodeError(400, "Bad Request"))
		handler = &UpdateHandler{botMock, redisMock, redmineMock}
		handler.Handle(input.command, input.message, input.chatID)
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}
