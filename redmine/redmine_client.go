package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Redmine interface {
	SetToken(token string)
	SetHost(host string)
	FillHoursRequest(message []string) (*RequestBody, error)
}

func NewRedmineClient(client HTTPClient) *RedmineClient {
	return &RedmineClient{"", "", client}
}

func WrongRedmineStatusCodeError(statusCode int, statusText string) error {
	return fmt.Errorf("Wrong response from redmine server %d - %s", statusCode, statusText)
}

type RedmineClient struct {
	token         string
	host          string
	networkClient HTTPClient
}

func (r *RedmineClient) SetToken(token string) {
	r.token = token
}

func (r *RedmineClient) SetHost(host string) {
	r.host = host
}

func (t *RedmineClient) FillHoursRequest(message []string) (*RequestBody, error) {
	requestBody := &RequestBody{
		&TimeEntry{
			message[1],
			message[2],
			strings.Join(message[3:], " "),
		},
	}

	json, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", t.host+"/time_entries.json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Redmine-API-Key", t.token)
	request.Header.Set("Content-Type", "application/json")
	response, err := t.networkClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, WrongRedmineStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	return requestBody, nil
}
