package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func HandleTokenMessage(message string, redisClient redis.Cmdable, chatID int64) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return WrongTokenMessageResponse
	}
	redisClient.Set(fmt.Sprint(chatID)+"_token", splittedMessage[1], 0)
	return SuccessTokenMessageResponse
}

func HandleHostMessage(message string, redisClient redis.Cmdable, chatID int64) (string, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return "", errors.New(WrongHostMessageResponse)
	}
	_, err := url.ParseRequestURI(splittedMessage[1])
	if err != nil {
		return "", err
	}
	redisClient.Set(fmt.Sprint(chatID)+"_host", splittedMessage[1], 0)
	return SuccessHostMessageResponse, nil
}
