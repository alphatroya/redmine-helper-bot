package commands

type Command interface {
	Handle(message string) (*CommandResult, error)
	IsCompleted() bool
}

type CommandResult struct {
	message string
	buttons []string
}

func (c *CommandResult) Buttons() []string {
	return c.buttons
}

func (c *CommandResult) Message() string {
	return c.message
}

func NewCommandResult(message string) *CommandResult {
	return &CommandResult{message: message}
}

func NewCommandResultWithKeyboard(message string, buttons []string) *CommandResult {
	return &CommandResult{message: message, buttons: buttons}
}
