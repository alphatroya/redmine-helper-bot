package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func TestActivities_HandleFirstPhase(t *testing.T) {
	redmineMock := &RedmineMock{mockActivities: []*redmine.Activities{
		{Id: 5, Name: "Test"},
		{Id: 1, Name: "Test2"},
	}}
	storageMock := storage.NewStorageMock()
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	result, err := sut.Handle("")
	expectedResults := []string{
		"`+----+----------+\n| ID | НАЗВАНИЕ |\n+----+----------+\n|  5 | Test     |\n|  1 | Test2    |\n+----+----------+\n`",
		"_Введите номер новой активности по умолчанию_",
	}
	if err != nil {
		t.Errorf("handling result should not be nil, got: %s", err)
	}
	if len(result.Messages()) != len(expectedResults) {
		t.Errorf("result count (%d) not match expected count (%d)", len(result.Messages()), len(expectedResults))
	}
	for i, result := range result.Messages() {
		expected := expectedResults[i]
		if result != expected {
			t.Errorf("getting wrong result text expected: %s, got: %s", expected, result)
		}
	}
	if len(result.buttons) != len(redmineMock.mockActivities) {
		t.Errorf("result buttons count (%d) not match expected count (%d)", len(result.buttons), len(redmineMock.mockActivities))
	}
	for i, button := range result.buttons {
		mockActivity := redmineMock.mockActivities[i]
		if button != fmt.Sprintf("%d - %s", mockActivity.Id, mockActivity.Name) {
			t.Errorf("received wrong buttons title, got: %s", button)
		}
	}
}

func TestActivities_HandleSecondPhase(t *testing.T) {
	successMessage := "_Активность по умолчанию успешно сохранена. Она будет использоваться при каждой команде заполнения часов_"
	errorMessage := "Вы ввели неправильный номер активности"
	inputData := []struct {
		input         string
		output        string
		isErr         bool
		savedActivity string
	}{
		{"5 - Test", successMessage, false, "5"},
		{"5", successMessage, false, "5"},
		{"Test", errorMessage, true, ""},
		{"6 - Test", errorMessage, true, ""},
		{"6", errorMessage, true, ""},
	}

	for _, data := range inputData {
		redmineMock := &RedmineMock{mockActivities: []*redmine.Activities{
			{Id: 5, Name: "Test"},
			{Id: 1, Name: "Test2"},
		}}
		storageMock := storage.NewStorageMock()
		sut := newActivitiesCommand(redmineMock, storageMock, 1)
		_, _ = sut.Handle("")

		result, err := sut.Handle(data.input)
		expectedResults := []string{successMessage}
		if data.isErr && err == nil {
			t.Errorf("method should fail but it's not, error: %s", err)
		}
		if err != nil {
			if data.isErr && err.Error() != data.output {
				t.Errorf("result error is not correct, got: %s, expected: %s", err, data.output)
			}
			continue
		}
		if activity, _ := storageMock.GetActivity(1); activity != data.savedActivity {
			t.Errorf("new activity id is not set, current value: %s", activity)
		}
		if len(result.Messages()) != len(expectedResults) {
			t.Errorf("result count (%d) not match expected count (%d)", len(result.Messages()), len(expectedResults))
		}
		for i, result := range result.Messages() {
			expected := expectedResults[i]
			if result != expected {
				t.Errorf("getting wrong result text expected: %s, got: %s", expected, result)
			}
		}
	}
}

func TestActivities_HandleFirstPhaseError(t *testing.T) {
	redmineMock := &RedmineMock{err: fmt.Errorf("test")}
	storageMock := storage.NewStorageMock()
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
	redmineMock := &RedmineMock{mockActivities: []*redmine.Activities{
		{Id: 5, Name: "Test"},
		{Id: 1, Name: "Test2"},
	}}
	storageMock := storage.NewStorageMock()
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	if sut.IsCompleted() {
		t.Error("activities command should not be completed after initialization")
	}
	_, _ = sut.Handle("Ttt")
	if sut.IsCompleted() {
		t.Error("activities command should not be completed after first phase command")
	}
	_, _ = sut.Handle("6 - Test2")
	if sut.IsCompleted() {
		t.Error("activities command should not be completed after not correct second phase command")
	}
	_, _ = sut.Handle("5 - Test")
	if !sut.IsCompleted() {
		t.Error("activities command should be completed after second phase command")
	}
}
