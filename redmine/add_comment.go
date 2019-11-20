package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (r *ClientManager) AddComment(issueID string, comment string, assignedTo int) error {
	type Issue struct {
		Notes      string `json:"notes"`
		AssignedTo int    `json:"assigned_to_id,omitempty"`
	}
	type AddCommentPayload struct {
		Issue Issue `json:"issue"`
	}

	body, err := json.Marshal(AddCommentPayload{Issue: Issue{Notes: comment, AssignedTo: assignedTo}})
	if err != nil {
		return err
	}

	_, err = r.sendMessage(bytes.NewBuffer(body), "PUT", fmt.Sprintf("/issues/%s.json", issueID))
	if err != nil {
		return err
	}
	return nil
}
