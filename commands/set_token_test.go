package commands

import (
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestSetTokenCommand_Handle(t *testing.T) {
	const wrongArgumentsCountText = "Неправильное количество аргументов, введите токен доступа к АПИ в формате `/token <токен>`"
	data := []struct {
		text     string
		chatID   int64
		expected string
		error    error
	}{
		{"", 1, wrongArgumentsCountText, nil},
		{"test test", 1, wrongArgumentsCountText, nil},
		{"fdsjfdsj", 1, "Токен успешно обновлен", nil},
		{"  ", 1, wrongArgumentsCountText, nil},
	}

	for _, message := range data {
		storageMock := storage.NewStorageMock()
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
	storageMock := storage.NewStorageMock()
	command := newSetTokenCommand(storageMock, 0)
	if command.IsCompleted() != true {
		t.Error("set token command should always be completed")
	}
}
