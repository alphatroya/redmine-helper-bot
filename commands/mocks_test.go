package commands

import "github.com/alphatroya/redmine-helper-bot/redmine"

type RedmineMock struct {
	mockActivities  []*redmine.Activities
	mockTimeEntries []*redmine.TimeEntryResponse
	err             error
}

func (r RedmineMock) TodayTimeEntries() ([]*redmine.TimeEntryResponse, error) {
	return r.mockTimeEntries, r.err
}

func (r RedmineMock) Activities() ([]*redmine.Activities, error) {
	return r.mockActivities, r.err
}

func (r RedmineMock) SetToken(token string) {
}

func (r RedmineMock) SetHost(host string) {
}

func (r RedmineMock) FillHoursRequest(issueID string, hours string, comment string, activityID string) (*redmine.TimeEntryBodyResponse, error) {
	return &redmine.TimeEntryBodyResponse{}, nil
}

func (r RedmineMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	return &redmine.IssueContainer{
		Issue: &redmine.Issue{},
	}, nil
}

func (r RedmineMock) AssignedIssues() ([]*redmine.Issue, error) {
	return nil, nil
}
