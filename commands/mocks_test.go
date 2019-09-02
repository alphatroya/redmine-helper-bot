package commands

import "github.com/alphatroya/redmine-helper-bot/redmine"

type StorageMock struct {
}

func (s StorageMock) ResetData(chat int64) error {
	panic("implement me")
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
	mockActivities []*redmine.Activities
	err            error
}

func (r RedmineMock) Activities() ([]*redmine.Activities, error) {
	return r.mockActivities, r.err
}

func (r RedmineMock) SetToken(token string) {
	panic("implement me")
}

func (r RedmineMock) SetHost(host string) {
	panic("implement me")
}

func (r RedmineMock) FillHoursRequest(issueID string, hours string, comment string, activityID string) (*redmine.TimeEntryBodyResponse, error) {
	panic("implement me")
}

func (r RedmineMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	panic("implement me")
}

func (r RedmineMock) AssignedIssues() ([]*redmine.Issue, error) {
	panic("implement me")
}
