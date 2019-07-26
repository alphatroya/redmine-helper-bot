package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

func HandleFillMessage(message string, chatID int64, redisClient redis.Cmdable, client HTTPClient) (string, error) {
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
		return "", fmt.Errorf(WrongFillHoursWrongNumberOfArgumentsResponse)
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(splitted[1]) == false {
		return "", fmt.Errorf(WrongFillHoursWrongIssueIDResponse)
	}

	_, conversionError := strconv.ParseFloat(splitted[2], 32)
	if conversionError != nil {
		return "", fmt.Errorf(WrongFillHoursWrongHoursCountResponse)
	}

	requestBody, err := MakeFillHoursRequest(token, host, splitted, client)
	if err != nil {
		return "", err
	}
	resultMessage := SuccessFillHoursMessageResponse(requestBody.TimeEntry.IssueID, requestBody.TimeEntry.Hours, host)
	return resultMessage, nil
}

func MakeFillHoursRequest(token string, host string, message []string, client HTTPClient) (*RequestBody, error) {
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
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Wrong response from redmine server %d - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	if err != nil {
		return nil, err
	}

	return requestBody, nil
}
