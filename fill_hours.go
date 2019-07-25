package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
)

func HandleFillMessage(message string, chatID int64, redisClient redis.Cmdable, client HttpClient) (string, error) {
	chatIDString := fmt.Sprint(chatID)

	token, err := redisClient.Get(chatIDString + "_token").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursTokenNilResponse)
	}

	host, err := redisClient.Get(chatIDString + "_host").Result()
	if err != nil {
		return "", fmt.Errorf(WrongFillHoursHostNilResponse)
	}

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf("Неправильное количество аргументов")
	}

	regex := regexp.MustCompile(`[0-9]+`)
	if regex.MatchString(splitted[1]) == false {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIdResponse)
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
