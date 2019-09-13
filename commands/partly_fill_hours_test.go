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
		{message: "test", isCompleted: true, isHoursSet: true, result: NewCommandResult("Операция выполнена"), err: nil},
	}
	for _, item := range data {
		redmineMock := &RedmineMock{}
		storageMock := storage.NewStorageMock()
		sut := newPartlyFillHoursCommand(redmineMock, storageMock, 1)
		sut.isCompleted = item.isCompleted
		sut.isHoursSet = item.isHoursSet
		result, err := sut.Handle(item.message)
		if result != nil && result.message != item.result.message {
			t.Errorf("wrong result from handle method, got: %s, expected: %s", result.message, item.result.message)
		}
		if err != nil && err != item.err {
			t.Errorf("wrong error from handle method, got: %s, expected: %s", err, item.err)
		}
	}
}

func TestNewFillHoursCommand(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := storage.NewStorageMock()
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
			},
			false,
		},
		{
			"54422 8 Test test",
			5,
			&PartlyFillHoursCommand{
				redmineClient:   redmineMock,
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
			},
			false,
		},
		{
			"#54422 8 Test test",
			5,
			&PartlyFillHoursCommand{
				redmineClient:   redmineMock,
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
		result, err := NewFillHoursCommand(redmineMock, storageMock, item.chatID, item.message)
		if item.isErr && err == nil {
			t.Errorf("not getting error when should")
		}
		if result != nil && !reflect.DeepEqual(result, item.result) {
			t.Errorf("getting wrong result, expected: %v, got: %v", item.result, result)
		}
	}
}
