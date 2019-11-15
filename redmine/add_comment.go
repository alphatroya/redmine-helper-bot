package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (r *ClientManager) AddComment(issueID string, comment string) error {
	type Notes struct {
		Notes string `json:"notes"`
	}
	type AddCommentPayload struct {
		Issue Notes `json:"issue"`
	}

	body, err := json.Marshal(AddCommentPayload{Issue: Notes{Notes: comment}})
	if err != nil {
		return err
	}

	_, err = r.sendMessage(bytes.NewBuffer(body), "PUT", fmt.Sprintf("/issues/%s.json", issueID))
	if err != nil {
		return err
	}
	return nil
}
