package redmine

type TimeEntryBody struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	IssueID    string `json:"issue_id"`
	Hours      string `json:"hours"`
	Comments   string `json:"comments"`
	ActivityID string `json:"activity_id,omitempty"`
}

type TimeEntryBodyResponse struct {
	TimeEntry TimeEntryResponse `json:"time_entry"`
}

type TimeEntriesBodyResponse struct {
	TimeEntries []*TimeEntryResponse `json:"time_entries"`
}

type TimeEntryResponseIssue struct {
	ID int `json:"id"`
}

type TimeEntryResponseActivity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TimeEntryResponse struct {
	Activity  TimeEntryResponseActivity `json:"activity"`
	Comments  string                    `json:"comments"`
	CreatedOn string                    `json:"created_on"`
	Hours     float32                   `json:"hours"`
	ID        int                       `json:"id"`
	Issue     TimeEntryResponseIssue    `json:"issue"`
	Project   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	SpentOn   string `json:"spent_on"`
	UpdatedOn string `json:"updated_on"`
	User      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}
