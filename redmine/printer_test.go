package redmine

import (
	"strings"
	"testing"
)

func TestTablePrinter_Print(t *testing.T) {
	testCases := []struct {
		issue Issue
	}{
		{
			issue: Issue{
				AssignedTo: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{},
				Author: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{},
				Description: "",
				Project: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{},
				SpentHours: 0,
				Status: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{},
				Subject: "",
				Tracker: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{},
			},
		},
	}
	for _, testCase := range testCases {
		printer := TablePrinter{}
		table := printer.Print(testCase.issue, true)
		if len(table) != 3 {
			t.Errorf("table should contains 3 elements")
		}
		if !strings.Contains(table[0], testCase.issue.Subject) {
			t.Errorf("first element should contain subject")
		}
		if !strings.Contains(table[2], testCase.issue.Description) {
			t.Errorf("last element should contain description")
		}
	}
}
