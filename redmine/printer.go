package redmine

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Printer interface {
	Print(issue Issue, printDescription bool) []string
}

type TablePrinter struct {
}

func (t TablePrinter) Print(issue Issue, printDescription bool) []string {
	result := []string{fmt.Sprintf("*ЗАДАЧА #%d*: %s", issue.ID, issue.Subject)}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.Append([]string{fmt.Sprintf("ПРОЕКТ"), issue.Project.Name})
	table.Append([]string{fmt.Sprintf("СТАТУС"), issue.Status.Name})
	table.Append([]string{fmt.Sprintf("ТРЕКЕР"), issue.Tracker.Name})
	table.Append([]string{fmt.Sprintf("ПРИОРИТЕТ"), issue.Priority.Name})
	table.Append([]string{fmt.Sprintf("АВТОР"), issue.Author.Name})
	table.Append([]string{fmt.Sprintf("НАЗНАЧЕНО"), issue.AssignedTo.Name})
	table.Append([]string{fmt.Sprintf("ЧАСЫ"), fmt.Sprintf("%.1f", issue.SpentHours)})
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.Render()
	result = append(result, "`"+tableString.String()+"`")

	if printDescription {
		result = append(result, "ОПИСАНИЕ:\n\n"+issue.Description)
	}
	return result
}
