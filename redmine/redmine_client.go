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
	FillHoursRequest(issueID string, hours string, comment string, activityID string) (*TimeEntryBodyResponse, error)
	Issue(issueID string) (*IssueContainer, error)
	AssignedIssues() ([]*Issue, error)
	Activities() ([]*Activities, error)
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

func (r *ClientManager) AssignedIssues() ([]*Issue, error) {
	bytesResponse, err := r.sendMessage(nil, "GET", r.host+"/issues.json?assigned_to_id=me")
	if err != nil {
		return nil, err
	}

	issues := new(IssuesList)
	err = json.Unmarshal(bytesResponse, issues)
	if err != nil {
		return nil, err
	}

	return issues.Issues, nil
}

func (r *ClientManager) Issue(issueID string) (*IssueContainer, error) {
	bytesResponse, err := r.sendMessage(nil, "GET", r.host+"/issues/"+issueID+".json")
	if err != nil {
		return nil, err
	}
	issue := new(IssueContainer)
	err = json.Unmarshal(bytesResponse, issue)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (r *ClientManager) FillHoursRequest(issueID string, hours string, comment string, activityID string) (*TimeEntryBodyResponse, error) {
	requestBody := &TimeEntryBody{
		&TimeEntry{
			issueID,
			hours,
			comment,
			activityID,
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	bytesResponse, err := r.sendMessage(bytes.NewBuffer(body), "POST", r.host+"/time_entries.json")
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

func (r *ClientManager) sendMessage(bodyBuffer *bytes.Buffer, requestMethod string, requestURL string) ([]byte, error) {
	request, err := http.NewRequest(requestMethod, requestURL, bodyBuffer)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Redmine-API-Key", r.token)
	request.Header.Set("Content-Type", "application/json")
	response, err := r.networkClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil, WrongStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}
	return ioutil.ReadAll(response.Body)
}
