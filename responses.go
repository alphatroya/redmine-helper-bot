package main

import (
	"fmt"
	"strconv"
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
	WrongFillHoursWrongNumberOfArgumentsResponse = "Неправильное количество аргументов"
	UnknownCommandResponse                       = "Введена неправильная команда"
)

func SuccessFillHoursMessageResponse(issueID string, hours string, host string) string {
	message := fmt.Sprintf("В задачу [#%s](%s/issues/%s/time_entries) добавлено часов: *%s*", issueID, host, issueID, hours)
	number, _ := strconv.ParseInt(hours, 10, 64)
	if number > 0 {
		message += " "
		var i int64
		for i = 0; i < number; i++ {
			message += "🕺"
		}
	}
	return message
}
