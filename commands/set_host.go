package commands

import (
	"net/url"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

type SetHostCommand struct {
	storage storage.Manager
	chatID  int64
}

func (s SetHostCommand) IsCompleted() bool {
	return true
}

func newSetHostCommand(storage storage.Manager, chatID int64) *SetHostCommand {
	return &SetHostCommand{storage: storage, chatID: chatID}
}

func (s SetHostCommand) Handle(message string) (*CommandResult, error) {
	split := strings.Split(message, " ")
	if len(split) > 1 || len(message) == 0 {
		return NewCommandResult("Неправильное количество аргументов, введите адрес в формате `/host <адрес сервера>`(например, `/host https://google.ru`)"), nil
	}
	_, err := url.ParseRequestURI(split[0])
	if err != nil {
		return nil, err
	}
	s.storage.SetHost(split[0], s.chatID)
	return NewCommandResult("Адрес сервера успешно обновлен"), nil
}
