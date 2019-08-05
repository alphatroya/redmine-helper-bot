package main

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func setupSubTest(t *testing.T) func(t *testing.T) {
	setup()
	return func(t *testing.T) {
		tearDown()
	}
}

var botMock *MockBotSender
var redisMock *RedisMock
var redmineMock *RedmineClientMock
var handler *UpdateHandler

func setup() {
	botMock = &MockBotSender{}
	redisMock = NewRedisMock()
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
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		handler.Handle(message.command, message.chatID)
		if botMock.text != message.expected {
			t.Errorf("Wrong response expected: %s received: %s", message.expected, botMock.text)
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

		handler.Handle("/token "+message.command, message.chatID)
		tokenValue := redisMock.Get(fmt.Sprint(message.chatID) + "_token").Val()
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

	handler.Handle("/token "+command, chatID)
	handler.Handle("/token "+command2, chatID2)

	tokenValue := redisMock.Get(fmt.Sprint(chatID2) + "_token").Val()
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
		handler.Handle("/host"+" "+message.url, message.chatID)
		hostValue := redisMock.Get(fmt.Sprint(message.chatID) + "_host").Val()
		if hostValue != message.url {
			t.Errorf("Wrong saved host value %s is not %s", hostValue, message.url)
		}
	}
}

func TestHandleFillHoursSuccessCommand(t *testing.T) {
	host := "https://test_host.com"
	tables := []struct {
		issue    string
		hours    string
		message  string
		chatID   int64
		expected string
	}{
		{"43212", "8", "Test", 44, SuccessFillHoursMessageResponse("43212", nil, "8", host)},
		{"51293", "8.0", "Test", 44, SuccessFillHoursMessageResponse("51293", nil, "8.0", host)},
		{"51293", "9.6", "Test", 44, SuccessFillHoursMessageResponse("51293", nil, "9.6", host)},
	}

	for _, message := range tables {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)
		redmineBody := &redmine.TimeEntryBody{
			&redmine.TimeEntry{
				message.issue,
				message.hours,
				message.message,
			},
		}
		redmineMock.SetFillHoursResponse(redmineBody, nil)
		redisMock.Set(fmt.Sprint(message.chatID)+"_token", "Test_TOKEN", 0)
		redisMock.Set(fmt.Sprint(message.chatID)+"_host", host, 0)
		handler.Handle(fmt.Sprintf("/fillhours %s %s %s", message.issue, message.hours, message.hours), message.chatID)
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
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursTokenNilResponse}

	handler.Handle(input.message, input.chatID)
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
	}
}

func TestHandleFillHoursNilHostFailCommand(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursHostNilResponse}

	redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
	handler.Handle(input.message, input.chatID)
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
		{"/fillhours aaaa 8 Test", 44, WrongFillHoursWrongIssueIDResponse},
		{"/fillhours <51293 8 Test", 44, WrongFillHoursWrongIssueIDResponse},
		{"/fillhours 51293 8a Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293 ff Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293 9,6 Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293", 44, WrongFillHoursWrongNumberOfArgumentsResponse},
		{"/fillhours 51293 6", 44, WrongFillHoursWrongNumberOfArgumentsResponse},
	}

	for _, input := range inputs {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		redisMock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		handler.Handle(input.message, input.chatID)
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}

func TestFillHoursWrongResponse(t *testing.T) {
	inputs := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 51293 8 Test", 44, fmt.Sprintf(WrongFillHoursWrongStatusCodeResponse, 400, "Bad Request")},
	}

	for _, input := range inputs {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		redisMock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		redmineMock.SetFillHoursResponse(&redmine.TimeEntryBody{nil}, redmine.WrongRedmineStatusCodeError(400, "Bad Request"))
		handler = &UpdateHandler{botMock, redisMock, redmineMock}
		handler.Handle(input.message, input.chatID)
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}
