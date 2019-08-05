package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Redmine interface {
	SetToken(token string)
	SetHost(host string)
	FillHoursRequest(issueID string, hours string, comment string) (*TimeEntryBody, error)
	Issue(issueID string) (*Issue, error)
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

func (t *RedmineClient) Issue(issueID string) (*Issue, error) {
	request, err := http.NewRequest("GET", t.host+"/issues/"+issueID+".json", nil)
	if err != nil {
		return nil, err
	}

	t.configure(request)
	response, err := t.networkClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, WrongRedmineStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	bytes, err := ioutil.ReadAll(response.Body)

	issue := new(Issue)
	err = json.Unmarshal(bytes, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (t *RedmineClient) FillHoursRequest(issueID string, hours string, comment string) (*TimeEntryBody, error) {
	requestBody := &TimeEntryBody{
		&TimeEntry{
			issueID,
			hours,
			comment,
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

	t.configure(request)
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

func (t *RedmineClient) configure(request *http.Request) {
	request.Header.Set("X-Redmine-API-Key", t.token)
	request.Header.Set("Content-Type", "application/json")
}
