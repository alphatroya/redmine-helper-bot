package commands

type CommandsBuilder interface {
	Build(command string, message string, previousCommand Command)
}

type BotCommandsBuilder struct {
}

func (b BotCommandsBuilder) Build(command string, message string, previousCommand Command) {
}
