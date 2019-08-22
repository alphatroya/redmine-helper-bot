package commands

import (
	"fmt"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

const (
	WrongFillHoursWrongIssueIDResponse    = "Введен неправильный номер задачи"
	WrongFillHoursWrongHoursCountResponse = "Введено неправильное количество часов"
)

func SuccessFillHoursMessageResponse(issueID int, issue *redmine.IssueContainer, hours float32, host string) string {
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
