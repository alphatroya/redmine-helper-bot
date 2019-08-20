package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/mocks"
)

func TestSetHostCommand_Handle(t *testing.T) {
	data := []struct {
		text     string
		chatID   int64
		expected string
		error    error
	}{
		{"", 1, "", fmt.Errorf("Неправильное количество аргументов")},
		{" ", 1, "", fmt.Errorf("Неправильное количество аргументов")},
		{"test test", 1, "", fmt.Errorf("Неправильное количество аргументов")},
		{"test", 1, "", fmt.Errorf("parse test: invalid URI for request")},
		{"https://www.google.com", 1, "Адрес сервера успешно обновлен", nil},
	}

	for _, message := range data {
		storageMock := mocks.NewStorageMock()
		command := newSetHostCommand(storageMock, message.chatID)
		result, err := command.Handle(message.text)
		if result != nil && result.message != message.expected {
			t.Errorf("Wrong success response, input: %s, expected: %s, got: %s", message.text, message.expected, result)
		}
		if err != nil && err.Error() != message.error.Error() {
			t.Errorf("Wrong error response, expected: %s, got: %s", message.error, err.Error())
		}
	}
}

func TestSetHostCommand_IsCompleted(t *testing.T) {
	storageMock := mocks.NewStorageMock()
	command := newSetHostCommand(storageMock, 0)
	if command.IsCompleted() != true {
		t.Error("set host command should always be completed")
	}
}
