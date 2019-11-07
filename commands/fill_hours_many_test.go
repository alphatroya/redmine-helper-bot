package commands

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"testing"
)

func TestFillHoursMany_Handle(t *testing.T) {
	data := []struct {
		input       string
		output      string
		storedHours float32
		err         string
	}{
		{
			input: "52233 Исправление",
			output: successMessage(
				"https://google.com",
				[]struct {
					id    string
					hours string
				}{
					{"52233", "4"},
				}),
			storedHours: 4,
			err:         "",
		},
		{
			input: "52233 54223 53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				[]struct {
					id    string
					hours string
				}{
					{"52233", "1"},
					{"54223", "1"},
					{"53312", "1"},
					{"52551", "1"},
				}),
			storedHours: 4,
			err:         "",
		},
		{
			input: "52233 54223 53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				[]struct {
					id    string
					hours string
				}{
					{"52233", "2"},
					{"54223", "2"},
					{"53312", "2"},
					{"52551", "2"},
				}),
			storedHours: 0,
			err:         "",
		},
		{
			input: "52233 54223 53312 Исправление",
			output: successMessage(
				"https://google.com",
				[]struct {
					id    string
					hours string
				}{
					{"52233", "3"},
					{"54223", "3"},
					{"53312", "2"},
				}),
			storedHours: 0,
			err:         "",
		},
		{
			input:       "52233 54223 53312 Исправление",
			output:      "",
			storedHours: 8,
			err:         "Вы сегодня уже работали 8 часов",
		},
		{
			input:       "54221 44221 44421 Test",
			output:      "",
			storedHours: 6,
			err:         "Вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов",
		},
		{
			input:  "",
			output: "",
			err:    "Введена неправильная команда",
		},
		{
			input:  "53223 43231 33321",
			output: "",
			err:    "Вы не ввели комментарий для задач",
		},
	}

	for _, item := range data {
		storageMock := storage.NewStorageMock()
		var chatID int64 = 5
		storageMock.SetHost("https://google.com", chatID)

		redmineMock := &RedmineMock{}
		redmineMock.mockTimeEntries = []*redmine.TimeEntryResponse{
			{Hours: item.storedHours},
		}

		command := NewFillHoursMany(redmineMock, storageMock, chatID)
		result, err := command.Handle(item.input)

		if err != nil {
			if err.Error() != item.err {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, item.err)
			}
			if len(redmineMock.filledIssues) != 0 {
				t.Errorf("error command should not edit issues")
			}
			return
		}

		if len(result.buttons) != 0 {
			t.Error("success command should not return buttons")
		}

		if result.Message() != item.output {
			t.Errorf("command return wrong message\ngot: %s\nexpected: %s", result.Message(), item.output)
		}
	}
}

func successMessage(host string, tasks []struct {
	id    string
	hours string
}) string {
	result := fmt.Sprintf("[Задачи](%s/time_entries) успешно обновлены!\n", host)
	result += "`+--------+------+\n"
	result += "| ЗАДАЧА | ЧАСЫ |\n"
	result += "+--------+------+\n"
	for _, task := range tasks {
		result += fmt.Sprintf("|  %s |    %s |\n", task.id, task.hours)
	}
	result += "+--------+------+\n"
	result += "`"
	return result
}

func TestFillHoursMany_IsCompleted(t *testing.T) {
	storageMock := storage.NewStorageMock()
	redmineMock := &RedmineMock{}
	command := NewFillHoursMany(redmineMock, storageMock, 5)
	if !command.IsCompleted() {
		t.Errorf("Fill command should always be completed")
	}
}
