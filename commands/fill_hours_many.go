package commands

import (
	"errors"
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
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
		if regexp.MustCompile(`^#?[0-9]+$`).MatchString(fragment) {
			fragment = strings.TrimLeft(fragment, "#")
			issues = append(issues, fragment)
		} else {
			comment = strings.Join(fragments[i:], "")
			break
		}
	}
	if len(comment) == 0 {
		return nil, errors.New("Вы не ввели комментарий для задач")
	}

	timeEntries, err := f.redmineClient.TodayTimeEntries()
	if err != nil {
		return nil, err
	}

	var storedHours = 0
	for _, entry := range timeEntries {
		storedHours += int(math.Ceil(float64(entry.Hours)))
	}

	if storedHours >= 8 {
		return nil, errors.New("Вы сегодня уже работали 8 часов")
	}

	remainingHours := 8 - storedHours
	issuesCount := len(issues)
	if issuesCount > remainingHours {
		return nil, errors.New("Вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов")
	}

	var hours []string
	var remainingIssuesCount = issuesCount
	for range issues {
		hour := int(math.Ceil(float64(remainingHours) / float64(remainingIssuesCount)))
		hours = append(hours, fmt.Sprintf("%d", hour))
		remainingHours -= hour
		remainingIssuesCount--
	}

	host, err := f.storage.GetHost(f.chatID)
	if err != nil {
		return nil, err
	}

	responseMessage := fmt.Sprintf("[Задачи](%s/time_entries) успешно обновлены!\n", host)
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Задача", "Часы"})
	for i, issue := range issues {
		if i < len(hours) {
			fillHoursResponse, err := f.redmineClient.FillHoursRequest(issue, hours[i], comment, "")
			if err != nil {
				continue
			}
			table.Append([]string{fmt.Sprintf("%d", fillHoursResponse.TimeEntry.Issue.ID), fmt.Sprintf("%.0f", fillHoursResponse.TimeEntry.Hours)})
		}
	}
	table.Render()
	responseMessage += "`" + tableString.String() + "`"

	return NewCommandResult(responseMessage), nil
}

func (f FillHoursMany) IsCompleted() bool {
	return true
}
