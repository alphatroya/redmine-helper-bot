package redmine

import (
	"regexp"
	"strings"
)

func CheckAndExtractIssueID(source string) (trimmed string, success bool) {
	success = regexp.MustCompile(IssueIDRegex).MatchString(source)
	if success {
		trimmed = strings.TrimLeft(source, "#")
	} else {
		trimmed = source
	}
	return
}
