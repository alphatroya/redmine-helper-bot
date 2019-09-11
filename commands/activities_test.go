package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func TestActivities_Handle(t *testing.T) {
	redmineMock := &RedmineMock{mockActivities: []*redmine.Activities{
		{Id: 5, Name: "Test"},
		{Id: 1, Name: "Test2"},
	}}
	storageMock := &StorageMock{}
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	result, err := sut.Handle("")
	expectedResult := "`+----+----------+\n| ID | НАЗВАНИЕ |\n+----+----------+\n|  5 | Test     |\n|  1 | Test2    |\n+----+----------+\n`"
	if result != nil && result.Message() != expectedResult {
		t.Errorf("getting wrong result text expected: %q, got: %q", expectedResult, result.Message())
	}

	if err != nil {
		t.Errorf("handling result should not be nil, got: %s", err)
	}
}

func TestActivities_Handle2(t *testing.T) {
	redmineMock := &RedmineMock{err: fmt.Errorf("test")}
	storageMock := &StorageMock{}
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	result, err := sut.Handle("")
	expectedResult := "test"
	if err != nil && err.Error() != expectedResult {
		t.Errorf("getting wrong result text expected: %s, got: %s", expectedResult, err)
	}

	if result != nil {
		t.Errorf("handling result should not be nil, got: %s", err)
	}
}

func TestActivities_IsCompleted(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := &StorageMock{}
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	if sut.IsCompleted() != true {
		t.Error("activities command should always be completed")
	}
}
