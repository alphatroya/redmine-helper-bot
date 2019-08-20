package commands

type IntroCommand struct {
}

func newIntroCommand() *IntroCommand {
	return &IntroCommand{}
}

func (i IntroCommand) Handle(message string) (*CommandResult, error) {
	return NewCommandResult(`
*Бот активирован*

Перед тем как его использовать, необходимо установить:

1. адрес сервера c помощью команды /host <адрес сервера>
2. ключ доступа к API c помощью команды /token <API токен>

Ключ доступа к API можно найти в боковом меню на экране настроек профиля в Redmine.

Заполнение часов осуществляется с помощью команды /fillhours
`), nil
}

func (i IntroCommand) IsCompleted() bool {
	return true
}
