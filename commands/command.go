package commands

type Command interface {
	Handle(message string) (*CommandResult, error)
	IsCompleted() bool
}

type CommandResult struct {
	message string
}

func (c *CommandResult) Message() string {
	return c.message
}

func NewCommandResult(message string) *CommandResult {
	return &CommandResult{message: message}
}
