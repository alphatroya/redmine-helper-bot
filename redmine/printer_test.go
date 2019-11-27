package redmine

import (
	"strings"
	"testing"
)

func TestTablePrinter_PrintIssues(t *testing.T) {
	testCases := []struct {
		issues []*Issue
		count  int
	}{
		{
			issues: []*Issue{
				{
					ID: 1,
					Project: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{
						Name: "Project1",
					},
					Subject: "Foo",
				},
				{
					ID: 3,
					Project: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{
						Name: "Project2",
					},
					Subject: "Bar",
				},
				{
					ID: 2,
					Project: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{
						Name: "Project2",
					},
					Subject: "Bar",
				},
			},
			count: 2,
		},
	}
	for _, testCase := range testCases {
		printer := TablePrinter{}
		table := printer.PrintIssues(testCase.issues)
		if len(table) != testCase.count {
			t.Errorf("result message count should be equal to number of projects, got: %d, expected: %d", len(table), testCase.count)
		}
	}
}

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
