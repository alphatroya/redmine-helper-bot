package commands

type UnknownCommand struct {
}

func newUnknownCommand() *UnknownCommand {
	return &UnknownCommand{}
}

const unknownCommandResponse = "Введена неправильная команда"

func (u UnknownCommand) Handle(message string) (string, error) {
	return unknownCommandResponse, nil
}
