package commands

import (
	"reflect"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestPartlyFillHoursCommand_Handle(t *testing.T) {
	data := []struct {
		message     string
		isCompleted bool
		isHoursSet  bool
		result      *CommandResult
		err         error
	}{
		{message: "test", isCompleted: true, isHoursSet: true, result: NewCommandResult("_Операция выполнена_"), err: nil},
	}
	for _, item := range data {
		redmineMock := &RedmineMock{}
		storageMock := storage.NewStorageMock()
		printerMock := PrinterMock{}
		sut := newPartlyFillHoursCommand(redmineMock, printerMock, storageMock, 1)
		sut.isCompleted = item.isCompleted
		sut.isHoursSet = item.isHoursSet
		result, err := sut.Handle(item.message)
		if result != nil && result.Message() != item.result.Message() {
			t.Errorf("wrong result from handle method, got: %s, expected: %s", result.Message(), item.result.Message())
		}
		if err != nil && err != item.err {
			t.Errorf("wrong error from handle method, got: %s, expected: %s", err, item.err)
		}
	}
}

func TestNewFillHoursCommand(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := storage.NewStorageMock()
	printerMock := PrinterMock{}
	data := []struct {
		message string
		chatID  int64
		result  *PartlyFillHoursCommand
		isErr   bool
	}{
		{
			"54422 8 Test",
			5,
			&PartlyFillHoursCommand{
				redmineClient:   redmineMock,
				printer:         printerMock,
				storage:         storageMock,
				chatID:          5,
				issuesRequested: true,
				isIssueIDSet:    true,
				isHoursSet:      true,
				isCompleted:     false,
				issueID:         "54422",
				activityID:      "",
				hours:           "8",
				comment:         "Test",
				shortVersion:    true,
			},
			false,
		},
		{
			"54422 8 Test test",
			5,
			&PartlyFillHoursCommand{
				redmineClient:   redmineMock,
				printer:         printerMock,
				storage:         storageMock,
				chatID:          5,
				issuesRequested: true,
				isIssueIDSet:    true,
				isHoursSet:      true,
				isCompleted:     false,
				issueID:         "54422",
				activityID:      "",
				hours:           "8",
				comment:         "Test test",
				shortVersion:    true,
			},
			false,
		},
		{
			"#54422 8 Test test",
			5,
			&PartlyFillHoursCommand{
				redmineClient:   redmineMock,
				printer:         printerMock,
				storage:         storageMock,
				chatID:          5,
				issuesRequested: true,
				isIssueIDSet:    true,
				isHoursSet:      true,
				isCompleted:     false,
				issueID:         "54422",
				activityID:      "",
				hours:           "8",
				comment:         "Test test",
				shortVersion:    true,
			},
			false,
		},
		{
			"54422",
			5,
			nil,
			true,
		},
		{
			"54422 8",
			5,
			nil,
			true,
		},
		{
			"5442s2 8",
			5,
			nil,
			true,
		},
		{
			"54422 8a",
			5,
			nil,
			true,
		},
		{
			"54422 8    ",
			5,
			nil,
			true,
		},
	}
	for _, item := range data {
		result, err := NewFillHoursCommand(redmineMock, printerMock, storageMock, item.chatID, item.message)
		if item.isErr && err == nil {
			t.Errorf("not getting error when should")
		}
		if result != nil && !reflect.DeepEqual(result, item.result) {
			t.Errorf("getting wrong result, expected: %v, got: %v", item.result, result)
		}
	}
}

func TestPartlyFillHoursCommand_HelpMessage(t *testing.T) {
	helpMessage := newPartlyFillHoursCommand(nil, nil, nil, 0).HelpMessage()
	if len(helpMessage) == 0 {
		t.Errorf("help message should now be nil")
	}
}
