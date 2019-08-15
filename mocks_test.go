package main

import (
	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RedmineClientMock struct {
	host          string
	token         string
	response      interface{}
	responseError error
}

func (r *RedmineClientMock) AssignedIssues() ([]*redmine.Issue, error) {
	panic("implement me")
}

func (r *RedmineClientMock) FillHoursRequest(issueID string, hours string, comment string) (*redmine.TimeEntryBodyResponse, error) {
	return r.response.(*redmine.TimeEntryBodyResponse), r.responseError
}

func (r *RedmineClientMock) SetToken(token string) {
	r.token = token
}

func (r *RedmineClientMock) SetHost(host string) {
	r.host = host
}

func (r *RedmineClientMock) SetFillHoursResponse(body *redmine.TimeEntryBodyResponse, responseError error) {
	r.response, r.responseError = body, responseError
}

func (r *RedmineClientMock) Issue(issueID string) (*redmine.IssueContainer, error) {
	return nil, r.responseError
}

type MockBotSender struct {
	text string
}

func (t *MockBotSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if config, ok := c.(tgbotapi.MessageConfig); ok == true {
		t.text = config.Text
	}
	return tgbotapi.Message{}, nil
}
