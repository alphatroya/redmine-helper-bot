package commands

import (
	"github.com/alphatroya/redmine-helper-bot/mocks"
	"testing"
)

func TestFillCommand(t *testing.T) {
	data := []struct {
		text     string
		chatID   int64
		expected string
	}{
		{"", 1, WrongTokenMessageResponse},
		{"test test", 1, WrongTokenMessageResponse},
		{"fdsjfdsj", 1, SuccessTokenMessageResponse},
		{"  ", 1, WrongTokenMessageResponse},
	}

	for _, message := range data {
		storageMock := mocks.NewStorageMock()
		command := NewSetTokenCommand(storageMock, message.chatID)
		result := command.Handle(message.text)
		if result != message.expected {
			t.Errorf("Wrong success response, expected: %s, got: %s", message.expected, result)
		}
	}
}
