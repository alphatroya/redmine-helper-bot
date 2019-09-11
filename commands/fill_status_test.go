package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func TestFillStatus_Handle(t *testing.T) {
	table := []struct {
		mockEntries []*redmine.TimeEntryResponse
		expected    string
	}{
		{
			[]*redmine.TimeEntryResponse{},
			"Вы сегодня работали *0.0 ч.*",
		},
		{
			[]*redmine.TimeEntryResponse{
				{
					Activity: redmine.TimeEntryResponseActivity{
						ID:   0,
						Name: "Разработка",
					},
					ID:       0,
					Hours:    5,
					Issue:    redmine.TimeEntryResponseIssue{ID: 55422},
					Comments: "Test 1",
				},
				{
					Activity: redmine.TimeEntryResponseActivity{
						ID:   1,
						Name: "Дизайн",
					},
					ID:       1,
					Hours:    4,
					Issue:    redmine.TimeEntryResponseIssue{ID: 55422},
					Comments: "Test 2",
				},
			},
			"Вы сегодня работали *9.0 ч.*\n\n`Часы | Задача | Активность | Комментарий`\n`-----+--------+------------+---------------------`\n` 5.0 | 55422  | Разработка | Test 1\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\n`` 4.0 | 55422  | Дизайн\x00\x00\x00\x00 | Test 2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\n`",
		},
	}
	for _, item := range table {
		storageMock := &StorageMock{}
		sut := NewFillStatus(RedmineMock{mockTimeEntries: item.mockEntries}, storageMock, 5)
		result, err := sut.Handle("")
		if err != nil {
			t.Errorf("success conditions should complete without error, got: %s", err)
		}

		expectedMessage := item.expected
		if result.Message() != expectedMessage {
			t.Errorf("success conditions should result correct output message, got %q, expected %q", result.Message(), expectedMessage)
		}
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
