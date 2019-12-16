package commands

import (
	"fmt"
)

const (
	WrongFillHoursWrongIssueIDResponse    = "Введен неправильный номер задачи"
	WrongFillHoursWrongHoursCountResponse = "Введено неправильное количество часов"
)

func SuccessFillHoursMessageResponse(issueID int, hours float32, host string) string {
	return fmt.Sprintf("_В задачу_ [#%d](%s/issues/%d/time_entries) _добавлено часов:_ *%.1f*", issueID, host, issueID, hours)
}
