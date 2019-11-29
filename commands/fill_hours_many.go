package commands

import (
	"errors"
	"fmt"
	"sort"
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

	hours := f.getHours(issues, remainingHours)

	host, err := f.storage.GetHost(f.chatID)
	if err != nil {
		return nil, err
	}

	fillSuccess, fillErrors, remain := f.fillIssuesResult(issues, hours, comment)

	successTableString := &strings.Builder{}
	successTable := tablewriter.NewWriter(successTableString)
	successTable.SetHeader([]string{"Задача", "Часы"})

	for _, timeEntry := range fillSuccess {
		successTable.Append([]string{fmt.Sprintf("%d", timeEntry.TimeEntry.Issue.ID), fmt.Sprintf("%.1f", timeEntry.TimeEntry.Hours)})
	}
	successTable.Render()

	var responseMessage string
	if len(fillErrors) != 0 {
		responseMessage = fmt.Sprintf("Задачи([%d](%s/time_entries)) *частично* обновлены, обновленные задачи\n\n", len(issues)-len(fillErrors), host)
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
		responseMessage += fmt.Sprintf("Не удалось распределить %.1f ч.", remain)
	} else {
		responseMessage = fmt.Sprintf("Задачи([%d](%s/time_entries)) успешно обновлены!\n\n", len(issues), host)
		responseMessage += "`" + successTableString.String() + "`"
	}

	return NewCommandResult(responseMessage), nil
}

func (f FillHoursMany) fillIssuesResult(issues []string, hours []float64, comment string) (fillSuccess []*redmine.TimeEntryBodyResponse, fillErrors []string, remain float64) {
	type Result struct {
		success       *redmine.TimeEntryBodyResponse
		failure       string
		failureRemain float64
	}
	resultChan := make(chan Result)
	for i, issue := range issues {
		if i < len(hours) {
			hour := hours[i]
			go func(hour float64, issue string) {
				hourString := fmt.Sprintf("%.3f", hour)
				fillHoursResponse, err := f.redmineClient.FillHoursRequest(issue, hourString, comment, "")
				if err != nil {
					resultChan <- Result{nil, issue, hour}
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

func (f FillHoursMany) getHours(issues []string, remainingHours float64) (hours []float64) {
	issuesCount := len(issues)
	hours = make([]float64, issuesCount)
	for i := range issues {
		hours[i] = remainingHours / float64(issuesCount)
	}
	return
}

func (f FillHoursMany) getRemainingHours(timeEntries []*redmine.TimeEntryResponse) (float64, error) {
	var storedHours float64 = 0
	for _, entry := range timeEntries {
		storedHours += float64(entry.Hours)
	}
	const workDayLength = 8.0
	if workDayLength-storedHours <= 0.02 {
		return 0, errors.New("Вы сегодня уже работали 8 часов")
	}
	remainingHours := workDayLength - storedHours
	return remainingHours, nil
}

func (f FillHoursMany) getIssuesAndComment(message string) ([]string, string, error) {
	if len(message) == 0 {
		return nil, "", errors.New(f.HelpMessage())
	}

	fragments := strings.Split(message, " ")
	issuesMap := make(map[string]bool)
	var comment string
	for i, fragment := range fragments {
		if len(fragment) == 0 {
			continue
		}
		if trimmed, ok := redmine.CheckAndExtractIssueID(fragment); ok {
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

func (f FillHoursMany) HelpMessage() string {
	message := `
*Команда служит для распределения свободных незаполненных часов между введенными задачами*

_Синтаксис:_ '/fhm <один или несколько идентификаторов задач разделенных пробелом> <комментарий>'

- Один комментарий будет установлен для всех перечисленных задач
- Число задач не может быть больше числа свободных за сегодня часов

_Пример:_ "/fhm 1 2 3 5 Исправление" при свободных *8 часах* установит значения:

`
	stringBuilder := &strings.Builder{}
	table := tablewriter.NewWriter(stringBuilder)
	table.SetHeader([]string{"ID", "Часы", "Комментарий"})
	table.Append([]string{"1", "2", "Исправление"})
	table.Append([]string{"2", "2", "Исправление"})
	table.Append([]string{"3", "2", "Исправление"})
	table.Append([]string{"5", "2", "Исправление"})
	table.Render()
	return message + "`" + stringBuilder.String() + "`"
}

func (f FillHoursMany) IsCompleted() bool {
	return true
}
