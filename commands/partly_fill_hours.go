package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

type PartlyFillHoursCommand struct {
	redmineClient   redmine.Client
	printer         redmine.Printer
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
	shortVersion    bool
}

func (p *PartlyFillHoursCommand) IsCompleted() bool {
	return p.isCompleted
}

func newPartlyFillHoursCommand(redmineClient redmine.Client, printer redmine.Printer, storage storage.Manager, chatID int64) *PartlyFillHoursCommand {
	return &PartlyFillHoursCommand{redmineClient: redmineClient, printer: printer, storage: storage, chatID: chatID}
}

func NewFillHoursCommand(redmineClient redmine.Client, printer redmine.Printer, storage storage.Manager, chatID int64, message string) (*PartlyFillHoursCommand, error) {
	command := newPartlyFillHoursCommand(redmineClient, printer, storage, chatID)
	split := strings.Split(message, " ")
	if len(split) < 3 {
		return nil, fmt.Errorf(command.HelpMessage())
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
	err = command.setComment(strings.Join(split[2:], " "))
	if err != nil {
		return nil, err
	}
	command.shortVersion = true
	return command, nil
}

func (p *PartlyFillHoursCommand) Handle(message string) (*CommandResult, error) {
	if p.isCompleted {
		return NewCommandResult("_Операция выполнена_"), nil
	}
	if len(p.comment) > 0 {
		return p.makeFillRequest()
	}
	if p.isHoursSet {
		err := p.setComment(message)
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
	return p.makeIssuesRequest()
}

func (p *PartlyFillHoursCommand) makeIssuesRequest() (*CommandResult, error) {
	issues, err := p.redmineClient.AssignedIssues()
	if err != nil {
		return nil, err
	}
	messages := []string{
		"_Введите номер задачи_",
	}

	var buttons []string
	for _, issue := range issues {
		var subject string
		maxLength := 30
		runes := []rune(issue.Subject)
		if len(runes) <= maxLength {
			subject = issue.Subject
		} else {
			subject = string(runes[:maxLength]) + "..."
		}
		buttons = append(buttons, fmt.Sprintf("#%d - %s", issue.ID, subject))
	}

	p.issuesRequested = true
	return NewCommandResultWithMessagesAndKeyboard(messages, buttons), nil
}

func (p *PartlyFillHoursCommand) setIssueID(issueID string) (*CommandResult, error) {
	issueID = strings.TrimLeft(issueID, "#")
	p.activityID, _ = p.storage.GetActivity(p.chatID)

	if regexp.MustCompile(`^[0-9]+ - .+$`).MatchString(issueID) {
		searchResult := regexp.MustCompile(`^[0-9]+`).Find([]byte(issueID))
		if len(searchResult) != 0 {
			return p.issueIDSuccess(string(searchResult))
		}
	}

	if !regexp.MustCompile(`^[0-9]+$`).MatchString(issueID) {
		return nil, fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}
	return p.issueIDSuccess(issueID)
}

func (p *PartlyFillHoursCommand) issueIDSuccess(issueID string) (*CommandResult, error) {
	p.issueID = issueID
	p.isIssueIDSet = true
	var buttons []string
	for i := 1; i <= 8; i++ {
		buttons = append(buttons, fmt.Sprintf("%d", i))
	}
	successMessage := p.printIssueInfo("Номер задачи установлен, введите число часов", true)
	return NewCommandResultWithMessagesAndKeyboard(successMessage, buttons), nil
}

func (p *PartlyFillHoursCommand) printIssueInfo(message string, italic bool) (successMessage []string) {
	issue, _ := p.redmineClient.Issue(p.issueID)
	if issue != nil {
		successMessage = append(successMessage, p.printer.Print(*issue.Issue, false)...)
		if italic {
			successMessage = append(successMessage, "_"+message+"_")
		} else {
			successMessage = append(successMessage, message)
		}
	} else {
		successMessage = append(successMessage, message)
	}
	return
}

func (p *PartlyFillHoursCommand) setHours(hours string) (*CommandResult, error) {
	_, err := strconv.ParseFloat(hours, 32)
	if err != nil {
		return nil, fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}
	p.hours = hours
	p.isHoursSet = true
	activities, err := p.redmineClient.Activities()
	successMessage := "_Количество часов установлено, введите комментарий_"
	if err != nil {
		return NewCommandResult(successMessage), nil
	}
	activitiesButtons := []string{"Исправление бага"}
	for _, activity := range activities {
		activitiesButtons = append(activitiesButtons, activity.Name)
	}
	return NewCommandResultWithKeyboard(successMessage, activitiesButtons), nil
}

func (p *PartlyFillHoursCommand) setComment(comment string) error {
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return errors.New("Введена пустая команда")
	}
	p.comment = comment
	return nil
}

func (p *PartlyFillHoursCommand) makeFillRequest() (*CommandResult, error) {
	host, _ := p.storage.GetHost(p.chatID)
	requestBody, err := p.redmineClient.FillHoursRequest(p.issueID, p.hours, p.comment, p.activityID)
	if err != nil {
		return nil, err
	}
	p.isCompleted = true
	var successMessage []string
	if p.shortVersion {
		successMessage = p.printIssueInfo(SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, requestBody.TimeEntry.Hours, host), false)
	} else {
		successMessage = []string{SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, requestBody.TimeEntry.Hours, host)}
	}
	return NewCommandResultWithMessages(successMessage), nil
}

func (p *PartlyFillHoursCommand) HelpMessage() string {
	message := `
*Команда служит для быстрого заполнения часов для указанной задачи*

_Синтаксис:_ '/fh <номер задачи> <часы> <комментарий>'

_Пример:_ "/fh 1 8 Исправление" установит значения:

`
	stringBuilder := &strings.Builder{}
	table := tablewriter.NewWriter(stringBuilder)
	table.SetHeader([]string{"ID", "Часы", "Комментарий"})
	table.Append([]string{"1", "8", "Исправление"})
	table.Render()
	return message + "`" + stringBuilder.String() + "`"
}
