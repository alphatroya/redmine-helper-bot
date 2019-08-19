package commands

type UnknownCommand struct {
}

func NewUnknownCommand() *UnknownCommand {
	return &UnknownCommand{}
}

const unknownCommandResponse = "Введена неправильная команда"

func (u UnknownCommand) Handle(message string) (*CommandResult, error) {
	return NewCommandResult(unknownCommandResponse), nil
}
