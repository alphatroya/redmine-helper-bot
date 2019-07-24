package main

import (
	"fmt"
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
		mock := NewRedisMock("_token")
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
		mock := NewRedisMock("_token")
		result := HandleTokenMessage("/token "+message.command, mock, message.chatID)
		tokenValue := mock.Get(fmt.Sprint(message.chatID)).Val()
		if tokenValue != message.command {
			t.Errorf("Wrong token storage logic: %s is not %s", tokenValue, message.command)
		}
		if result != SuccessTokenMessageResponse {
			t.Errorf("Wrong response from handling token command")
		}
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
		mock := NewRedisMock("_host")
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
		mock := NewRedisMock("_host")
		result, err := HandleHostMessage("/host"+" "+message.url, mock, message.chatID)
		if err != nil {
			t.Errorf("Correct input returns error result %s", err)
		}
		hostValue := mock.Get(fmt.Sprint(message.chatID)).Val()
		if hostValue != message.url {
			t.Errorf("Wrong saved host value %s is not %s", hostValue, message.url)
		}
		if result != SuccessHostMessageResponse {
			t.Errorf("Wrong success host update response")
		}
	}
}
