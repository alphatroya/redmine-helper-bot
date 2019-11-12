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
}

func NewAddComment(redmineClient redmine.Client, storage storage.Manager, chatID int64) *AddComment {
	return &AddComment{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (a *AddComment) Handle(message string) (*CommandResult, error) {
	issueID, ok := redmine.CheckAndExtractIssueID(message)
	if !ok {
		return nil, errors.New("Вы ввели неправильный номер задачи")
	}
	host, err := a.storage.GetHost(a.chatID)
	if err != nil {
		return nil, err
	}
	responseMessage := fmt.Sprintf("Добавьте комментарий к задаче #[%s](%s/issues/%s)\n", issueID, host, issueID)
	result, err := a.redmineClient.Issue(issueID)
	if err == nil {
		tableString := &strings.Builder{}
		tableString.WriteString(fmt.Sprintf("\n%s\n`", result.Issue.Subject))
		table := tablewriter.NewWriter(tableString)
		table.Append([]string{fmt.Sprintf("СТАТУС"), result.Issue.Status.Name})
		table.Append([]string{fmt.Sprintf("НАЗНАЧЕНО"), result.Issue.AssignedTo.Name})
		table.Render()
		responseMessage += tableString.String() + "`"
	}
	a.issueID = issueID
	return NewCommandResult(responseMessage), nil
}

func (a *AddComment) IsCompleted() bool {
	return false
}
