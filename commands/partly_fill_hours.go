package commands

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"regexp"
)

type PartlyFillHoursCommand struct {
	redmineClient   redmine.Client
	storage         storage.Manager
	chatID          int64
	issuesRequested bool
	isIssueIDSet    bool
	issueID         string
}

func NewPartlyFillHoursCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *PartlyFillHoursCommand {
	return &PartlyFillHoursCommand{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (p *PartlyFillHoursCommand) Handle(message string) (*CommandResult, error) {
	if p.issuesRequested {
		return p.saveIssueID(message)
	}
	return p.makeIssuesRequest(message)
}

func (p *PartlyFillHoursCommand) makeIssuesRequest(message string) (*CommandResult, error) {
	token, err := p.storage.GetToken(p.chatID)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	p.redmineClient.SetToken(token)
	host, err := p.storage.GetHost(p.chatID)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	p.redmineClient.SetHost(host)
	issues, err := p.redmineClient.AssignedIssues()
	if err != nil {
		return nil, err
	}
	for _, issue := range issues {
		message += fmt.Sprintf("%d: %s\n", issue.ID, issue.Subject)
	}
	p.issuesRequested = true
	return NewCommandResult(message), err
}

func (p *PartlyFillHoursCommand) saveIssueID(issueID string) (*CommandResult, error) {
	regex := regexp.MustCompile(`^[0-9]+$`)
	if !regex.MatchString(issueID) {
		return nil, fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}
	p.issueID = issueID
	p.isIssueIDSet = true
	return NewCommandResult("Номер issue id установлен"), nil
}
