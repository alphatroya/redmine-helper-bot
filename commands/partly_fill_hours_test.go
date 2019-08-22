package commands

import (
	"testing"
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
		storageMock := &StorageMock{}
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
