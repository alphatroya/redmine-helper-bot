package commands

type Command interface {
	Handle(message string) (string, error)
	Cancel()
}
