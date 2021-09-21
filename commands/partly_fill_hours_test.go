package commands

import (
	"reflect"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestNewFillHoursCommand(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := storage.NewStorageMock()
	printerMock := PrinterMock{}
	data := []struct {
		message string
		chatID  int64
		result  *FillHoursCommand
		isErr   bool
	}{
		{
			"54422 8 Test",
			5,
			&FillHoursCommand{
				redmineClient: redmineMock,
				printer:       printerMock,
				storage:       storageMock,
				chatID:        5,
				isIssueIDSet:  true,
				issueID:       "54422",
				activityID:    "",
				hours:         "8",
				comment:       "Test",
			},
			false,
		},
		{
			"54422 8 Test test",
			5,
			&FillHoursCommand{
				redmineClient: redmineMock,
				printer:       printerMock,
				storage:       storageMock,
				chatID:        5,
				isIssueIDSet:  true,
				issueID:       "54422",
				activityID:    "",
				hours:         "8",
				comment:       "Test test",
			},
			false,
		},
		{
			"#54422 8 Test test",
			5,
			&FillHoursCommand{
				redmineClient: redmineMock,
				printer:       printerMock,
				storage:       storageMock,
				chatID:        5,
				isIssueIDSet:  true,
				issueID:       "54422",
				activityID:    "",
				hours:         "8",
				comment:       "Test test",
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
