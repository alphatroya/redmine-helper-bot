package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

type Client interface {
	FillHoursRequest(issueID string, hours string, comment string, activityID string) (*TimeEntryBodyResponse, error)
	Issue(issueID string) (*IssueContainer, error)
	AssignedIssues() ([]*Issue, error)
	Activities() ([]*Activities, error)
	TodayTimeEntries() ([]*TimeEntryResponse, error)
	AddComment(issueID string, comment string) error
}

func WrongStatusCodeError(statusCode int, statusText string) error {
	return fmt.Errorf("получен ошибочный статус от сервера: %d - %s", statusCode, statusText)
}

type ClientManager struct {
	networkClient HTTPClient
	storage       storage.Manager
	chatID        int64
}

func (r *ClientManager) TodayTimeEntries() ([]*TimeEntryResponse, error) {
	bytesResponse, err := r.sendMessage(nil, "GET", "/time_entries.json?user_id=me&spent_on=today")
	if err != nil {
		return nil, err
	}

	timeEntriesBody := new(TimeEntriesBodyResponse)
	err = json.Unmarshal(bytesResponse, timeEntriesBody)
	if err != nil {
		return nil, err
	}

	return timeEntriesBody.TimeEntries, nil
}

func NewClientManager(networkClient HTTPClient, storage storage.Manager, chatID int64) *ClientManager {
	return &ClientManager{networkClient: networkClient, storage: storage, chatID: chatID}
}

func (r *ClientManager) AssignedIssues() ([]*Issue, error) {
	bytesResponse, err := r.sendMessage(nil, "GET", "/issues.json?assigned_to_id=me")
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
	bytesResponse, err := r.sendMessage(nil, "GET", "/issues/"+issueID+".json")
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

	bytesResponse, err := r.sendMessage(bytes.NewBuffer(body), "POST", "/time_entries.json")
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

func (r *ClientManager) sendMessage(bodyBuffer io.Reader, requestMethod string, requestURL string) ([]byte, error) {
	host, err := r.storage.GetHost(r.chatID)
	if err != nil {
		return nil, fmt.Errorf("Адрес сервера не задан! Пожалуйста задайте его с помощью команды /host <адрес сервера>")
	}
	request, err := http.NewRequest(requestMethod, host+requestURL, bodyBuffer)
	if err != nil {
		return nil, err
	}
	token, err := r.storage.GetToken(r.chatID)
	if err != nil {
		return nil, fmt.Errorf("Токен доступа к API не задан! Пожалуйста задайте его с помощью команды /token <токен>")
	}
	request.Header.Set("X-Redmine-API-Key", token)
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
