package commands

import "testing"

func TestStop_Handle(t *testing.T) {
	sut := StopCommand{}
	result, err := sut.Handle("")
	expectedResult := "Бот остановлен"
	if result != nil && result.Message() != expectedResult {
		t.Errorf("stop command should return correct message %s", result.Message())
	}
	if err != nil {
		t.Errorf("stop command should not return error, got: %s", err)
	}
}

func TestStopCommand_IsCompleted(t *testing.T) {
	sut := StopCommand{}
	if sut.IsCompleted() != true {
		t.Errorf("stop command should always be completed")
	}
}

func TestStopConstructor(t *testing.T) {
	if newStopCommand() == nil {
		t.Error("new stop command should return a new instance")
	}
}
