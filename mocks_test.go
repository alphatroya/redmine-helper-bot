package main

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RedisMock struct {
	storageToken map[int64]string
	storageHost  map[int64]string
}

func NewRedisMock() *RedisMock {
	mock := new(RedisMock)
	mock.storageToken = make(map[int64]string)
	mock.storageHost = make(map[int64]string)
	return mock
}

func (r RedisMock) SetToken(token string, chat int64) {
	r.storageToken[chat] = token
}

func (r RedisMock) GetToken(chat int64) (string, error) {
	token, ok := r.storageToken[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return token, nil
}

func (r RedisMock) SetHost(host string, chat int64) {
	r.storageHost[chat] = host
}

func (r RedisMock) GetHost(chat int64) (string, error) {
	host, ok := r.storageHost[chat]
	if !ok {
		return "", fmt.Errorf("storage value is nil")
	}
	return host, nil
}

type RedmineClientMock struct {
	host          string
	token         string
	response      interface{}
	responseError error
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

func (r *RedmineClientMock) Issue(issueID string) (*redmine.Issue, error) {
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
