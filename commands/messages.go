package commands

import (
	"fmt"
)

const (
	WrongFillHoursWrongIssueIDResponse    = "Введен неправильный номер задачи"
	WrongFillHoursWrongHoursCountResponse = "Введено неправильное количество часов"
)

func SuccessFillHoursMessageResponse(issueID int, hours float32, host string) string {
	return fmt.Sprintf("В задачу [#%d](%s/issues/%d/time_entries) добавлено часов: *%.1f*", issueID, host, issueID, hours)
}
