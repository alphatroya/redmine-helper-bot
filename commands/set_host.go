package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

const (
	wrongHostMessageResponse   = "Неправильное количество аргументов"
	successHostMessageResponse = "Адрес сервера успешно обновлен"
)

type SetHostCommand struct {
	storage storage.Manager
	chatID  int64
}

func newSetHostCommand(storage storage.Manager, chatID int64) *SetHostCommand {
	return &SetHostCommand{storage: storage, chatID: chatID}
}

func (s SetHostCommand) Handle(message string) (*CommandResult, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) > 1 || len(message) == 0 {
		return nil, fmt.Errorf(wrongHostMessageResponse)
	}
	_, err := url.ParseRequestURI(splittedMessage[0])
	if err != nil {
		return nil, err
	}
	s.storage.SetHost(splittedMessage[0], s.chatID)
	return NewCommandResult(successHostMessageResponse), nil
}
