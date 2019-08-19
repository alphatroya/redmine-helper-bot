package commands

import (
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"testing"
)

func TestPartlyFillHoursCommand_Handle(t *testing.T) {
	data := []struct {
		message     string
		isCompleted bool
		isHoursSet  bool
		result      *CommandResult
		err         error
	}{
		{message: "test", isCompleted: true, isHoursSet: true, result: NewCommandResult("Операция выполнена"), err: nil},
	}
	for _, item := range data {
		redmineMock := &RedmineMock{}
		storageMock := &StorageMock{}
		sut := newPartlyFillHoursCommand(redmineMock, storageMock, 1)
		sut.isCompleted = item.isCompleted
		sut.isHoursSet = item.isHoursSet
		result, err := sut.Handle(item.message)
		if result != nil && result.message != item.result.message {
			t.Errorf("wrong result from handle method, got: %s, expected: %s", result.message, item.result.message)
		}
		if err != nil && err != item.err {
			t.Errorf("wrong error from handle method, got: %s, expected: %s", err, item.err)
		}
	}
}

type StorageMock struct {
}

func (s StorageMock) SetToken(token string, chat int64) {
	panic("implement me")
}

func (s StorageMock) GetToken(int64) (string, error) {
	panic("implement me")
}

func (s StorageMock) SetHost(host string, chat int64) {
	panic("implement me")
}

func (s StorageMock) GetHost(chat int64) (string, error) {
	panic("implement me")
}

type RedmineMock struct {
}

func (r RedmineMock) SetToken(token string) {
	panic("implement me")
}

func (r RedmineMock) SetHost(host string) {
	panic("implement me")
}

func (r RedmineMock) FillHoursRequest(issueID string, hours string, comment string) (*redmine.TimeEntryBodyResponse, error) {
	panic("implement me")
}

func (r RedmineMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	panic("implement me")
}

func (r RedmineMock) AssignedIssues() ([]*redmine.Issue, error) {
	panic("implement me")
}
