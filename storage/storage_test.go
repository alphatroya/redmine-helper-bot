package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

type RedisMock struct {
	redis.Cmdable
	mockStorage map[string]string
}

func newRedisMock() *RedisMock {
	return &RedisMock{mockStorage: make(map[string]string)}
}

func (t *RedisMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	t.mockStorage[key] = value.(string)
	return redis.NewStatusCmd(value)
}

func (t *RedisMock) Get(key string) *redis.StringCmd {
	result, ok := t.mockStorage[key]
	if !ok && len(result) > 0 {
		return redis.NewStringResult("", fmt.Errorf("storage value is nil"))
	}
	return redis.NewStringResult(result, nil)
}

func TestHostStorage(t *testing.T) {
	host := "www.google.com"
	var chat int64 = 5
	mock := newRedisMock()
	sut := RedisStorage{mock}
	sut.SetHost(host, chat)
	if mock.mockStorage["5_host"] != host {
		t.Errorf("storing value in redis failed, got: \"%s\"", mock.mockStorage["5_host"])
	}

	restoredHost, err := sut.GetHost(chat)
	if err != nil {
		t.Errorf("getting error during host obtaining, got: %s", err)
	}
	if restoredHost != host {
		t.Errorf("getting value from redis failed, expected: \"%s\", got: \"%s\"", host, restoredHost)
	}
}

func TestTokenStorage(t *testing.T) {
	token := "d3i3j423432"
	var chat int64 = 5
	mock := newRedisMock()
	sut := RedisStorage{mock}
	sut.SetToken(token, chat)
	if mock.mockStorage["5_token"] != token {
		t.Errorf("storing value in redis failed, got: \"%s\"", mock.mockStorage["5_token"])
	}

	restoredToken, err := sut.GetToken(chat)
	if err != nil {
		t.Errorf("getting error during token obtaining, got: %s", err)
	}
	if restoredToken != token {
		t.Errorf("getting value from redis failed, expected: \"%s\", got: \"%s\"", token, restoredToken)
	}
}
