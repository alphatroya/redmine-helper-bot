package commands

import "testing"

func TestIntro_Handle(t *testing.T) {
	sut := IntroCommand{}
	result, err := sut.Handle("")
	expectedResult := `
*Бот активирован*

Перед тем как начать полноценно его использовать, его необходимо настроить. Это делается с помощью 2х команд:

1. Установите адрес сервера c помощью команды /host <адрес сервера>
2. Задайте токен доступа к API c помощью команды /token <API токен>. Токен можно найти в боковом меню на экране настроек профиля в Redmine. Он будет сохранен в базе в зашифрованном виде.

- Заполнение часов с помощью одной команды делается с помощью /fh (удобно при использовании отложенных сообщений в Telegram)
- Для того чтобы заполнить оставшиеся за сегодня часы можно воспользоваться командой /fhm
- Изменить активность по-умолчанию можно с помощью команды /activity
`
	if result != nil && result.Message() != expectedResult {
		t.Errorf("intro command should return correct message %s", result.Message())
	}
	if err != nil {
		t.Errorf("intro command should not return error, got: %s", err)
	}
}

func TestIntro_IsCompleted(t *testing.T) {
	sut := IntroCommand{}
	if sut.IsCompleted() != true {
		t.Errorf("intro command should always be completed")
	}
}

func TestNewIntroConstructor(t *testing.T) {
	if newIntroCommand() == nil {
		t.Error("new intro should return a new instance")
	}
}
