package commands

import (
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
)

type Activities struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
}

func newActivitiesCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *Activities {
	return &Activities{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (a Activities) Handle(message string) (*CommandResult, error) {
	activities, err := a.redmineClient.Activities()
	if err != nil {
		return nil, err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"ID", "Название"})

	for _, activity := range activities {
		data := []string{
			fmt.Sprintf("%d", activity.Id),
			wrap(activity.Name, 20),
		}
		table.Append(data)
	}
	table.Render()
	return NewCommandResult("`" + tableString.String() + "`"), nil
}

func (a Activities) IsCompleted() bool {
	return true
}
