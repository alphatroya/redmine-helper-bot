package commands

import (
	"fmt"
	"strconv"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

type RedmineMock struct {
	mockActivities        []*redmine.Activities
	mockTimeEntries       []*redmine.TimeEntryResponse
	err                   error
	fillHoursSuccessCount *struct{ count int }
	filledIssues          []string
}

func (r *RedmineMock) TodayTimeEntries() ([]*redmine.TimeEntryResponse, error) {
	return r.mockTimeEntries, r.err
}

func (r *RedmineMock) Activities() ([]*redmine.Activities, error) {
	return r.mockActivities, r.err
}

func (r *RedmineMock) SetToken(token string) {
}

func (r *RedmineMock) SetHost(host string) {
}

func (r *RedmineMock) FillHoursRequest(issueID string, hours string, comment string, activityID string) (*redmine.TimeEntryBodyResponse, error) {
	if r.fillHoursSuccessCount != nil {
		if r.fillHoursSuccessCount.count != 0 {
			r.fillHoursSuccessCount.count--
			return r.mockResponse(issueID, hours)
		} else {
			return nil, fmt.Errorf("mock error")
		}
	}
	return r.mockResponse(issueID, hours)
}

func (r *RedmineMock) mockResponse(issueID string, hours string) (*redmine.TimeEntryBodyResponse, error) {
	r.filledIssues = append(r.filledIssues, issueID)
	hoursInt, _ := strconv.Atoi(hours)
	intIssueID, _ := strconv.Atoi(issueID)
	timeEntry := redmine.TimeEntryResponse{
		Hours: float32(hoursInt),
		Issue: redmine.TimeEntryResponseIssue{ID: intIssueID},
	}
	return &redmine.TimeEntryBodyResponse{
		TimeEntry: timeEntry,
	}, nil
}

func (r *RedmineMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	return &redmine.IssueContainer{
		Issue: &redmine.Issue{},
	}, nil
}

func (r *RedmineMock) AssignedIssues() ([]*redmine.Issue, error) {
	return nil, nil
}
