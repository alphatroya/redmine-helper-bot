package commands

type UnknownCommand struct {
	response string
}

func NewUnknownCommand() *UnknownCommand {
	return &UnknownCommand{unknownCommandResponse}
}

func NewUnknownCommandWithMessage(message string) *UnknownCommand {
	return &UnknownCommand{message}
}

const unknownCommandResponse = "Введена неправильная команда"

func (u UnknownCommand) Handle(message string) (*CommandResult, error) {
	return NewCommandResult(u.response), nil
}

func (u UnknownCommand) IsCompleted() bool {
	return true
}
