package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

type ClientRequestMock struct {
	statusCode int
}

func (c *ClientRequestMock) Do(req *http.Request) (*http.Response, error) {
	response := &http.Response{}
	if c.statusCode != 0 {
		response.StatusCode = c.statusCode
	} else {
		response.StatusCode = 200
	}
	response.Body = &bodyMock{}
	return response, nil
}

type bodyMock struct{}

func (b *bodyMock) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (b *bodyMock) Close() error {
	return nil
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
