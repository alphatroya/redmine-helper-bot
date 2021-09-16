package commands

import "github.com/alphatroya/redmine-helper-bot/redmine"

type Printer interface {
	Print(issue redmine.Issue, printDescription bool) []string
}
