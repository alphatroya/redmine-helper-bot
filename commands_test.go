package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandleTokenMessageWithWrongCommand(t *testing.T) {
	tables := []struct {
		command string
		failure string
	}{
		{"/token", "empty"},
		{"/token test test", "too many arguments"},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		if HandleTokenMessage(message.command, mock, 0) != WrongTokenMessageResponse {
			t.Errorf("Handling token should fail if it is %s", message.failure)
		}
	}
}

func TestHandleTokenMessageWithCorrectCommand(t *testing.T) {
	tables := []struct {
		command string
		chatID  int64
	}{
		{"431", 44},
		{"23", 45},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		result := HandleTokenMessage("/token "+message.command, mock, message.chatID)
		tokenValue := mock.Get(fmt.Sprint(message.chatID) + "_token").Val()
		if tokenValue != message.command {
			t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, message.command)
		}
		if result != SuccessTokenMessageResponse {
			t.Errorf("Wrong response from handling token command")
		}
	}
}

func TestHandleTokenMessageWithMultipleCommands(t *testing.T) {
	command := "4433"
	var chatID int64 = 1

	command2 := "4433"
	var chatID2 int64 = 2

	mock := NewRedisMock()

	_ = HandleTokenMessage("/token "+command, mock, chatID)
	result := HandleTokenMessage("/token "+command2, mock, chatID2)

	tokenValue := mock.Get(fmt.Sprint(chatID2) + "_token").Val()
	if tokenValue != command2 {
		t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, command2)
	}
	if result != SuccessTokenMessageResponse {
		t.Errorf("Wrong response from handling token command")
	}
}

func TestHandleHostMessageWithWrongCommand(t *testing.T) {
	tables := []struct {
		message string
		failure string
	}{
		{"/host", "empty command"},
		{"/host test test", "too many arguments"},
		{"/host test", "not correct URL"},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		_, err := HandleHostMessage(message.message, mock, 0)
		if err == nil {
			t.Errorf("Handling host should fail if %s", message.failure)
		}
	}
}

func TestHandleHostMessageWithCorrectCommand(t *testing.T) {
	tables := []struct {
		url    string
		chatID int64
	}{
		{"https://www.google.com", 44},
		{"https://www.tt.com", 45},
		{"https://tt.com", 46},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		result, err := HandleHostMessage("/host"+" "+message.url, mock, message.chatID)
		if err != nil {
			t.Errorf("Correct input returns error result %s", err)
		}
		hostValue := mock.Get(fmt.Sprint(message.chatID) + "_host").Val()
		if hostValue != message.url {
			t.Errorf("Wrong saved host value %s is not %s", hostValue, message.url)
		}
		if result != SuccessHostMessageResponse {
			t.Errorf("Wrong success host update response")
		}
	}
}

func TestHandleFillHoursSuccessCommand(t *testing.T) {
	tables := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 43212 8 Test", 44, fmt.Sprintf(SuccessFillHoursMessageResponse, "43212", "8")},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		mock.Set(fmt.Sprint(message.chatID)+"_token", "TestToken", 0)
		mock.Set(fmt.Sprint(message.chatID)+"_host", "https://test_host.com", 0)
		text, err := HandleFillMessage(message.message, message.chatID, mock, &ClientRequestMock{})
		if err != nil {
			t.Errorf("Success function should not return err %s", err)
		}
		if text != message.expected {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", text, message.expected)
		}
	}
}

func TestHandleFillHoursNilTokenFailCommand(t *testing.T) {
	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursTokenNilResponse}

	mock := NewRedisMock()
	_, err := HandleFillMessage(input.message, input.chatID, mock, &ClientRequestMock{})
	if err == nil {
		t.Errorf("Wrong command should return non-nil err")
	}
	if input.expected != err.Error() {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", err.Error(), input.expected)
	}
}

type ClientRequestMock struct {
}

func (c *ClientRequestMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}
