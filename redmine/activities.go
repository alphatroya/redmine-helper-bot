package redmine

import (
	"encoding/json"
)

type ActivitiesRoot struct {
	TimeEntryActivities []*Activities `json:"time_entry_activities"`
}

type Activities struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (r *ClientManager) Activities() ([]*Activities, error) {
	bytesResponse, err := r.sendMessage(nil, "GET", "/enumerations/time_entry_activities.json")
	if err != nil {
		return nil, err
	}

	activities := new(ActivitiesRoot)
	err = json.Unmarshal(bytesResponse, activities)
	if err != nil {
		return nil, err
	}
	return activities.TimeEntryActivities, nil
}
