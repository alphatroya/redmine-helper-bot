package commands

type Command interface {
	Handle(message string) (*CommandResult, error)
	IsCompleted() bool
}

type CommandResult struct {
	messages []string
	buttons  []string
}

func (c *CommandResult) Buttons() []string {
	return c.buttons
}

func (c *CommandResult) Message() string {
	return c.messages[0]
}

func (c *CommandResult) Messages() []string {
	return c.messages
}

func NewCommandResult(message string) *CommandResult {
	return &CommandResult{messages: []string{message}}
}

func NewCommandResultWithMessages(messages []string) *CommandResult {
	return &CommandResult{messages: messages}
}

func NewCommandResultWithKeyboard(message string, buttons []string) *CommandResult {
	return &CommandResult{messages: []string{message}, buttons: buttons}
}
