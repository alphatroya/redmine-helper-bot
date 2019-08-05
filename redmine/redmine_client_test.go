package redmine

import (
	"fmt"
	"testing"
)

func TestSuccessRequest(t *testing.T) {
	networkClient := &ClientRequestMock{}
	redmine := RedmineClient{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err != nil {
		t.Errorf("Success request should return nil error, got: %s", err)
	}
}

func TestNetworkErrorResponse(t *testing.T) {
	networkClient := &ClientRequestMock{400, fmt.Errorf("Error")}
	redmine := RedmineClient{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err.Error() != "Error" {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}

func TestWrongResponseRequest(t *testing.T) {
	networkClient := &ClientRequestMock{400, nil}
	redmine := RedmineClient{"", "", networkClient}
	redmine.SetHost("http://google.com")
	redmine.SetToken("fdjsdfjs")
	_, err := redmine.FillHoursRequest("foo", "bar", "baz")
	if err.Error() != WrongRedmineStatusCodeError(400, "Bad Request").Error() {
		t.Errorf("Wrong error instance after received wrong status code, got: %s", err)
	}
}
