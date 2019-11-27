package commands

import (
	"github.com/alphatroya/redmine-helper-bot/redmine"
)

type MyIssues struct {
	redmineClient redmine.Client
	printer       redmine.Printer
}

func NewMyIssuesCommand(redmineClient redmine.Client, printer redmine.Printer) *MyIssues {
	return &MyIssues{redmineClient: redmineClient, printer: printer}
}

func (m *MyIssues) Handle(message string) (*CommandResult, error) {
	issues, err := m.redmineClient.AssignedIssues()
	if err != nil {
		return nil, err
	}
	messages := m.printer.PrintIssues(issues)
	return NewCommandResultWithMessages(messages), nil
}

func (m *MyIssues) IsCompleted() bool {
	return true
}
