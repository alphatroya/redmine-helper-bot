package commands

type Stop struct {
}

func newStopCommand() *Stop {
	return &Stop{}
}

func (s Stop) Handle(message string) (*CommandResult, error) {
	return NewCommandResult("Бот остановлен"), nil
}

func (s Stop) IsCompleted() bool {
	return true
}
