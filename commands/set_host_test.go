package commands

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/mocks"
	"testing"
)

func TestSetHostCommand(t *testing.T) {
	data := []struct {
		text     string
		chatID   int64
		expected string
		error    error
	}{
		{"", 1, "", fmt.Errorf(WrongHostMessageResponse)},
		{" ", 1, "", fmt.Errorf(WrongHostMessageResponse)},
		{"test test", 1, "", fmt.Errorf(WrongHostMessageResponse)},
		{"test", 1, "", fmt.Errorf("parse test: invalid URI for request")},
		{"https://www.google.com", 1, SuccessHostMessageResponse, nil},
	}

	for _, message := range data {
		storageMock := mocks.NewStorageMock()
		command := NewSetHostCommand(storageMock, message.chatID)
		result, err := command.Handle(message.text)
		if result != message.expected {
			t.Errorf("Wrong success response, input: %s, expected: %s, got: %s", message.text, message.expected, result)
		}
		if err != nil && err.Error() != message.error.Error() {
			t.Errorf("Wrong error response, expected: %s, got: %s", message.error, err.Error())
		}
	}
}
