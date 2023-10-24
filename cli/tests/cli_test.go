package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/cli/config"
)

func ensureCancel(t *testing.T, done <-chan struct{}) {
	select {
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "context must be cancelled")
		return
	case <-done:
	}
}

func TestSigWaitTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	s, err := cli.SigWait(ctx)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.Nil(t, s)
}

func TestSigWaitCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	s, err := cli.SigWait(ctx)
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func TestCancel(t *testing.T) {
	ctx, cancel := cli.StartTerm(context.Background(), config.Server())
	res := make(chan bool)
	go func(res chan<- bool) {
		cancel()
		<-ctx.Done()
		res <- true
	}(res)
	assert.True(t, <-res, "context must be cancelled")
}

func TestParentCancel(t *testing.T) {
	pctx, pcancel := context.WithCancel(context.Background())
	ctx, _ := cli.StartTerm(pctx, config.Server())
	res := make(chan bool)
	go func(res chan<- bool) {
		pcancel()
		<-ctx.Done()
		res <- true
	}(res)
	assert.True(t, <-res, "context must be cancelled by parent")
}
