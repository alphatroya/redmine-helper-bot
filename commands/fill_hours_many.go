package commands

import (
	"errors"
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"math"
	"regexp"
	"strings"
)

type FillHoursMany struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
}

func NewFillHoursMany(redmineClient redmine.Client, storage storage.Manager, chatID int64) *FillHoursMany {
	return &FillHoursMany{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (f FillHoursMany) Handle(message string) (*CommandResult, error) {
	if len(message) == 0 {
		return nil, errors.New("Введена неправильная команда")
	}
	fragments := strings.Split(message, " ")
	var issues []string
	comment := ""
	for i, fragment := range fragments {
		if regexp.MustCompile(`^#?[0-9]+ - .+$`).MatchString(fragment) {
			fragment = strings.TrimLeft(fragment, "#")
			issues = append(issues, fragment)
		} else {
			comment = strings.Join(fragments[i:], "")
			break
		}
	}
	if len(comment) == 0 {
		return nil, errors.New("Ошибка: вы не ввели комментарий для задач")
	}

	timeEntries, err := f.redmineClient.TodayTimeEntries()
	if err != nil {
		return nil, err
	}

	var storedHours float32 = 0
	for _, entry := range timeEntries {
		storedHours += entry.Hours
	}

	remainingHours := 8 - storedHours
	if remainingHours > float32(len(issues)) {
		return nil, errors.New("Ошибка: вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов")
	}

	var hours []string
	for range issues {
		hour := math.Ceil(float64(remainingHours / float32(len(issues))))
		hours = append(hours, fmt.Sprintf("%d", int(hour)))
		remainingHours -= hour
	}
	//for _, issue := range issues {
	//	timeEntryResponse, err := f.redmineClient.FillHoursRequest(issue, "1", comment, "")
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return NewCommandResult("Test"), nil
}

func (f FillHoursMany) IsCompleted() bool {
	panic("implement me")
}
