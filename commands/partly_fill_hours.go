package commands

import (
	"errors"
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
	activityID      string
	hours           string
	comment         string
}

func (p *PartlyFillHoursCommand) IsCompleted() bool {
	return p.isCompleted
}

func newPartlyFillHoursCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *PartlyFillHoursCommand {
	return &PartlyFillHoursCommand{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func NewFillHoursCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64, message string) (*PartlyFillHoursCommand, error) {
	command := newPartlyFillHoursCommand(redmineClient, storage, chatID)
	split := strings.Split(message, " ")
	if len(split) < 3 {
		return nil, fmt.Errorf("Введена неправильная команда")
	}
	command.issuesRequested = true
	_, err := command.setIssueID(split[0])
	if err != nil {
		return nil, err
	}
	_, err = command.setHours(split[1])
	if err != nil {
		return nil, err
	}
	_, err = command.setComment(strings.Join(split[2:], " "))
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (p *PartlyFillHoursCommand) Handle(message string) (*CommandResult, error) {
	if p.isCompleted {
		return NewCommandResult("Операция выполнена"), nil
	}
	if len(p.comment) > 0 {
		return p.makeFillRequest()
	}
	if p.isHoursSet {
		_, err := p.setComment(message)
		if err != nil {
			return nil, err
		}
		return p.makeFillRequest()
	}
	if p.isIssueIDSet {
		return p.setHours(message)
	}
	if p.issuesRequested {
		return p.setIssueID(message)
	}
	return p.makeIssuesRequest(message)
}

func (p *PartlyFillHoursCommand) makeIssuesRequest(message string) (*CommandResult, error) {
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
	message += fmt.Sprintln("_Вы можете также ввести через пробел код активности, если хотите установить ее отличной от дефолтной (например \"54323 15\")._")
	message += fmt.Sprintln("_Список кодов можно получить с помощью команды /activities_")
	p.issuesRequested = true
	var buttons []string
	for _, issue := range issues {
		buttons = append(buttons, fmt.Sprintf("#%d - %s", issue.ID, issue.Subject))
	}
	return NewCommandResultWithKeyboard(message, buttons), err
}

func (p *PartlyFillHoursCommand) setIssueID(issueID string) (*CommandResult, error) {
	issueID = strings.TrimLeft(issueID, "#")

	if regexp.MustCompile(`^[0-9]+ - .+$`).MatchString(issueID) {
		searchResult := regexp.MustCompile(`^[0-9]+`).Find([]byte(issueID))
		if len(searchResult) != 0 {
			return p.issueIDSuccess(issueID)
		}
	}

	parts := strings.Split(issueID, " ")
	if len(parts) == 2 {
		p.activityID = parts[1]
		issueID = parts[0]
	}

	if !regexp.MustCompile(`^[0-9]+$`).MatchString(issueID) {
		return nil, fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}
	return p.issueIDSuccess(issueID)
}

func (p *PartlyFillHoursCommand) issueIDSuccess(issueID string) (*CommandResult, error) {
	p.issueID = issueID
	p.isIssueIDSet = true
	var hourButtons []string
	for i := 1; i <= 8; i++ {
		hourButtons = append(hourButtons, fmt.Sprintf("%d", i))
	}
	return NewCommandResultWithKeyboard("Номер задачи установлен, введите число часов", hourButtons), nil
}

func (p *PartlyFillHoursCommand) setHours(hours string) (*CommandResult, error) {
	_, err := strconv.ParseFloat(hours, 32)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}
	p.hours = hours
	p.isHoursSet = true
	activities, err := p.redmineClient.Activities()
	if err != nil {
		return NewCommandResult("Количество часов установлено, введите комментарий"), nil
	}
	activitiesButtons := []string{"Исправление бага"}
	for _, activity := range activities {
		activitiesButtons = append(activitiesButtons, activity.Name)
	}
	return NewCommandResultWithKeyboard("Количество часов установлено, введите комментарий", activitiesButtons), nil
}

func (p *PartlyFillHoursCommand) setComment(comment string) (*CommandResult, error) {
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return nil, errors.New("Введена пустая команда")
	}
	p.comment = comment
	return NewCommandResult("Комментарий сохранен"), nil
}

func (p *PartlyFillHoursCommand) makeFillRequest() (*CommandResult, error) {
	host, _ := p.storage.GetHost(p.chatID)
	requestBody, err := p.redmineClient.FillHoursRequest(p.issueID, p.hours, p.comment, p.activityID)
	if err != nil {
		return nil, err
	}
	issue, _ := p.redmineClient.Issue(p.issueID)
	p.isCompleted = true
	return NewCommandResult(SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, issue, requestBody.TimeEntry.Hours, host)), nil
}
