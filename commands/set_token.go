package commands

import (
	"github.com/alphatroya/redmine-helper-bot/storage"
	"strings"
)

const (
	WrongTokenMessageResponse   = "Неправильное количество аргументов"
	SuccessTokenMessageResponse = "Токен успешно обновлен"
)

type SetToken struct {
	storage storage.Manager
	chatID  int64
}

func NewSetTokenCommand(storage storage.Manager, chatID int64) *SetToken {
	return &SetToken{storage: storage, chatID: chatID}
}

func (s SetToken) Handle(message string) string {
	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) > 1 || len(message) == 0 {
		return WrongTokenMessageResponse
	}
	s.storage.SetToken(splittedMessage[0], s.chatID)
	return SuccessTokenMessageResponse
}

func (s SetToken) Cancel() {
}
