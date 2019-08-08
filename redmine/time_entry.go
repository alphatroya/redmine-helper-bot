package redmine

type TimeEntryBody struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	IssueID  string `json:"issue_id"`
	Hours    string `json:"hours"`
	Comments string `json:"comments"`
}

type TimeEntryBodyResponse struct {
	TimeEntry TimeEntryResponse `json:"time_entry"`
}

type TimeEntryResponse struct {
	Activity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"activity"`
	Comments  string  `json:"comments"`
	CreatedOn string  `json:"created_on"`
	Hours     float32 `json:"hours"`
	ID        int     `json:"id"`
	Issue     struct {
		ID int `json:"id"`
	} `json:"issue"`
	Project struct {
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

