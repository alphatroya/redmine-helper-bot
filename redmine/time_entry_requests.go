package redmine

type RequestBody struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	IssueID  string `json:"issue_id"`
	Hours    string `json:"hours"`
	Comments string `json:"comments"`
}
