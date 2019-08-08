package redmine

import (
	"fmt"
	"testing"
)

func TestSuccessRequest(t *testing.T) {
	networkClient := NewClientRequestMock(200, nil, `{ "test": "test" }`)
	redmine := ClientManager{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err != nil {
		t.Errorf("Success request should return nil error, got: %s", err)
	}
}

func TestNetworkErrorResponse(t *testing.T) {
	networkClient := NewClientRequestMock(400, fmt.Errorf("error"), "")
	redmine := ClientManager{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err != nil && err.Error() != "error" {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}

func TestWrongResponseRequest(t *testing.T) {
	networkClient := NewClientRequestMock(400, nil, "")
	redmine := ClientManager{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err != nil && err.Error() != WrongStatusCodeError(400, "Bad Request").Error() {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}
