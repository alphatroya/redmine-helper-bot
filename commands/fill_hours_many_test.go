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
			input: "52233 54223 53312 Исправление",
			output: successMessage(
				"https://google.com",
				[]struct {
					id    string
					name  string
					hours string
				}{
					{"52233", "Test0", "3"},
					{"54223", "Test1", "3"},
					{"53312", "Test2", "2"},
				}),
			storedHours: 0,
			err:         "",
		},
		{
			input:       "54221 44221 44421 Test",
			output:      "",
			storedHours: 6,
			err:         "Ошибка: вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов",
		},
		{
			input:  "",
			output: "",
			err:    "Введена неправильная команда",
		},
		{
			input:  "53223 43231 33321",
			output: "",
			err:    "Ошибка: вы не ввели комментарий для задач",
		},
	}

	for _, item := range data {
		storageMock := storage.NewStorageMock()
		var chatID int64 = 5
		storageMock.SetHost("https://google.com", chatID)

		redmineMock := &RedmineMock{}
		redmineMock.mockTimeEntries = []*redmine.TimeEntryResponse{
			{Hours: 8 - item.storedHours},
		}

		command := NewFillHoursMany(redmineMock, storageMock, chatID)
		result, err := command.Handle(item.input)

		if err != nil {
			if err.Error() != item.err {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, item.err)
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
	name  string
	hours string
}) string {
	result := fmt.Sprintf("[Задачи](%s/time_entries) успешно обновлены!\n", host)
	for _, task := range tasks {
		result += fmt.Sprintf("%s | %s | %s\n", task.id, task.name, task.hours)
	}
	return result
}

func TestFillHoursMany_IsCompleted(t *testing.T) {
}
