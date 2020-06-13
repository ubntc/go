package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Commands provides convenient functions on a list of commands.
type Commands []Command

// Get returns a command or nil.
func (c Commands) Get(r rune) *Command {
	for _, cmd := range c {
		if r == cmd.Key {
			return &cmd
		}
	}
	return nil
}

// String returns the command names and keys on a single line.
func (c Commands) String() string {
	return "Commands: " + strings.Join(c.Info(inlineFormatter), " ")
}

// Help returns the command names and keys, one command per line.
func (c Commands) Help() string {
	res := append([]string{"Keyboard Commands:"}, c.Info(helpFormatter)...)
	return strings.Join(res, "\n\r  ")
}

// Info returns the command names and keys as strings.
func (c Commands) Info(fn CommandInfoFormatter) []string {
	var res []string
	cmds := make(map[string]Command)
	for _, cmd := range c {
		if _, ok := cmds[cmd.Name]; !ok || cmd.StartsWithKey() {
			cmds[cmd.Name] = cmd
		}
	}
	for _, cmd := range cmds {
		res = append(res, fn(cmd))
	}
	sort.Strings(res)
	return res
}

// Run runs a command or returns a NotFound error.
func (c Commands) Run(r rune) error {
	if cmd := c.Get(r); cmd != nil {
		cmd.Run()
		return nil
	}
	return fmt.Errorf("Command not found for Key=%q", r)
}

// RunScript runs the given script with the Commands.
func (c Commands) RunScript(s string) error {
	for _, r := range []rune(s) {
		num := int(r - '0')
		switch {
		case num >= 0 && num < 10:
			time.Sleep(time.Duration(num) * time.Second)
		default:
			if err := c.Run(r); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetCommands sets the global commands.
func SetCommands(c Commands) {
	// store a copy of the given commands
	term.mu.Lock()
	defer term.mu.Unlock()
	term.commands = append(Commands{}, c...)
}

// GetCommands returns a copy of the global commands.
func GetCommands() Commands {
	term.mu.RLock()
	defer term.mu.RUnlock()
	return append(Commands{}, term.commands...)
}

// QuitCommands returns quit commands.
func QuitCommands(fn func()) Commands {
	return Commands{
		Command{"quit", 'q', fn},
		Command{"quit", 'Q', fn},
		Command{"quit", 3, fn}, // CTRL-C
		Command{"quit", 4, fn}, // CTRL-D
	}
}
