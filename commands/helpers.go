package commands

import (
	"fmt"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

const StandardWrapValue = 20

func wrap(line string, limit int) string {
	if limit <= 3 {
		return line
	}
	wrapped := []rune(line)
	if len(wrapped) > limit {
		wrapped = append(wrapped[:limit-3], []rune("...")...)
	}
	return string(wrapped)
}

func makeIssuesRequest(redmineClient redmine.Client) (*CommandResult, error) {
	issues, err := redmineClient.AssignedIssues()
	if err != nil {
		return nil, err
	}
	messages := []string{
		"_Введите номер задачи_",
	}

	var buttons []string
	for _, issue := range issues {
		var subject string
		maxLength := 30
		runes := []rune(issue.Subject)
		if len(runes) <= maxLength {
			subject = issue.Subject
		} else {
			subject = string(runes[:maxLength]) + "..."
		}
		buttons = append(buttons, fmt.Sprintf("#%d - %s", issue.ID, subject))
	}

	return NewCommandResultWithMessagesAndKeyboard(messages, buttons), nil
}
