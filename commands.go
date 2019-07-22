package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	tokens[chatID] = splittedMessage[1]
	return "Адрес сервера успешно обновлен", nil
}

func HandleFillMessage(message string, update tgbotapi.Update, bot *tgbotapi.BotAPI, tokens map[int64]string, hosts map[int64]string) (string, error) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	token, ok := tokens[update.Message.Chat.ID]
	if !ok {
		return "", fmt.Errorf("Токен доступа для текущего пользователя не найден")
	}

	host, ok := hosts[update.Message.Chat.ID]
	if !ok {
		return "", fmt.Errorf("Адрес сервера не найден")
	}

	splitted := strings.Split(message, " ")
	if len(splitted) < 4 {
		return "", fmt.Errorf("Неправильное количество аргументов")
	}

	err := MakeFillHoursRequest(token, host, splitted)
	if err != nil {
		return "", err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Test updated")
	bot.Send(msg)
	return "Success", nil
}

func MakeFillHoursRequest(token string, host string, message []string) error {
	requestBody := new(RequestBody)
	requestBody.TimeEntry = new(TimeEntry)
	requestBody.TimeEntry.IssueID = message[1]
	requestBody.TimeEntry.Comments = strings.Join(message[3:], " ")
	requestBody.TimeEntry.Hours = message[2]

	json, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	log.Printf("data %s", json)

	request, err := http.NewRequest("POST", host+"/time_entries.json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	request.Header.Set("X-Redmine-API-Key", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, _ := client.Do(request)

	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
