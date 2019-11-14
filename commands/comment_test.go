package commands

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestAddComment_Handle(t *testing.T) {
	host := "https://google.com"
	mockIssue := &redmine.Issue{
		AssignedTo: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			Name: "Иванов Иван",
		},
		Author: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			Name: "Сидоров Лев",
		},
		Status: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			Name: "Сделано",
		},
		Subject: "Название",
	}

	testData := []struct {
		command   string
		result    string
		issue     *redmine.Issue
		issueErr  error
		resultErr string
		completed bool
	}{
		{
			resultErr: "Вы ввели неправильный номер задачи",
		},
		{
			command:   "34fdsd111",
			resultErr: "Вы ввели неправильный номер задачи",
		},
		{
			command: "43213",
			issue:   mockIssue,
			result:  fmt.Sprintf("Добавьте комментарий к задаче [#43213](%s/issues/43213)"+issuePrintMessage(), host),
		},
		{
			command: "#43214",
			issue:   mockIssue,
			result:  fmt.Sprintf("Добавьте комментарий к задаче [#43214](%s/issues/43214)"+issuePrintMessage(), host),
		},
		{
			command:  "#43214",
			issueErr: fmt.Errorf("error"),
			result:   fmt.Sprintf("Добавьте комментарий к задаче [#43214](%s/issues/43214)", host),
		},
	}

	for _, testCase := range testData {
		storageMock := storage.NewStorageMock()
		var chatID int64 = 5
		storageMock.SetHost(host, chatID)

		redmineMock := &RedmineMock{}
		if testCase.issue != nil {
			redmineMock.mockIssue = &redmine.IssueContainer{Issue: testCase.issue}
		}
		if testCase.issueErr != nil {
			redmineMock.mockIssueErr = testCase.issueErr
		}

		command := NewAddComment(redmineMock, storageMock, chatID)
		result, err := command.Handle(testCase.command)
		completed := command.IsCompleted()

		if completed != testCase.completed {
			t.Errorf("completed status is not same to expected, got: %t, expected: %t", completed, testCase.completed)
		}

		if err != nil {
			if err.Error() != testCase.resultErr {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, testCase.resultErr)
			}
			continue
		}

		if len(result.buttons) != 0 {
			t.Error("success command should not return buttons")
		}

		if result.Message() != testCase.result {
			t.Errorf("command: %s\nreturn wrong message\ngot: \"%s\"\nexpected: \"%s\"", testCase.command, result.Message(), testCase.result)
		}
	}
}

func TestAddComment_Handle_Phase2(t *testing.T) {
	host := "https://google.com"
	issueID := "43213"
	mockIssue := &redmine.Issue{
		AssignedTo: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			Name: "Иванов Иван",
		},
		Status: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			Name: "Сделано",
		},
		Subject: "Название",
	}

	testData := []struct {
		command       string
		result        string
		resultErr     string
		addCommentErr error
		completed     bool
	}{
		{
			command:   "Test",
			result:    fmt.Sprintf("Комментарий добавлен в задачу [#%s](%s/issues/%s)", issueID, host, issueID),
			completed: true,
		},
		{
			command:   "  ",
			resultErr: "Введен пустой комментарий",
			completed: false,
		},
		{
			command:   "",
			resultErr: "Введен пустой комментарий",
			completed: false,
		},
	}

	for _, testCase := range testData {
		storageMock := storage.NewStorageMock()
		var chatID int64 = 5
		storageMock.SetHost(host, chatID)

		redmineMock := &RedmineMock{}
		redmineMock.mockIssue = &redmine.IssueContainer{Issue: mockIssue}

		command := NewAddComment(redmineMock, storageMock, chatID)
		_, _ = command.Handle(issueID)
		result, err := command.Handle(testCase.command)
		completed := command.IsCompleted()

		if completed != testCase.completed {
			t.Errorf("completed status is not same to expected, got: %t, expected: %t", completed, testCase.completed)
		}

		if err != nil {
			if err.Error() != testCase.resultErr {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, testCase.resultErr)
			}
			continue
		}

		if len(result.buttons) != 0 {
			t.Error("success command should not return buttons")
		}

		if result.Message() != testCase.result {
			t.Errorf("command: %s\nreturn wrong message\ngot: \"%s\"\nexpected: \"%s\"", testCase.command, result.Message(), testCase.result)
		}
	}
}

func issuePrintMessage() string {
	result := "\n\n\n"
	result += "*Название*\n\n"
	result += "`+-----------+-------------+\n"
	result += "| СТАТУС    | Сделано     |\n"
	result += "| АВТОР     | Сидоров Лев |\n"
	result += "| НАЗНАЧЕНО | Иванов Иван |\n"
	result += "+-----------+-------------+\n"
	result += "`"
	return result
}
