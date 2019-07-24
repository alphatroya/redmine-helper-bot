package main

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisMock struct {
	redis.Cmdable
	lastValue   string
	storage     map[string]string
	storageType string
}

func NewRedisMock(storageType string) *RedisMock {
	mock := new(RedisMock)
	mock.storage = make(map[string]string)
	mock.storageType = storageType
	return mock
}

func (r *RedisMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	r.storage[key] = value.(string)
	r.lastValue = value.(string)
	return redis.NewStatusCmd(value)
}

func (r *RedisMock) Get(key string) *redis.StringCmd {
	return redis.NewStringResult(r.storage[key+r.storageType], nil)
}
