package commands

import (
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/olekukonko/tablewriter"
)

type FillStatus struct {
	redmineClient redmine.Client
	chatID        int64
}

func NewFillStatus(redmineClient redmine.Client, chatID int64) *FillStatus {
	return &FillStatus{redmineClient: redmineClient, chatID: chatID}
}

func (f FillStatus) Handle(message string) (*CommandResult, error) {
	timeEntries, err := f.redmineClient.TodayTimeEntries()
	if err != nil {
		return nil, err
	}

	var sum float32
	for _, timeEntry := range timeEntries {
		sum += timeEntry.Hours
	}
	resultMessage := fmt.Sprintf(`Вы сегодня работали *%.1f ч.*`, sum)

	if len(timeEntries) > 0 {
		resultMessage += "\n\n"

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetHeader([]string{"Часы", "Задача", "Активность", "Комментарий"})

		for _, entry := range timeEntries {
			data := []string{
				fmt.Sprintf("%.1f", entry.Hours),
				fmt.Sprintf("%d", entry.Issue.ID),
				wrap(entry.Activity.Name, 10),
				wrap(entry.Comments, 20),
			}
			table.Append(data)
		}
		table.Render()
		resultMessage += "`" + tableString.String() + "`"
	}
	return NewCommandResult(resultMessage), nil
}

func (f FillStatus) IsCompleted() bool {
	return true
}
