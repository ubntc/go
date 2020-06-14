package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewlineExpr(t *testing.T) {
	v := rePendingNL.FindString("test")
	assert.Equal(t, "test", v)

	v = rePendingNL.FindString("test\n")
	assert.Equal(t, "", v)

	rest := rePendingNL.FindString("test\n\r\rrest")
	assert.Equal(t, "\r\rrest", rest)

	cr := reStartCR.FindString(rest)
	assert.Equal(t, "\r\r", cr)
}
