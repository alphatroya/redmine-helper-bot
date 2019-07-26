package main

import (
	"fmt"
	"testing"
)

func TestHandleFillHoursSuccessCommand(t *testing.T) {
	host := "https://test_host.com"
	tables := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 43212 8 Test", 44, SuccessFillHoursMessageResponse("43212", "8", host)},
		{"/fillhours 51293 8.0 Test", 44, SuccessFillHoursMessageResponse("51293", "8.0", host)},
		{"/fillhours 51293 9.6 Test", 44, SuccessFillHoursMessageResponse("51293", "9.6", host)},
	}

	for _, message := range tables {
		mock := NewRedisMock()
		mock.Set(fmt.Sprint(message.chatID)+"_token", "TestToken", 0)
		mock.Set(fmt.Sprint(message.chatID)+"_host", host, 0)
		text, err := HandleFillMessage(message.message, message.chatID, mock, &ClientRequestMock{})
		if err != nil {
			t.Errorf("Success function should not return err %s %s", err, message.message)
		}
		if text != message.expected {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", text, message.expected)
		}
	}
}

func TestHandleFillHoursNilTokenFailCommand(t *testing.T) {
	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursTokenNilResponse}

	mock := NewRedisMock()
	_, err := HandleFillMessage(input.message, input.chatID, mock, &ClientRequestMock{})
	if err == nil {
		t.Errorf("Wrong command should return non-nil err")
	}
	if input.expected != err.Error() {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", err.Error(), input.expected)
	}
}

func TestHandleFillHoursNilHostFailCommand(t *testing.T) {
	input := struct {
		message  string
		chatID   int64
		expected string
	}{"/fillhours 43212 8 Test", 44, WrongFillHoursHostNilResponse}

	mock := NewRedisMock()
	mock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
	_, err := HandleFillMessage(input.message, input.chatID, mock, &ClientRequestMock{})
	if err == nil {
		t.Errorf("Wrong command should return non-nil err")
	}
	if input.expected != err.Error() {
		t.Errorf("Wrong response from fill hours method got %s, expected %s", err.Error(), input.expected)
	}
}

func TestFillHoursWrongInput(t *testing.T) {
	inputs := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours aaaa 8 Test", 44, WrongFillHoursWrongIssueIDResponse},
		{"/fillhours <51293 8 Test", 44, WrongFillHoursWrongIssueIDResponse},
		{"/fillhours 51293 8a Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293 ff Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293 9,6 Test", 44, WrongFillHoursWrongHoursCountResponse},
		{"/fillhours 51293", 44, WrongFillHoursWrongNumberOfArgumentsResponse},
		{"/fillhours 51293 6", 44, WrongFillHoursWrongNumberOfArgumentsResponse},
	}

	for _, input := range inputs {
		mock := NewRedisMock()
		mock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		mock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		_, err := HandleFillMessage(input.message, input.chatID, mock, &ClientRequestMock{})
		if err == nil {
			t.Errorf("Wrong command should return non-nil err")
		}
		if input.expected != err.Error() {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", err.Error(), input.expected)
		}
	}
}

func TestFillHoursWrongRsponse(t *testing.T) {
	inputs := []struct {
		message  string
		chatID   int64
		expected string
	}{
		{"/fillhours 51293 8 Test", 44, fmt.Sprintf(WrongFillHoursWrongStatusCodeResponse, 400, "Bad Request")},
	}

	for _, input := range inputs {
		mock := NewRedisMock()
		mock.Set(fmt.Sprint(input.chatID)+"_token", "TestToken", 0)
		mock.Set(fmt.Sprint(input.chatID)+"_host", "https://test_host.com", 0)
		_, err := HandleFillMessage(input.message, input.chatID, mock, &ClientRequestMock{400})
		if err == nil {
			t.Errorf("Bad response status code should produce error")
		}
		if input.expected != err.Error() {
			t.Errorf("Wrong response from fill hours method got %s, expected %s", err.Error(), input.expected)
		}
	}
}
