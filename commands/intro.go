package commands

type Intro struct {
}

func newIntroCommand() *Intro {
	return &Intro{}
}

func (i Intro) Handle(message string) (*CommandResult, error) {
	return NewCommandResult(`Бот активирован, перед использованием необходимо указать адрес сервера и авторизационный токен`), nil
}

func (i Intro) IsCompleted() bool {
	return true
}
