package commands

import "github.com/alphatroya/redmine-helper-bot/storage"

type StopCommand struct {
	storage storage.Manager
	chatID  int64
}

func newStopCommand(storage storage.Manager, chatID int64) *StopCommand {
	return &StopCommand{storage: storage, chatID: chatID}
}

func (s StopCommand) Handle(message string) (*CommandResult, error) {
	err := s.storage.ResetData(s.chatID)
	if err != nil {
		return nil, err
	}
	return NewCommandResult("Бот остановлен, сохраненные данные очищены"), nil
}

func (s StopCommand) IsCompleted() bool {
	return true
}
