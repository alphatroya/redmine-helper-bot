package main

import "fmt"

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
	return fmt.Sprintf("В задачу %s добавлено часов: %s (%s/issues/%s/time_entries)", issueID, hours, host, issueID)
}
