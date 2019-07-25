package main

import (
	"fmt"
	"net/http"
	"time"

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

type ClientRequestMock struct {
}

func (c *ClientRequestMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}
