package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type AddComment struct {
	redmineClient redmine.Client
	storage       storage.Manager
	printer       redmine.Printer
	chatID        int64
	issueID       string
	updatingIssue *redmine.Issue
	comment       string
	completed     bool
	isReject      bool
}

func NewAddComment(redmineClient redmine.Client, storage storage.Manager, printer redmine.Printer, chatID int64) *AddComment {
	return &AddComment{redmineClient: redmineClient, storage: storage, printer: printer, chatID: chatID}
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
	var responseMessage []string
	result, err := a.redmineClient.Issue(issueID)
	if err == nil {
		a.updatingIssue = result.Issue
		printedIssue := a.printer.Print(*result.Issue, false)
		for _, element := range printedIssue {
			responseMessage = append(responseMessage, element)
		}
		responseMessage = append(responseMessage, "_Напишите комментарий к задаче_")
	} else {
		responseMessage = []string{fmt.Sprintf("Напишите комментарий к задаче [#%s](%s/issues/%s)", issueID, host, issueID)}
	}
	a.issueID = issueID
	return NewCommandResultWithMessages(responseMessage), nil
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
		return NewCommandResultWithKeyboard(
			fmt.Sprintf("Комментарий *не был* добавлен в задачу [#%s](%s/issues/%s) 😞\n\nПовторить запрос?", a.issueID, host, a.issueID),
			[]string{"Да", "Нет"},
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
