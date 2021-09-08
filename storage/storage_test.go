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

func (t *RedisMock) Del(keys ...string) *redis.IntCmd {
	for _, key := range keys {
		t.mockStorage[key] = ""
	}
	return redis.NewIntCmd()
}

func (t *RedisMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if valueString, ok := value.(string); ok {
		t.mockStorage[key] = valueString
	} else if bytes, ok := value.([]byte); ok {
		t.mockStorage[key] = string(bytes)
	}
	return redis.NewStatusCmd(value)
}

func (t *RedisMock) Get(key string) *redis.StringCmd {
	result, ok := t.mockStorage[key]
	if !ok && len(result) > 0 {
		return redis.NewStringResult("", fmt.Errorf("storage value is nil"))
	}
	return redis.NewStringResult(result, nil)
}

func TestActivityStorage(t *testing.T) {
	activity := "6"
	var chat int64 = 5
	mock := newRedisMock()
	sut := RedisStorage{mock, "123"}
	sut.SetActivity(activity, chat)
	if mock.mockStorage["5_activity"] != activity {
		t.Errorf("storing value in redis failed, got: \"%s\"", mock.mockStorage["5_activity"])
	}

	restoredHost, err := sut.GetActivity(chat)
	if err != nil {
		t.Errorf("getting error during host obtaining, got: %s", err)
	}
	if restoredHost != activity {
		t.Errorf("getting value from redis failed, expected: \"%s\", got: \"%s\"", activity, restoredHost)
	}
}

func TestHostStorage(t *testing.T) {
	host := "www.google.com"
	var chat int64 = 5
	mock := newRedisMock()
	sut := RedisStorage{mock, "123"}
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

const token = "d3i3j423432"

func TestTokenStorage(t *testing.T) {
	var chat int64 = 5
	mock := newRedisMock()
	passphrase := "123"
	sut := RedisStorage{mock, passphrase}
	sut.SetToken(token, chat)
	restoredToken, err := sut.GetToken(chat)
	if mock.mockStorage["5_encrypted"] == "" {
		t.Error("storing value in redis failed, encrypted field is nil")
	}

	if err != nil {
		t.Errorf("getting error during token obtaining, got: %s", err)
	}
	if restoredToken != token {
		t.Errorf("getting value from redis failed, expected: \"%s\", got: \"%s\"", token, restoredToken)
	}
}

func TestMigrationStorage(t *testing.T) {
	var chat int64 = 5
	mock := newRedisMock()
	passphrase := "123"
	sut := RedisStorage{mock, passphrase}
	sut.SetToken(token, chat)
	restoredToken, err := sut.GetToken(chat)
	if err != nil {
		t.Errorf("getting error during token obtaining, got: %s", err)
	}
	if restoredToken != token {
		t.Errorf("getting value from redis failed, expected: \"%s\", got: \"%s\"", token, restoredToken)
	}
}

func TestRedisStorage_ResetData(t *testing.T) {
	var chat int64 = 5
	mock := newRedisMock()
	sut := RedisStorage{mock, "123"}
	sut.SetToken(token, chat)
	sut.SetHost("https://google.com", chat)
	sut.SetActivity("555", chat)

	if err := sut.ResetData(chat); err != nil {
		t.Errorf("reset data should no reset data, got err: %s", err)
	}

	if mock.mockStorage["5_encrypted"] != "" || mock.mockStorage["5_host"] != "" || mock.mockStorage["5_activity"] != "" {
		t.Errorf("storage data is not nil after resetting")
	}
}
