package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
)

func TestCommands(t *testing.T) {
	var numCalls int
	fn := func() {
		numCalls++
	}

	quit := cli.QuitCommands(fn)
	cli.SetCommands(quit)
	err := cli.GetCommands().RunScript("qQ\x04\x03")
	assert.NoError(t, err, "script must not fail")

	assert.Equal(t, 4, numCalls, "all quit keys must have triggered a command")

	// trigger sleep code in script runner
	err = cli.GetCommands().RunScript("0")
	assert.NoError(t, err, "sleep codes must not fail")

	err = cli.GetCommands().RunScript("x")
	assert.Error(t, err, "missing commands should fail")
}

func TestCustomCommands(t *testing.T) {
	var numCalls int
	cmds := []cli.Command{
		{Name: "func1", Key: 'a', Fn: func() { numCalls += 10 }},
		{Name: "func2", Key: 'f', Fn: func() { numCalls += 100 }},
	}
	cli.SetCommands(cmds)
	cli.GetCommands().RunScript("af")
	line := cli.GetCommands().String()

	assert.Contains(t, line, "func1:'a'")
	assert.Contains(t, line, "(f)unc2")
	assert.Equal(t, 110, numCalls, "all commands must have run")

	help := cli.GetCommands().Help()
	assert.Contains(t, help, "func1")
	assert.Contains(t, help, "func2")
}
