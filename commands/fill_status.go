package commands

import (
	"fmt"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
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
		message += "`Часы | Задача | Активность | Комментарий`\n"
		message += "`-----+--------+------------+---------------------`\n"
		for _, entry := range timeEntries {
			message += fmt.Sprintf("` %.1f | %d  | %-10s | %-20s\n`", entry.Hours, entry.Issue.ID, string([]rune(entry.Activity.Name)[:10]), string([]rune(entry.Comments)[:20]))
		}
	}
	return NewCommandResult(message), nil
}

func (f FillStatus) IsCompleted() bool {
	return true
}
