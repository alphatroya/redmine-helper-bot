package commands

import (
	"errors"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

func TestMyIssues_Handle(t *testing.T) {
	testCases := []struct {
		issues    []*redmine.Issue
		issuesErr error
		messages  []string
	}{
		{
			issues: []*redmine.Issue{
				{Subject: "Test1"},
				{Subject: "Test2"},
			},
			messages: []string{
				"Test1",
				"Test2",
			},
		},
		{
			issuesErr: errors.New("mock error"),
		},
	}
	for _, testCase := range testCases {
		redmineMock := &RedmineMock{}
		redmineMock.mockAssignedIssues = testCase.issues
		redmineMock.mockAssignedIssuesErr = testCase.issuesErr
		printerMock := PrinterMock{}
		sut := NewMyIssuesCommand(redmineMock, printerMock)
		messages, err := sut.Handle("")
		if testCase.issuesErr != nil {
			if testCase.issuesErr != err {
				t.Errorf("mock errors not the same, got: %s, expected: %s", err, testCase.issuesErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("success command should not return error, got: %s", err)
		}
		for i, expected := range testCase.messages {
			got := messages.Messages()[i]
			if got != expected {
				t.Errorf("result messages check failed, got: %s, expected: %s", got, expected)
			}
		}
	}
}

func TestMyIssues_IsCompleted(t *testing.T) {
	sut := NewMyIssuesCommand(nil, nil)
	if !sut.IsCompleted() {
		t.Error("my issues command should always be completed")
	}
}
