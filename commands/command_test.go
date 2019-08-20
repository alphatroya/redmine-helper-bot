package commands

import "testing"

func TestNewCommandResult(t *testing.T) {
	testMessage := "Test"
	sut := NewCommandResult(testMessage)
	if sut.Message() != testMessage || sut.message != testMessage {
		t.Error("constructor should set correct message property")
	}
}
