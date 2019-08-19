package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type PartlyFillHoursCommand struct {
	redmineClient   redmine.Client
	storage         storage.Manager
	chatID          int64
	issuesRequested bool
	isIssueIDSet    bool
	isHoursSet      bool
	isCompleted     bool
	issueID         string
	hours           string
	comment         string
}

func (p *PartlyFillHoursCommand) IsCompleted() bool {
	return p.isCompleted
}

func newPartlyFillHoursCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *PartlyFillHoursCommand {
	return &PartlyFillHoursCommand{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (p *PartlyFillHoursCommand) Handle(message string) (*CommandResult, error) {
	if p.isCompleted {
		return NewCommandResult("Операция выполнена"), nil
	}
	if p.isHoursSet {
		return p.setComment(message)
	}
	if p.isIssueIDSet {
		return p.saveHours(message)
	}
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
	message += fmt.Sprintln("*Введите номер задачи*")
	message += fmt.Sprintln("-----------------------------")
	message += fmt.Sprintln("")
	message += fmt.Sprintln("_Вы можете выбрать номер из списка снизу или ввести свой_")
	message += fmt.Sprintln("")
	projects := make(map[string][]*redmine.Issue)
	for _, issue := range issues {
		projects[issue.Project.Name] = append(projects[issue.Project.Name], issue)
	}
	for key, value := range projects {
		message += fmt.Sprintf("*%s*\n", key)
		for _, issue := range value {
			message += fmt.Sprintf("    *#%d* %s\n", issue.ID, issue.Subject)
		}
		message += fmt.Sprintln("")
	}
	p.issuesRequested = true
	return NewCommandResult(message), err
}

func (p *PartlyFillHoursCommand) saveIssueID(issueID string) (*CommandResult, error) {
	issueID = strings.TrimLeft(issueID, "#")
	regex := regexp.MustCompile(`^[0-9]+$`)
	if !regex.MatchString(issueID) {
		return nil, fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}
	p.issueID = issueID
	p.isIssueIDSet = true
	return NewCommandResult("Номер задачи установлен, введите число часов"), nil
}

func (p *PartlyFillHoursCommand) saveHours(hours string) (*CommandResult, error) {
	_, err := strconv.ParseFloat(hours, 32)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}
	p.hours = hours
	p.isHoursSet = true
	return NewCommandResult("Количество часов установлено, введите комментарий"), nil
}

func (p *PartlyFillHoursCommand) setComment(comment string) (*CommandResult, error) {
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return NewCommandResult("Введена пустая команда"), nil
	}
	p.isCompleted = true
	p.comment = comment
	plainCommand := newFillHoursCommand(p.storage, p.chatID, p.redmineClient)
	command := []string{p.issueID, p.hours, p.comment}
	return plainCommand.Handle(strings.Join(command, " "))
}
