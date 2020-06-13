package cli

import "fmt"

// Command define a command.
type Command struct {
	Name string
	Key  rune
	Fn   func()
}

// CommandInfoFormatter formats a command.
type CommandInfoFormatter func(Command) string

func helpFormatter(c Command) string {
	return fmt.Sprintf("Key: %q, Command: %s", c.Key, c.Name)
}

func inlineFormatter(c Command) string {
	if c.StartsWithKey() {
		return fmt.Sprintf("(%s)%s", c.Name[0:1], c.Name[1:])
	}
	return fmt.Sprintf("%s:%q", c.Name, c.Key)
}

// Run runs the command and sets the Prompt to indicate that the command was run.
func (c *Command) Run() {
	Prompt("Running command: %s (%q)", c.Name, c.Key)
	c.Fn()
}

// StartsWithKey tells if a command name starts with the command key.
func (c *Command) StartsWithKey() bool {
	return len(c.Name) > 0 && c.Name[0] == byte(c.Key)
}
