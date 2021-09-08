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
			"_Вы сегодня работали_ *0.0 ч.*",
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
			"_Вы сегодня работали_ *9.0 ч.*\n\n`+------+--------+------------+-------------+\n| ЧАСЫ | ЗАДАЧА | АКТИВНОСТЬ | КОММЕНТАРИЙ |\n+------+--------+------------+-------------+\n|  5.0 |  55422 | Разработка | Test 1      |\n|  4.0 |  55422 | Дизайн     | Test 2      |\n+------+--------+------------+-------------+\n`",
		},
	}
	for _, item := range table {
		sut := NewFillStatus(&RedmineMock{mockTimeEntries: item.mockEntries}, 5)
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
	sut := NewFillStatus(redmineMock, 1)
	result, err := sut.Handle("")
	if expectedResult := "test"; err != nil && err.Error() != expectedResult {
		t.Errorf("getting wrong result text expected: %s, got: %s", expectedResult, err)
	}

	if result != nil {
		t.Errorf("handling result should not be nil, got: %s", err)
	}
}

func TestFillStatus_IsCompleted(t *testing.T) {
	redmineMock := &RedmineMock{}
	sut := NewFillStatus(redmineMock, 5)
	if sut.IsCompleted() != true {
		t.Error("fill status command should always be completed")
	}
}
