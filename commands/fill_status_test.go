package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func TestFillStatus_Handle(t *testing.T) {
	redmineMock := &RedmineMock{mockTimeEntries: []*redmine.TimeEntryResponse{
		{
			ID:    0,
			Hours: 5,
		},
		{
			ID:    1,
			Hours: 4,
		},
	},
	}
	storageMock := &StorageMock{}
	sut := NewFillStatus(redmineMock, storageMock, 5)
	result, err := sut.Handle("")
	if err != nil {
		t.Errorf("success conditions should complete without error, got: %s", err)
	}

	expectedMessage := "Вы сегодня работали 9.0 ч."
	if result.Message() != expectedMessage {
		t.Errorf("success conditions should result correct output message, got %s, expected %s", result.Message(), expectedMessage)
	}
}

func TestFillStatus_Handle2(t *testing.T) {
	redmineMock := &RedmineMock{err: fmt.Errorf("test")}
	storageMock := &StorageMock{}
	sut := NewFillStatus(redmineMock, storageMock, 1)
	result, err := sut.Handle("")
	expectedResult := "test"
	if err != nil && err.Error() != expectedResult {
		t.Errorf("getting wrong result text expected: %s, got: %s", expectedResult, err)
	}

	if result != nil {
		t.Errorf("handling result should not be nil, got: %s", err)
	}
}

func TestFillStatus_IsCompleted(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := &StorageMock{}
	sut := NewFillStatus(redmineMock, storageMock, 5)
	if sut.IsCompleted() != true {
		t.Error("fill status command should always be completed")
	}
}
