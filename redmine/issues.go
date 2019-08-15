package redmine

type IssuesList struct {
	Issues     []*Issue `json:"issues"`
	TotalCount int      `json:"total_count"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
}

type Issue struct {
	AssignedTo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"assigned_to"`
	Author struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	CreatedOn   string `json:"created_on"`
	Description string `json:"description"`
	DoneRatio   int    `json:"done_ratio"`
	DueDate     string `json:"due_date"`
	ID          int    `json:"id"`
	Priority    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"priority"`
	Project struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	SpentHours float32 `json:"spent_hours"`
	Status     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"status"`
	Subject string `json:"subject"`
	Tracker struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tracker"`
	UpdatedOn string `json:"updated_on"`
}

type IssueContainer struct {
	Issue *Issue `json:"issue"`
}
