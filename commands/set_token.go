package commands

import (
	"fmt"
	"strings"

	"github.com/alphatroya/redmine-helper-bot/storage"
)

const (
	WrongTokenMessageResponse   = "Неправильное количество аргументов"
	SuccessTokenMessageResponse = "Токен успешно обновлен"
)

type SetTokenCommand struct {
	storage storage.Manager
	chatID  int64
}

func newSetTokenCommand(storage storage.Manager, chatID int64) *SetTokenCommand {
	return &SetTokenCommand{storage: storage, chatID: chatID}
}

func (s SetTokenCommand) Handle(message string) (string, error) {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) > 1 || len(message) == 0 {
		return "", fmt.Errorf(WrongTokenMessageResponse)
	}
	s.storage.SetToken(splittedMessage[0], s.chatID)
	return SuccessTokenMessageResponse, nil
}
