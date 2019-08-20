package commands

type StopCommand struct {
}

func newStopCommand() *StopCommand {
	return &StopCommand{}
}

func (s StopCommand) Handle(message string) (*CommandResult, error) {
	return NewCommandResult("Бот остановлен"), nil
}

func (s StopCommand) IsCompleted() bool {
	return true
}
