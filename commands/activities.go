package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"github.com/olekukonko/tablewriter"
)

type Activities struct {
	redmineClient redmine.Client
	storage       storage.Manager
	chatID        int64
	activities    []*redmine.Activities
	completed     bool
}

func newActivitiesCommand(redmineClient redmine.Client, storage storage.Manager, chatID int64) *Activities {
	return &Activities{redmineClient: redmineClient, storage: storage, chatID: chatID}
}

func (a *Activities) Handle(message string) (*CommandResult, error) {
	if len(a.activities) == 0 {
		return a.firstPhase()
	}
	activityID := string(regexp.MustCompile(`^[0-9]+`).Find([]byte(message)))
	errorMessage := "Вы ввели неправильный номер активности"
	if len(activityID) == 0 {
		return nil, errors.New(errorMessage)
	}

	for _, activity := range a.activities {
		if fmt.Sprintf("%d", activity.ID) == activityID {
			a.storage.SetActivity(activityID, a.chatID)
			a.completed = true
			return NewCommandResult("_Активность по умолчанию успешно сохранена. Она будет использоваться при каждой команде заполнения часов_"), nil
		}
	}
	return nil, errors.New(errorMessage)
}

func (a *Activities) firstPhase() (*CommandResult, error) {
	activities, err := a.redmineClient.Activities()
	if err != nil {
		return nil, err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"ID", "Название"})

	var buttons []string
	for _, activity := range activities {
		activityID := fmt.Sprintf("%d", activity.ID)
		data := []string{
			activityID,
			wrap(activity.Name, 20),
		}
		buttons = append(buttons, activityID+" - "+activity.Name)
		table.Append(data)
	}
	table.Render()
	a.activities = activities
	return NewCommandResultWithMessagesAndKeyboard([]string{
		"`" + tableString.String() + "`",
		"_Введите номер новой активности по умолчанию_",
	}, buttons), nil
}

func (a *Activities) IsCompleted() bool {
	return a.completed
}
