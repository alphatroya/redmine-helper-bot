package commands

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/alphatroya/redmine-helper-bot/redmine"
)

type RedmineMock struct {
	sync.RWMutex
	mockActivities        []*redmine.Activities
	mockTimeEntries       []*redmine.TimeEntryResponse
	err                   error
	fillHoursErrorsMap    map[string]bool
	filledIssues          []string
	mockIssue             *redmine.IssueContainer
	mockIssueErr          error
	mockAddCommentError   error
	mockAssignedIssues    []*redmine.Issue
	mockAssignedIssuesErr error
}

func (r *RedmineMock) AddComment(issueID string, comment string, assignedTo int) error {
	return r.mockAddCommentError
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
	if isError := r.fillHoursErrorsMap[issueID]; isError {
		return nil, fmt.Errorf("mock error")
	}
	return r.mockResponse(issueID, hours)
}

func (r *RedmineMock) mockResponse(issueID string, hours string) (*redmine.TimeEntryBodyResponse, error) {
	r.Lock()
	r.filledIssues = append(r.filledIssues, issueID)
	r.Unlock()
	floatHours, _ := strconv.ParseFloat(hours, 64)
	intIssueID, _ := strconv.Atoi(issueID)
	timeEntry := redmine.TimeEntryResponse{
		Hours: float32(floatHours),
		Issue: redmine.TimeEntryResponseIssue{ID: intIssueID},
	}
	return &redmine.TimeEntryBodyResponse{
		TimeEntry: timeEntry,
	}, nil
}

func (r *RedmineMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	if r.mockIssueErr != nil {
		return nil, r.mockIssueErr
	}
	return r.mockIssue, nil
}

func (r *RedmineMock) AssignedIssues() ([]*redmine.Issue, error) {
	return r.mockAssignedIssues, r.mockAssignedIssuesErr
}

type PrinterMock struct {
}

func (p PrinterMock) PrintIssues(issues []*redmine.Issue) []string {
	var messages []string
	for _, issue := range issues {
		messages = append(messages, issue.Subject)
	}
	return messages
}

func (p PrinterMock) Print(issue redmine.Issue, printDescription bool) []string {
	if printDescription {
		return []string{"description"}
	}
	return []string{"empty"}
}
