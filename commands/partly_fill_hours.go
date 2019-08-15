package commands

import (
	"fmt"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type PartlyFillHoursCommand struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
}

func NewPartlyFillHoursCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *PartlyFillHoursCommand {
	return &PartlyFillHoursCommand{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (p PartlyFillHoursCommand) Handle(message string) (string, error) {
	return p.makeIssuesRequest(message)
}

func (p PartlyFillHoursCommand) makeIssuesRequest(message string) (string, error) {
	token, err := p.storage.GetToken(p.chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	p.redmineClient.SetToken(token)
	host, err := p.storage.GetHost(p.chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	p.redmineClient.SetHost(host)
	issues, err := p.redmineClient.AssignedIssues()
	if err != nil {
		return "", err
	}
	for _, issue := range issues {
		message += fmt.Sprintf("[%d](tg://fillhours %d): %s\n", issue.ID, issue.ID, issue.Subject)
	}
	return message, err
}
