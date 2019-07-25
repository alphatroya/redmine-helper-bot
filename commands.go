package main

import (
	"bytes"
	"encoding/json"
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

func HandleFillMessage(message string, chatID int64, redisClient redis.Cmdable, client HttpClient) (string, error) {
	chatIDString := fmt.Sprint(chatID)

	token, err := redisClient.Get(chatIDString + "_token").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}

	host, err := redisClient.Get(chatIDString + "_host").Result()
	if err != nil {
		return "", fmt.Errorf("Адрес сервера не найден")
	}

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf("Неправильное количество аргументов")
	}

	requestBody, err := MakeFillHoursRequest(token, host, splitted, client)
	if err != nil {
		return "", err
	}
	resultMessage := fmt.Sprintf(SuccessFillHoursMessageResponse, requestBody.TimeEntry.IssueID, requestBody.TimeEntry.Hours)
	return resultMessage, nil
}

func MakeFillHoursRequest(token string, host string, message []string, client HttpClient) (*RequestBody, error) {
	requestBody := new(RequestBody)
	requestBody.TimeEntry = new(TimeEntry)
	requestBody.TimeEntry.IssueID = message[1]
	requestBody.TimeEntry.Comments = strings.Join(message[3:], " ")
	requestBody.TimeEntry.Hours = message[2]

	json, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", host+"/time_entries.json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Redmine-API-Key", token)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if response.Body != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return requestBody, nil
}
