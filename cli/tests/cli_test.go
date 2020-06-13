package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/cli/cli"
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

	go ensureCancel(t, ctx.Done())

	cli.SigWait(ctx, cancel)
}

func TestSigWaitCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go ensureCancel(t, ctx.Done())
	go cancel()

	cli.SigWait(ctx, cancel)
}

func TestCancel(t *testing.T) {
	ctx, cancel := cli.WithSigWait(context.Background())
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
	ctx, _ := cli.WithSigWait(pctx)
	res := make(chan bool)
	go func(res chan<- bool) {
		pcancel()
		<-ctx.Done()
		res <- true
	}(res)
	assert.True(t, <-res, "context must be cancelled by parent")
}
