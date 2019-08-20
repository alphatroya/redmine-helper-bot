package redmine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ActivitiesRoot struct {
	TimeEntryActivities []*Activities `json:"time_entry_activities"`
}

type Activities struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (r *ClientManager) Activities() ([]*Activities, error) {
	request, err := http.NewRequest("GET", r.host+"/enumerations/time_entry_activities.json", nil)
	if err != nil {
		return nil, err
	}

	r.configure(request)
	response, err := r.networkClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, WrongStatusCodeError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	readBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	activities := new(ActivitiesRoot)
	err = json.Unmarshal(readBytes, activities)
	if err != nil {
		return nil, err
	}
	return activities.TimeEntryActivities, nil
}
