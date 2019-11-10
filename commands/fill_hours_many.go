package commands

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
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
	issues, comment, err := f.getIssuesAndComment(message)
	if err != nil {
		return nil, err
	}

	timeEntries, err := f.redmineClient.TodayTimeEntries()
	if err != nil {
		return nil, err
	}

	remainingHours, err := f.getRemainingHours(timeEntries)
	if err != nil {
		return nil, err
	}

	issuesCount := len(issues)
	if issuesCount > remainingHours {
		return nil, errors.New("Вы ввели слишком много номеров задач. В целях точного распределения задач за день количество ограничено числом свободных за день часов")
	}

	hours := f.getHours(issuesCount, issues, remainingHours)

	host, err := f.storage.GetHost(f.chatID)
	if err != nil {
		return nil, err
	}

	fillSuccess, fillErrors, remain := f.fillIssuesResult(issues, hours, comment)

	successTableString := &strings.Builder{}
	successTable := tablewriter.NewWriter(successTableString)
	successTable.SetHeader([]string{"Задача", "Часы"})

	for _, timeEntry := range fillSuccess {
		successTable.Append([]string{fmt.Sprintf("%d", timeEntry.TimeEntry.Issue.ID), fmt.Sprintf("%.0f", timeEntry.TimeEntry.Hours)})
	}
	successTable.Render()

	var responseMessage string
	if len(fillErrors) != 0 {
		responseMessage = "Задачи *частично* обновлены, обновленные задачи\n\n"
		responseMessage += "`" + successTableString.String() + "`\n"
		responseMessage += "Не удалось обновить задачи\n\n"
		failureTableString := &strings.Builder{}
		failureTable := tablewriter.NewWriter(failureTableString)
		failureTable.SetHeader([]string{"Задача"})
		for _, issue := range fillErrors {
			failureTable.Append([]string{issue})
		}
		failureTable.Render()
		responseMessage += "`" + failureTableString.String() + "`\n"
		responseMessage += fmt.Sprintf("Не удалось распределить %d ч.", remain)
	} else {
		responseMessage = fmt.Sprintf("Задачи([%d](%s/time_entries)) успешно обновлены!\n\n", len(issues), host)
		responseMessage += "`" + successTableString.String() + "`"
	}

	return NewCommandResult(responseMessage), nil
}

func (f FillHoursMany) fillIssuesResult(issues []string, hours []string, comment string) (fillSuccess []*redmine.TimeEntryBodyResponse, fillErrors []string, remain int) {
	type Result struct {
		success       *redmine.TimeEntryBodyResponse
		failure       string
		failureRemain int
	}
	resultChan := make(chan Result)
	for i, issue := range issues {
		if i < len(hours) {
			hour := hours[i]
			go func(hour string, issue string) {
				fillHoursResponse, err := f.redmineClient.FillHoursRequest(issue, hour, comment, "")
				if err != nil {
					hourInt, _ := strconv.Atoi(hour)
					resultChan <- Result{nil, issue, hourInt}
					return
				}
				resultChan <- Result{fillHoursResponse, "", 0}
			}(hour, issue)
		}
	}
	for range issues {
		data := <-resultChan
		if data.success != nil {
			fillSuccess = append(fillSuccess, data.success)
		} else {
			remain += data.failureRemain
			fillErrors = append(fillErrors, data.failure)
		}
	}
	sort.Slice(fillSuccess, func(i, j int) bool {
		return fillSuccess[i].TimeEntry.Issue.ID < fillSuccess[j].TimeEntry.Issue.ID
	})
	sort.Slice(fillErrors, func(i, j int) bool {
		return fillErrors[i] < fillErrors[j]
	})
	return
}

func (f FillHoursMany) getHours(issuesCount int, issues []string, remainingHours int) []string {
	var hours []string
	var remainingIssuesCount = issuesCount
	for range issues {
		hour := int(math.Ceil(float64(remainingHours) / float64(remainingIssuesCount)))
		hours = append(hours, fmt.Sprintf("%d", hour))
		remainingHours -= hour
		remainingIssuesCount--
	}
	return hours
}

func (f FillHoursMany) getRemainingHours(timeEntries []*redmine.TimeEntryResponse) (int, error) {
	var storedHours = 0
	for _, entry := range timeEntries {
		storedHours += int(math.Ceil(float64(entry.Hours)))
	}
	const workDayLength = 8
	if storedHours >= workDayLength {
		return 0, errors.New("Вы сегодня уже работали 8 часов")
	}
	remainingHours := workDayLength - storedHours
	return remainingHours, nil
}

func (f FillHoursMany) getIssuesAndComment(message string) ([]string, string, error) {
	if len(message) == 0 {
		return nil, "", errors.New("Введена неправильная команда")
	}

	fragments := strings.Split(message, " ")
	issuesMap := make(map[string]bool)
	var comment string
	for i, fragment := range fragments {
		if regexp.MustCompile(`^#?[0-9]+$`).MatchString(fragment) {
			trimmed := strings.TrimLeft(fragment, "#")
			issuesMap[trimmed] = true
		} else {
			comment = strings.Join(fragments[i:], " ")
			break
		}
	}
	if len(comment) == 0 {
		return nil, "", errors.New("Вы не ввели комментарий для задач")
	}

	var issues []string
	for key := range issuesMap {
		issues = append(issues, key)
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i] < issues[j]
	})

	return issues, comment, nil
}

func (f FillHoursMany) IsCompleted() bool {
	return true
}
