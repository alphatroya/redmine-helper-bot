package commands

import (
	"reflect"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/mocks"
)

func TestBotCommandsBuilder_Build(t *testing.T) {
	checkData := []struct {
		objectType string
		command    string
	}{
		{"*commands.SetHostCommand", "host"},
		{"*commands.SetTokenCommand", "token"},
		{"*commands.PartlyFillHoursCommand", "fillhours"},
		{"*commands.Activities", "activities"},
		{"*commands.IntroCommand", "start"},
		{"*commands.StopCommand", "stop"},
		{"*commands.UnknownCommand", ""},
		{"*commands.UnknownCommand", "sss"},
	}
	for _, input := range checkData {
		mockStorage := mocks.NewStorageMock()
		sut := NewBotCommandsBuilder(mockStorage)
		buildResult := sut.Build(input.command, "", 0)
		if reflect.TypeOf(buildResult).String() != input.objectType {
			t.Errorf("get wrong object from builder, got %T, expected %s", buildResult, input.objectType)
		}
	}
}
