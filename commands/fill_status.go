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
	return NewCommandResult(fmt.Sprintf("Вы сегодня работали %.1f ч.", sum)), nil
}

func (f FillStatus) IsCompleted() bool {
	return true
}
