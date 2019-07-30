package main

import (
	"fmt"
	"testing"
)

func TestHandleFillHoursSuccessCommand(t *testing.T) {
	host := "https://test_host.com"
	tables := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 43212 8 Test", 44, SuccessFillHoursMessageResponse("43212", "8", host)},
		{"/fillhours 51293 8.0 Test", 44, SuccessFillHoursMessageResponse("51293", "8.0", host)},
		{"/fillhours 51293 9.6 Test", 44, SuccessFillHoursMessageResponse("51293", "9.6", host)},
	}

	for _, message := range tables {
		botMock := &MockBotSender{}
		redisMock := NewRedisMock()
		redisMock.Set(fmt.Sprint(message.chatID)+"_token", "Test_TOKEN", 0)
		redisMock.Set(fmt.Sprint(message.chatID)+"_host", host, 0)
		HandleUpdate(botMock, message.message, message.chatID, redisMock, &ClientRequestMock{})
		if botMock.text != message.expected {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, message.expected)
		}
	}
}

func TestHandleFillHoursNilTokenFailCommand(t *testing.T) {
	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursTokenNilResponse}

	botMock := &MockBotSender{}
	redisMock := NewRedisMock()
	HandleUpdate(botMock, input.message, input.chatID, redisMock, &ClientRequestMock{})
	if input.expected != botMock.text {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
	}
}

func TestHandleFillHoursNilHostFailCommand(t *testing.T) {
	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursHostNilResponse}

	botMock := &MockBotSender{}
	redisMock := NewRedisMock()
	redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
	HandleUpdate(botMock, input.message, input.chatID, redisMock, &ClientRequestMock{})
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
		botMock := &MockBotSender{}
		redisMock := NewRedisMock()
		redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		redisMock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		HandleUpdate(botMock, input.message, input.chatID, redisMock, &ClientRequestMock{})
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}

func TestFillHoursWrongRsponse(t *testing.T) {
	inputs := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 51293 8 Test", 44, fmt.Sprintf(WrongFillHoursWrongStatusCodeResponse, 400, "Bad Request")},
	}

	for _, input := range inputs {
		botMock := &MockBotSender{}
		redisMock := NewRedisMock()
		redisMock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		redisMock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		HandleUpdate(botMock, input.message, input.chatID, redisMock, &ClientRequestMock{400})
		if input.expected != botMock.text {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", botMock.text, input.expected)
		}
	}
}
