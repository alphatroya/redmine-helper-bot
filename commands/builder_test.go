package commands

import (
	"reflect"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestBotCommandsBuilder_Build(t *testing.T) {
	checkData := []struct {
		objectType string
		message    string
		command    string
	}{
		{"*commands.SetHostCommand", "", "host"},
		{"*commands.SetTokenCommand", "", "token"},
		{"*commands.Activities", "", "activity"},
		{"*commands.IntroCommand", "", "start"},
		{"*commands.StopCommand", "", "stop"},
		{"*commands.UnknownCommand", "", ""},
		{"*commands.UnknownCommand", "", "sss"},
		{"*commands.PartlyFillHoursCommand", "54223 4 Test", "fh"},
		{"*commands.UnknownCommand", "", "fh"},
		{"*commands.FillStatus", "", "fstatus"},
		{"*commands.FillHoursMany", "", "fhm"},
	}
	for _, input := range checkData {
		mockStorage := storage.NewStorageMock()
		sut := NewBotCommandsBuilder(mockStorage)
		buildResult := sut.Build(input.command, input.message, 0)
		if reflect.TypeOf(buildResult).String() != input.objectType {
			t.Errorf("get wrong object from builder, input %s, got %T, expected %s", input.command, buildResult, input.objectType)
		}
	}
}
