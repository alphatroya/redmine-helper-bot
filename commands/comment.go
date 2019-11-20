package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
)

type AddComment struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
	issueID       string
	updatingIssue *redmine.Issue
	comment       string
	completed     bool
	isReject      bool
}

func NewAddComment(redmineClient redmine.Client, storage storage.Manager, chatID int64) *AddComment {
	return &AddComment{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (a *AddComment) Handle(message string) (*CommandResult, error) {
	host, err := a.storage.GetHost(a.chatID)
	if err != nil {
		return nil, err
	}

	if len(a.issueID) == 0 {
		return a.firstPhase(message, host)
	}
	if len(a.comment) == 0 {
		return a.secondPhase(message, host)
	}
	switch strings.ToLower(message) {
	case "да":
		return a.secondPhase(a.comment, host)
	case "нет":
		a.completed = true
		return NewCommandResult("Операция отменена"), nil
	default:
		return NewCommandResult(`Вы должны написать "да" или "нет"`), nil
	}
}

func (a *AddComment) firstPhase(message string, host string) (*CommandResult, error) {
	issueID, ok := redmine.CheckAndExtractIssueID(message)
	if !ok {
		return nil, errors.New("Вы ввели неправильный номер задачи")
	}
	responseMessage := fmt.Sprintf("Добавьте комментарий к задаче [#%s](%s/issues/%s)", issueID, host, issueID)
	result, err := a.redmineClient.Issue(issueID)
	if err == nil {
		a.updatingIssue = result.Issue
		tableString := &strings.Builder{}
		tableString.WriteString(fmt.Sprintf("\n\n\n*%s*\n\n`", result.Issue.Subject))
		table := tablewriter.NewWriter(tableString)
		table.Append([]string{fmt.Sprintf("СТАТУС"), result.Issue.Status.Name})
		table.Append([]string{fmt.Sprintf("АВТОР"), result.Issue.Author.Name})
		table.Append([]string{fmt.Sprintf("НАЗНАЧЕНО"), result.Issue.AssignedTo.Name})
		table.Render()
		responseMessage += tableString.String() + "`"
	}
	a.issueID = issueID
	return NewCommandResult(responseMessage), nil
}

func (a *AddComment) secondPhase(message string, host string) (*CommandResult, error) {
	message = strings.TrimSpace(message)
	if len(message) == 0 {
		return nil, errors.New("Введен пустой комментарий")
	}
	var err error
	if a.isReject && a.updatingIssue != nil {
		err = a.redmineClient.AddComment(a.issueID, message, a.updatingIssue.Author.ID)
	} else {
		err = a.redmineClient.AddComment(a.issueID, message, 0)
	}
	if err != nil {
		a.comment = message
		return NewCommandResult(
			fmt.Sprintf("Комментарий *не* добавлен в задачу [#%s](%s/issues/%s)\nПовторить запрос? (*да*/нет)", a.issueID, host, a.issueID),
		), nil
	}
	a.completed = true
	message = fmt.Sprintf(
		"Комментарий добавлен в задачу [#%s](%s/issues/%s)",
		a.issueID,
		host,
		a.issueID,
	)
	if a.isReject && a.updatingIssue != nil {
		message += fmt.Sprintf(" и назначен на: %s", a.updatingIssue.Author.Name)
	}
	return NewCommandResult(message), nil
}

func (a *AddComment) IsCompleted() bool {
	return a.completed
}
