package redmine

import (
	"fmt"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

func TestSuccessRequest(t *testing.T) {
	storageMock := storage.NewStorageMock()
	storageMock.SetHost("http://google.com", 10)
	storageMock.SetToken("fdjsdfjs", 10)
	networkClient := NewClientRequestMock(200, nil, `{ "test": "test" }`)
	redmine := NewClientManager(networkClient, storageMock, 10)
	_, err := redmine.FillHoursRequest("foo", "bar", "baz", "")
	if err != nil {
		t.Errorf("Success request should return nil error, got: %s", err)
	}
}

func TestNetworkErrorResponse(t *testing.T) {
	networkClient := NewClientRequestMock(400, fmt.Errorf("error"), "")
	storageMock := storage.NewStorageMock()
	storageMock.SetHost("http://google.com", 10)
	storageMock.SetToken("fdjsdfjs", 10)
	redmine := NewClientManager(networkClient, storageMock, 10)
	_, err := redmine.FillHoursRequest("foo", "bar", "baz", "")
	if err != nil && err.Error() != "error" {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}

func TestWrongResponseRequest(t *testing.T) {
	networkClient := NewClientRequestMock(400, nil, "")
	storageMock := storage.NewStorageMock()
	storageMock.SetHost("http://google.com", 10)
	storageMock.SetToken("fdjsdfjs", 10)
	redmine := NewClientManager(networkClient, storageMock, 10)
	_, err := redmine.FillHoursRequest("foo", "bar", "baz", "")
	if err != nil && err.Error() != WrongStatusCodeError(400, "Bad Request").Error() {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}
