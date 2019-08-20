package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/mocks"
)

func TestSetTokenCommand_Handle(t *testing.T) {
	data := []struct {
		text     string
		chatID   int64
		expected string
		error    error
	}{
		{"", 1, "", fmt.Errorf(WrongTokenMessageResponse)},
		{"test test", 1, "", fmt.Errorf(WrongTokenMessageResponse)},
		{"fdsjfdsj", 1, SuccessTokenMessageResponse, nil},
		{"  ", 1, "", fmt.Errorf(WrongTokenMessageResponse)},
	}

	for _, message := range data {
		storageMock := mocks.NewStorageMock()
		command := newSetTokenCommand(storageMock, message.chatID)
		result, err := command.Handle(message.text)
		if result != nil && result.Message() != message.expected {
			t.Errorf("Wrong success response, expected: %s, got: %s", message.expected, result)
		}
		if err != nil && err.Error() != message.error.Error() {
			t.Errorf("Wrong error response, expected: %s, got: %s", message.error, err.Error())
		}
	}
}

func TestSetTokenCommand_IsCompleted(t *testing.T) {
	storageMock := mocks.NewStorageMock()
	command := newSetTokenCommand(storageMock, 0)
	if command.IsCompleted() != true {
		t.Error("set token command should always be completed")
	}
}
