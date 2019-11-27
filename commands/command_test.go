package commands

import "testing"

func TestNewCommandResult(t *testing.T) {
	testMessage := "Test"
	sut := NewCommandResultWithMessages([]string{testMessage})
	if sut.Message() != testMessage {
		t.Error("constructor should set correct message property")
	}
}
