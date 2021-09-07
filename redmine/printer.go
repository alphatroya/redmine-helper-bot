package redmine

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Printer interface {
	Print(issue Issue, printDescription bool) []string
	PrintIssues(issues []*Issue) []string
}

type TablePrinter struct {
}

func (t TablePrinter) PrintIssues(issues []*Issue) (messages []string) {
	projects := make(map[string][]*Issue)
	for _, issue := range issues {
		projects[issue.Project.Name] = append(projects[issue.Project.Name], issue)
	}
	for key, value := range projects {
		tableBuilder := &strings.Builder{}
		table := tablewriter.NewWriter(tableBuilder)
		for _, issue := range value {
			table.Append([]string{
				fmt.Sprintf("%d", issue.ID),
				issue.Status.Name,
				issue.Subject,
			})
		}
		table.SetCaption(true, key)
		table.Render()
		messages = append(messages, monospaced(tableBuilder.String()))
	}
	return
}

func (t TablePrinter) Print(issue Issue, printDescription bool) []string {
	result := []string{fmt.Sprintf("*ЗАДАЧА #%d*: %s", issue.ID, issue.Subject)}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.Append([]string{"ПРОЕКТ", issue.Project.Name})
	table.Append([]string{"СТАТУС", issue.Status.Name})
	table.Append([]string{"ТРЕКЕР", issue.Tracker.Name})
	table.Append([]string{"ПРИОРИТЕТ", issue.Priority.Name})
	table.Append([]string{"АВТОР", issue.Author.Name})
	table.Append([]string{"НАЗНАЧЕНО", issue.AssignedTo.Name})
	table.Append([]string{"ЧАСЫ", fmt.Sprintf("%.1f", issue.SpentHours)})
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.Render()
	result = append(result, monospaced(tableString.String()))

	if printDescription {
		result = append(result, "ОПИСАНИЕ:\n\n"+issue.Description)
	}
	return result
}

func monospaced(table string) string {
	return "`" + table + "`"
}
