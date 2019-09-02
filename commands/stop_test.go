package commands

import (
	"testing"

	"github.com/alphatroya/redmine-helper-bot/mocks"
)

func TestStop_Handle(t *testing.T) {
	storageMock := mocks.NewStorageMock()
	var chatID int64 = 5
	storageMock.SetHost("https://google.com", chatID)
	storageMock.SetToken("dddsad", chatID)
	sut := newStopCommand(storageMock, chatID)
	result, err := sut.Handle("")
	expectedResult := "Бот остановлен, сохраненные данные очищены"
	if result != nil && result.Message() != expectedResult {
		t.Errorf("stop command should return correct message %s", result.Message())
	}
	if err != nil {
		t.Errorf("stop command should not return error, got: %s", err)
	}
	host, _ := storageMock.GetHost(chatID)
	token, _ := storageMock.GetToken(chatID)
	if host != "" || token != "" {
		t.Errorf("token and host values should be resetted")
	}
}

func TestStopCommand_IsCompleted(t *testing.T) {
	sut := StopCommand{}
	if sut.IsCompleted() != true {
		t.Errorf("stop command should always be completed")
	}
}

func TestStopConstructor(t *testing.T) {
	storageMock := mocks.NewStorageMock()
	if newStopCommand(storageMock, 5) == nil {
		t.Error("new stop command should return a new instance")
	}
}
