package commands

import (
	"strings"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

type SetTokenCommand struct {
	storage storage.Manager
	chatID  int64
}

func (s SetTokenCommand) IsCompleted() bool {
	return true
}

func newSetTokenCommand(storage storage.Manager, chatID int64) *SetTokenCommand {
	return &SetTokenCommand{storage: storage, chatID: chatID}
}

func (s SetTokenCommand) Handle(message string) (*CommandResult, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) > 1 || len(message) == 0 {
		return NewCommandResult("Неправильное количество аргументов, введите токен доступа к АПИ в формате `/token <токен>`"), nil
	}
	s.storage.SetToken(splittedMessage[0], s.chatID)
	return NewCommandResult("Токен успешно обновлен"), nil
}
