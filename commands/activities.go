package commands

import (
	"fmt"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type Activities struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
}

func newActivitiesCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *Activities {
	return &Activities{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (a Activities) Handle(message string) (*CommandResult, error) {
	token, err := a.storage.GetToken(a.chatID)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	a.redmineClient.SetToken(token)

	host, err := a.storage.GetHost(a.chatID)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	a.redmineClient.SetHost(host)

	activities, err := a.redmineClient.Activities()
	if err != nil {
		return nil, err
	}

	message = "Найдены следующие активности:\n\n"
	for _, activity := range activities {
		message += fmt.Sprintf("%s - *%d*\n", activity.Name, activity.Id)
	}
	return NewCommandResult(message), nil
}

func (a Activities) IsCompleted() bool {
	return true
}
