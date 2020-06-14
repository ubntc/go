package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
)

func TestRawMode(t *testing.T) {
	term := cli.GetTerm()
	term.SetVerbose(true).SetDebug(true)
	restore, err := cli.ClaimTerminal()
	assert.Error(t, err)
	assert.Nil(t, restore)
	assert.Equal(t, false, cli.GetTerm().IsRaw())
	term.SetVerbose(false).SetDebug(false)
}
