package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	SetToken(token string)
	SetHost(host string)
	FillHoursRequest(issueID string, hours string, comment string) (*TimeEntryBodyResponse, error)
	Issue(issueID string) (*Issue, error)
}

func WrongStatusCodeError(statusCode int, statusText string) error {
	return fmt.Errorf("получен ошибочный статус от сервера: %d - %s", statusCode, statusText)
}

type ClientManager struct {
	token         string
	host          string
	networkClient HTTPClient
}

func NewClientManager(networkClient HTTPClient) *ClientManager {
	return &ClientManager{networkClient: networkClient}
}

func (r *ClientManager) SetToken(token string) {
	r.token = token
}

func (r *ClientManager) SetHost(host string) {
	r.host = host
}

func (r *ClientManager) Issue(issueID string) (*Issue, error) {
	request, err := http.NewRequest("GET", r.host+"/issues/"+issueID+".json", nil)
	if err != nil {
		return nil, err
	}

	r.configure(request)
	response, err := r.networkClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, WrongStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	readBytes, err := ioutil.ReadAll(response.Body)

	issue := new(Issue)
	err = json.Unmarshal(readBytes, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (r *ClientManager) FillHoursRequest(issueID string, hours string, comment string) (*TimeEntryBodyResponse, error) {
	requestBody := &TimeEntryBody{
		&TimeEntry{
			issueID,
			hours,
			comment,
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", r.host+"/time_entries.json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.configure(request)
	response, err := r.networkClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, WrongStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	bytesResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := new(TimeEntryBodyResponse)
	err = json.Unmarshal(bytesResponse, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ClientManager) configure(request *http.Request) {
	request.Header.Set("X-Redmine-API-Key", r.token)
	request.Header.Set("Content-Type", "application/json")
}
