package commands

import (
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
)

type FillStatus struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
}

func NewFillStatus(redmineClient redmine.Client, storage storage.Manager, chatID int64) *FillStatus {
	return &FillStatus{redmineClient: redmineClient, storage: storage, chatID: chatID}
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
	message = fmt.Sprintf(`Вы сегодня работали *%.1f ч.*`, sum)

	if len(timeEntries) > 0 {
		message += "\n\n"

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
		message += "`" + tableString.String() + "`"
	}
	return NewCommandResult(message), nil
}

func (f FillStatus) IsCompleted() bool {
	return true
}
