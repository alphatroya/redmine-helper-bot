package commands

import (
	"fmt"
	"github.com/alphatroya/redmine-helper-bot/storage"
	"net/url"
	"strings"
)

const (
	WrongHostMessageResponse   = "Неправильное количество аргументов"
	SuccessHostMessageResponse = "Адрес сервера успешно обновлен"
)

type SetHostCommand struct {
	storage storage.Manager
	chatID  int64
}

func NewSetHostCommand(storage storage.Manager, chatID int64) *SetHostCommand {
	return &SetHostCommand{storage: storage, chatID: chatID}
}

func (s SetHostCommand) Handle(message string) (string, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) > 1 || len(message) == 0 {
		return "", fmt.Errorf(WrongHostMessageResponse)
	}
	_, err := url.ParseRequestURI(splittedMessage[0])
	if err != nil {
		return "", err
	}
	s.storage.SetHost(splittedMessage[0], s.chatID)
	return SuccessHostMessageResponse, nil
}

func (s SetHostCommand) Cancel() {
}
