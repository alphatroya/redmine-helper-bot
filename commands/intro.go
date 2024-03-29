package commands

type IntroCommand struct {
}

func newIntroCommand() *IntroCommand {
	return &IntroCommand{}
}

func (i IntroCommand) Handle(message string) (*CommandResult, error) {
	return NewCommandResult(`
*Бот активирован*

Перед тем как начать полноценно его использовать, его необходимо настроить. Это делается с помощью 2х команд:

1. Установите адрес сервера c помощью команды /host <адрес сервера>
2. Задайте токен доступа к API c помощью команды /token <API токен>. Токен можно найти в боковом меню на экране настроек профиля в Redmine. Он будет сохранен в базе в зашифрованном виде.

- Заполнение часов с помощью одной команды делается с помощью /fh (удобно при использовании отложенных сообщений в Telegram)
- Для того чтобы заполнить оставшиеся за сегодня часы можно воспользоваться командой /fhm
- Изменить активность по-умолчанию можно с помощью команды /activity
`), nil
}

func (i IntroCommand) IsCompleted() bool {
	return true
}
