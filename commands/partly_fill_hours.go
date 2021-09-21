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

type FillHoursCommand struct {
	redmineClient redmine.Client
	printer       Printer
	storage       storage.Manager
	chatID        int64
	isIssueIDSet  bool
	issueID       string
	activityID    string
	hours         string
	comment       string
}

func (p *FillHoursCommand) IsCompleted() bool {
	return true
}

func NewFillHoursCommand(redmineClient redmine.Client, printer Printer, storage storage.Manager, chatID int64, message string) (*FillHoursCommand, error) {
	command := &FillHoursCommand{redmineClient: redmineClient, printer: printer, storage: storage, chatID: chatID}
	components := strings.Split(message, " ")
	if len(components) < 3 {
		return nil, fmt.Errorf(command.HelpMessage())
	}
	_, err := command.setIssueID(components[0])
	if err != nil {
		return nil, err
	}
	_, err = command.setHours(components[1])
	if err != nil {
		return nil, err
	}
	err = command.setComment(strings.Join(components[2:], " "))
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (p *FillHoursCommand) Handle(message string) (*CommandResult, error) {
	return p.makeFillRequest()
}

func (p *FillHoursCommand) setIssueID(issueID string) (*CommandResult, error) {
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

func (p *FillHoursCommand) issueIDSuccess(issueID string) (*CommandResult, error) {
	p.issueID = issueID
	p.isIssueIDSet = true
	var buttons []string
	for i := 1; i <= 8; i++ {
		buttons = append(buttons, fmt.Sprintf("%d", i))
	}
	successMessage := p.printIssueInfo("Номер задачи установлен, введите число часов", true)
	return NewCommandResultWithMessagesAndKeyboard(successMessage, buttons), nil
}

func (p *FillHoursCommand) printIssueInfo(message string, italic bool) (successMessage []string) {
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

func (p *FillHoursCommand) setHours(hours string) (*CommandResult, error) {
	if _, err := strconv.ParseFloat(hours, 32); err != nil {
		return nil, fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}
	p.hours = hours
	successMessage := "_Количество часов установлено, введите комментарий_"
	activities, err := p.redmineClient.Activities()
	if err != nil {
		return NewCommandResult(successMessage), nil
	}
	activitiesButtons := []string{"Исправление бага"}
	for _, activity := range activities {
		activitiesButtons = append(activitiesButtons, activity.Name)
	}
	return NewCommandResultWithKeyboard(successMessage, activitiesButtons), nil
}

func (p *FillHoursCommand) setComment(comment string) error {
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return errors.New("Введена пустая команда")
	}
	p.comment = comment
	return nil
}

func (p *FillHoursCommand) makeFillRequest() (*CommandResult, error) {
	host, _ := p.storage.GetHost(p.chatID)
	requestBody, err := p.redmineClient.FillHoursRequest(p.issueID, p.hours, p.comment, p.activityID)
	if err != nil {
		return nil, err
	}
	successMessage := p.printIssueInfo(SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, requestBody.TimeEntry.Hours, host), false)
	return NewCommandResultWithMessages(successMessage), nil
}

func (p *FillHoursCommand) HelpMessage() string {
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
