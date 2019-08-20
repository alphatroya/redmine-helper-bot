package commands

import "testing"

func TestIntro_Handle(t *testing.T) {
	sut := Intro{}
	result, err := sut.Handle("")
	expectedResult := `
*Бот активирован*

Перед тем как его использовать, необходимо установить:

1. адрес сервера c помощью команды /host <адрес сервера>
2. ключ доступа к API c помощью команды /token <API токен>

Ключ доступа к API можно найти в боковом меню на экране настроек профиля в Redmine.

Заполнение часов осуществляется с помощью команды /fillhours
`
	if result != nil && result.Message() != expectedResult {
		t.Errorf("intro command should return correct message %s", result.Message())
	}
	if err != nil {
		t.Errorf("intro command should not return error, got: %s", err)
	}
}

func TestIntro_IsCompleted(t *testing.T) {
	sut := Intro{}
	if sut.IsCompleted() != true {
		t.Errorf("intro command should always be completed")
	}
}

func TestNewIntroConstructor(t *testing.T) {
	if newIntroCommand() == nil {
		t.Error("new intro should return a new instance")
	}
}
