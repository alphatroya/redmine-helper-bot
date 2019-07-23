package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestHandleTokenMessageWithWrongArgumentsCount(t *testing.T) {
	tables := []struct {
		message string
		failure string
	}{
		{"/token", "Empty command"},
		{"/token test test", "Too many command"},
	}

	for _, message := range tables {
		mock := NewRedisMock("_token")
		if HandleTokenMessage(message.message, mock, 0) != "Неправильное количество аргументов" {
			t.Errorf("Arguments check failed with wrong result %s", message.failure)
		}
	}
}

func TestHandleTokenMessageWithRightArgumentsCount(t *testing.T) {
	tables := []struct {
		token  string
		chatID int64
	}{
		{"431", 44},
		{"23", 45},
	}

	for _, message := range tables {
		mock := NewRedisMock("_token")
		result := HandleTokenMessage("/token"+" "+message.token, mock, message.chatID)
		tokenValue := mock.Get(fmt.Sprint(message.chatID)).Val()
		if tokenValue != message.token {
			t.Errorf("Arguments check failed with wrong result")
			t.Errorf(tokenValue)
		}
		if result != "Токен успешно обновлен" {
			t.Errorf("Wrong response from method")
		}
	}
}

// func TestHandleHostMessageWithWrongArgumentsCount(t *testing.T) {
// 	tables := []struct {
// 		message string
// 		failure string
// 	}{
// 		{"/host", "Empty command"},
// 		{"/host test test", "Too many command"},
// 		{"/host test", "Input is not correct URL"},
// 	}

// 	for _, message := range tables {
// 		_, err := HandleHostMessage(message.message, make(map[int64]string), 0)
// 		if err == nil {
// 			t.Errorf("Method should return error for wrong input %s", message.failure)
// 		}
// 	}
// }

// func TestHandleHostMessageWithRightArgumentsCount(t *testing.T) {
// 	hosts := make(map[int64]string)
// 	tables := []struct {
// 		url    string
// 		chatID int64
// 	}{
// 		{"https://www.google.com", 44},
// 		{"https://www.tt.com", 45},
// 		{"https://tt.com", 46},
// 	}

// 	for _, message := range tables {
// 		result, err := HandleHostMessage("/token"+" "+message.url, hosts, message.chatID)
// 		if err != nil {
// 			t.Errorf("Error result from correct input")
// 		}
// 		if hosts[message.chatID] != message.url {
// 			t.Errorf("Arguments check failed with wrong result")
// 		}
// 		if result != "Адрес сервера успешно обновлен" {
// 			t.Errorf("Wrong response from method")
// 		}
// 	}
// }

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
	r.storage[key+r.storageType] = value.(string)
	r.lastValue = value.(string)
	return redis.NewStatusCmd(value)
}

func (r *RedisMock) Get(key string) *redis.StringCmd {
	return redis.NewStringResult(r.storage[key+r.storageType], nil)
	// return redis.NewStringResult(r.lastValue, nil)
}
