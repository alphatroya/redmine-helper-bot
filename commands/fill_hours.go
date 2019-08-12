package commands

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"regexp"
	"strconv"
	"strings"
)

const (
	WrongFillHoursTokenNilResponse               = "Токен доступа для текущего пользователя не найден"
	WrongFillHoursHostNilResponse                = "Адрес сервера не найден"
	WrongFillHoursWrongIssueIDResponse           = "Введен неправильный номер задачи"
	WrongFillHoursWrongHoursCountResponse        = "Введено неправильное количество часов"
	WrongFillHoursWrongNumberOfArgumentsResponse = "Неправильное количество аргументов"
)

func SuccessFillHoursMessageResponse(issueID int, issue *redmine.Issue, hours float32, host string) string {
	message := fmt.Sprintf("В задачу [#%d](%s/issues/%d/time_entries) добавлено часов: *%.1f*\n", issueID, host, issueID, hours)
	if issue != nil {
		message += "\n"
		message += fmt.Sprintf("*Задача* %s", issue.Issue.Subject)
		message += "\n"
		message += fmt.Sprintf("*Статус* %s\n", issue.Issue.Status.Name)
		message += fmt.Sprintf("*Автор* %s\n", issue.Issue.Author.Name)
		message += fmt.Sprintf("*Назначена* %s\n", issue.Issue.AssignedTo.Name)
		message += fmt.Sprintf("*Всего часов* %.1f\n", issue.Issue.SpentHours)

		message += "\n"
		message += fmt.Sprintf("_%s_\n", issue.Issue.Description)
	}
	return message
}

type FillHoursCommand struct {
	storage       storage.Manager
	chatID        int64
	redmineClient redmine.Client
}

func NewFillHoursCommand(storage storage.Manager, chatID int64, redmineClient redmine.Client) *FillHoursCommand {
	return &FillHoursCommand{storage: storage, chatID: chatID, redmineClient: redmineClient}
}

func (f FillHoursCommand) Handle(message string) (string, error) {
	token, err := f.storage.GetToken(f.chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}
	f.redmineClient.SetToken(token)

	host, err := f.storage.GetHost(f.chatID)
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}
	f.redmineClient.SetHost(host)

	splitted := strings.Split(message, " ")
	if len(splitted) < 3 {
		return "", fmt.Errorf(WrongFillHoursWrongNumberOfArgumentsResponse)
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	issueID := splitted[0]
	if !regex.MatchString(issueID) {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}

	_, conversionError := strconv.ParseFloat(splitted[1], 32)
	if conversionError != nil {
		return "", fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}

	requestBody, err := f.redmineClient.FillHoursRequest(issueID, splitted[1], strings.Join(splitted[2:], " "))
	if err != nil {
		return "", err
	}

	issue, _ := f.redmineClient.Issue(issueID)

	return SuccessFillHoursMessageResponse(requestBody.TimeEntry.Issue.ID, issue, requestBody.TimeEntry.Hours, host), nil
}

func (f FillHoursCommand) Cancel() {
}
