package commands

import "testing"

func TestUnknownCommand_Handle(t *testing.T) {
	sut := NewUnknownCommand()
	result, err := sut.Handle("test")
	if result != nil && result.Message() != "Введена неправильная команда" {
		t.Errorf("wrong result message, got: %s", result)
	}
	if err != nil {
		t.Errorf("unknown command should always return nil error, got: %s", err)
	}
}

func TestUnknownCommand_Handle2(t *testing.T) {
	sut := NewUnknownCommandWithMessage("Foo")
	result, err := sut.Handle("test")
	if result != nil && result.Message() != "Foo" {
		t.Errorf("wrong result message, got: %s", result)
	}
	if err != nil {
		t.Errorf("unknown command should always return nil error, got: %s", err)
	}
}

func TestUnknownCommand_IsCompleted(t *testing.T) {
	sut := NewUnknownCommand()
	if !sut.IsCompleted() {
		t.Error("unknown command should be always completed")
	}
}
