package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RequestBody struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	IssueID  string `json:"issue_id"`
	Hours    string `json:"hours"`
	Comments string `json:"comments"`
}

func FillHoursRequest(token string, host string, message []string, client HTTPClient) (*RequestBody, error) {
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
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Wrong response from redmine server %d - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	return requestBody, nil
}
