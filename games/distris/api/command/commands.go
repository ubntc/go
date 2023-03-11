package command

type Command string

const (
	EMPTY = Command("")
	PLAY  = Command("p")
	PAUSE = Command("x")
	GET   = Command("g")
)
