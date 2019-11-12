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
		err       string
		completed bool
	}{
		{
			err: "Вы ввели неправильный номер задачи",
		},
		{
			command: "34fdsd111",
			err:     "Вы ввели неправильный номер задачи",
		},
		{
			command: "43213",
			issue:   mockIssue,
			result:  fmt.Sprintf("Добавьте комментарий к задаче #[43213](%s/issues/43213)\n"+issuePrintMessage(), host),
		},
		{
			command: "#43214",
			issue:   mockIssue,
			result:  fmt.Sprintf("Добавьте комментарий к задаче #[43214](%s/issues/43214)\n"+issuePrintMessage(), host),
		},
		{
			command:  "#43214",
			issueErr: fmt.Errorf("error"),
			result:   fmt.Sprintf("Добавьте комментарий к задаче #[43214](%s/issues/43214)\n", host),
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
			if err.Error() != testCase.err {
				t.Errorf("command return wrong error\ngot: %s\nexpected: %s", err, testCase.err)
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
	result := "\n"
	result += "Название\n"
	result += "`+-----------+-------------+\n"
	result += "| СТАТУС    | Сделано     |\n"
	result += "| НАЗНАЧЕНО | Иванов Иван |\n"
	result += "+-----------+-------------+\n"
	result += "`"
	return result
}
