package redmine

import (
	"regexp"
	"strings"
)

// CheckAndExtractIssueID check input source string and return check result and id string with trimmed # left symbol
func CheckAndExtractIssueID(source string) (trimmed string, success bool) {
	success = regexp.MustCompile(issueIDRegex).MatchString(source)
	if success {
		trimmed = strings.TrimLeft(source, "#")
	}
	return
}
