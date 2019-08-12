package commands

type Command interface {
	Handle(message string) string
	Cancel()
}
