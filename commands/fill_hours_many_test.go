package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestFillHoursMany_Handle(t *testing.T) {
	data := []struct {
		input              string
		output             string
		storedHours        float32
		resultErr          string
		fillHoursErrorsMap map[string]bool
	}{
		{
			input: "52233 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "4"},
				}),
			storedHours: 4,
			resultErr:   "",
		},
		{
			input: "52233 54223 53312 52551 #Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "1"},
					{"52551", "1"},
					{"53312", "1"},
					{"54223", "1"},
				}),
			storedHours: 4,
			resultErr:   "",
		},
		{
			input: "#52233 54223 #53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "1"},
					{"52551", "1"},
					{"53312", "1"},
					{"54223", "1"},
				}),
			storedHours: 4,
			resultErr:   "",
		},
		{
			input: "52233 54223 53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "1"},
					{"52551", "1"},
					{"53312", "1"},
					{"54223", "1"},
				}),
			storedHours: 4,
			resultErr:   "",
		},
		{
			input: "52233 52233 53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "3"},
					{"52551", "3"},
					{"53312", "2"},
				}),
			storedHours: 0,
			resultErr:   "",
		},
		{
			input: "52233 54223 53312 52551 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "2"},
					{"52551", "2"},
					{"53312", "2"},
					{"54223", "2"},
				}),
			storedHours: 0,
			resultErr:   "",
		},
		{
			input: "52233 54223 53312 Исправление",
			output: successMessage(
				"https://google.com",
				true,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "3"},
					{"53312", "3"},
					{"54223", "2"},
				}),
			storedHours: 0,
			resultErr:   "",
		},
		{
			input: "52233 63312 74223 Исправление",
			output: successMessage(
				"https://google.com",
				false,
				[]struct {
					id    string
					hours string
				}{
					{"52233", "3"},
				}) + "\n" + errorsMessage(5,
				[]struct {
					id    string
					hours string
				}{
					{"63312", "3"},
					{"74223", "2"},
				}),
			storedHours: 0,
			resultErr:   "",
			fillHoursErrorsMap: map[string]bool{
				"63312": true,
				"74223": true,
			},
		},
		{
			input:       "52233 54223 53312 Исправление",
			output:      "",
			storedHours: 8,
			resultErr:   "Вы сегодня уже работали 8 часов",
		},
		{
			input:       "54221 44221 44421 Test",
			output:      "",
			storedHours: 6,
			resultErr:   "Вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов",
		},
		{
			input:       "54221 44221 44421 73312 3332 21211 333211 44321 32321 Test",
			output:      "",
			storedHours: 0,
			resultErr:   "Вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов",
		},
		{
			input:     "",
			output:    "",
			resultErr: NewFillHoursMany(nil, nil, 0).HelpMessage(),
		},
		{
			input:     "53223 43231 33321",
			output:    "",
			resultErr: "Вы не ввели комментарий для задач",
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
		redmineMock.fillHoursErrorsMap = item.fillHoursErrorsMap

		command := NewFillHoursMany(redmineMock, storageMock, chatID)
		result, err := command.Handle(item.input)

		if err != nil {
			if err.Error() != item.resultErr {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, item.resultErr)
			} else if len(redmineMock.filledIssues) != 0 {
				t.Errorf("error command should not edit issues")
			}
			continue
		}

		if len(result.buttons) != 0 {
			t.Error("success command should not return buttons")
		}

		if result.Message() != item.output {
			t.Errorf("command: %s\nreturn wrong message\ngot: %s\nexpected: %s", item.input, result.Message(), item.output)
		}
	}
}

func successMessage(host string, success bool, tasks []struct {
	id    string
	hours string
}) string {
	var result string
	if success {
		result = fmt.Sprintf("Задачи([%d](%s/time_entries)) успешно обновлены!\n\n", len(tasks), host)
	} else {
		result = fmt.Sprintf("Задачи([%d](%s/time_entries)) *частично* обновлены, обновленные задачи\n\n", len(tasks), host)
	}
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

func errorsMessage(remain int, tasks []struct {
	id    string
	hours string
}) string {
	result := "Не удалось обновить задачи\n\n"
	result += "`+--------+\n"
	result += "| ЗАДАЧА |\n"
	result += "+--------+\n"
	for _, task := range tasks {
		result += fmt.Sprintf("|  %s |\n", task.id)
	}
	result += "+--------+\n"
	result += "`\n"
	result += fmt.Sprintf("Не удалось распределить %d ч.", remain)
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

func TestFillHoursMany_HelpMessage(t *testing.T) {
	storageMock := storage.NewStorageMock()
	redmineMock := &RedmineMock{}
	command := NewFillHoursMany(redmineMock, storageMock, 5)
	if len(command.HelpMessage()) == 0 {
		t.Errorf("Help message should not be empty")
	}
}
