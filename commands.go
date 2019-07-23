package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleTokenMessage(message string, tokens map[int64]string, chatID int64) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return "Неправильное количество аргументов"
	}
	redisClient.Set(fmt.Sprint(chatID)+"_token", splittedMessage[1], 0)
	tokens[chatID] = splittedMessage[1]
	return "Токен успешно обновлен"
}

func HandleHostMessage(message string, tokens map[int64]string, chatID int64) (string, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 {
		return "", errors.New("Неправильное количество аргументов")
	}
	_, err := url.ParseRequestURI(splittedMessage[1])
	if err != nil {
		return "", err
	}
	redisClient.Set(fmt.Sprint(chatID)+"_host", splittedMessage[1], 0)
	tokens[chatID] = splittedMessage[1]
	return "Адрес сервера успешно обновлен", nil
}

func HandleFillMessage(message string, update tgbotapi.Update, tokens map[int64]string, hosts map[int64]string) (string, error) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	chatIDString := fmt.Sprint(update.Message.From.ID)

	token, err := redisClient.Get(chatIDString + "_token").Result()
	if err != nil {
		return "", fmt.Errorf("Токен доступа для текущего пользователя не найден")
	}

	host, err := redisClient.Get(chatIDString + "_host").Result()
	if err != nil {
		return "", fmt.Errorf("Адрес сервера не найден")
	}

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf("Неправильное количество аргументов")
	}

	requestBody, err := MakeFillHoursRequest(token, host, splitted)
	if err != nil {
		return "", err
	}
	resultMessage := fmt.Sprintf("В задачу %s добавлено часов: %s", requestBody.TimeEntry.IssueID, requestBody.TimeEntry.Hours)
	return resultMessage, nil
}

func MakeFillHoursRequest(token string, host string, message []string) (*RequestBody, error) {
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
	client := &http.Client{}
	response, _ := client.Do(request)

	defer response.Body.Close()
	return requestBody, nil
}
