package main

import (
	"fmt"
	"time"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	tgbotapi "github.com/alphatroya/telegram-bot-api"
	"github.com/go-redis/redis"
)

type RedisMock struct {
	redis.Cmdable
	storage map[string]string
}

func NewRedisMock() *RedisMock {
	mock := new(RedisMock)
	mock.storage = make(map[string]string)
	return mock
}

func (r *RedisMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	r.storage[key] = value.(string)
	return redis.NewStatusCmd(value)
}

func (r *RedisMock) Get(key string) *redis.StringCmd {
	result, ok := r.storage[key]
	if !ok {
		return redis.NewStringResult("", fmt.Errorf("Storage value is nil"))
	}
	return redis.NewStringResult(result, nil)
}

type RedmineClientMock struct {
	host          string
	token         string
	response      interface{}
	responseError error
}

func (r *RedmineClientMock) SetToken(token string) {
	r.token = token
}

func (r *RedmineClientMock) SetHost(host string) {
	r.host = host
}

func (r *RedmineClientMock) SetFillHoursResponse(body *redmine.RequestBody, responseError error) {
	r.response, r.responseError = body, responseError
}

func (r *RedmineClientMock) Issue(issueID string) (*redmine.Issue, error) {
	return nil, r.responseError
}

func (r *RedmineClientMock) FillHoursRequest(issueID string, hours string, comment string) (*redmine.RequestBody, error) {
	return r.response.(*redmine.RequestBody), r.responseError
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
