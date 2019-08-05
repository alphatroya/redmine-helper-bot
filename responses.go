package main

import (
	"fmt"
	"strconv"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

const (
	WrongTokenMessageResponse                    = "Неправильное количество аргументов"
	SuccessTokenMessageResponse                  = "Токен успешно обновлен"
	WrongHostMessageResponse                     = "Неправильное количество аргументов"
	SuccessHostMessageResponse                   = "Адрес сервера успешно обновлен"
	WrongFillHoursTokenNilResponse               = "Токен доступа для текущего пользователя не найден"
	WrongFillHoursHostNilResponse                = "Адрес сервера не найден"
	WrongFillHoursWrongIssueIDResponse           = "Введен неправильный номер задачи"
	WrongFillHoursWrongHoursCountResponse        = "Введено неправильное количество часов"
	WrongFillHoursWrongStatusCodeResponse        = "Wrong response from redmine server %d - %s"
	WrongFillHoursWrongNumberOfArgumentsResponse = "Неправильное количество аргументов"
	UnknownCommandResponse                       = "Введена неправильная команда"
)

func SuccessFillHoursMessageResponse(issueID string, issue *redmine.Issue, hours string, host string) string {
	message := fmt.Sprintf("В задачу [#%s](%s/issues/%s/time_entries) добавлено часов: *%s*\n", issueID, host, issueID, hours)
	if issue != nil {
		message += "\n"
		message += fmt.Sprintf("Задача #%d", issue.Issue.ID)
		message += "\n"
	}

	number, _ := strconv.ParseInt(hours, 10, 64)
	if number > 0 {
		message += "\n"
		message += "_Сколько это Джонов Траволт? А вот сколько:_\n"
		message += "\n"
		var i int64
		for i = 0; i < number; i++ {
			message += "🕺"
		}
		message += "\n"
	}
	return message
}
