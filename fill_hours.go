package main

import (
	"fmt"
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

	requestBody, err := FillHoursRequest(token, host, splitted, client)
	if err != nil {
		return "", err
	}
	resultMessage := SuccessFillHoursMessageResponse(requestBody.TimeEntry.IssueID, requestBody.TimeEntry.Hours, host)
	return resultMessage, nil
}
