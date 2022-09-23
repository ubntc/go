package errchan

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestErrChanString(t *testing.T) {
	errch, ch := NewChan(10)
	ch <- errors.New("test")
	ch <- errors.New("test")
	assert.Equal(t, "test\ntest", errch.String())
	assert.Equal(t, "test\ntest", errch.String())
}

func TestErrChanList(t *testing.T) {
	errch, ch := NewChanGroup(10)
	errch.Add(1)
	go func() {
		defer errch.Done()
		ch <- errors.New("test")
		ch <- errors.New("test")
	}()
	errch.Wait()
	errs := []error{
		errors.New("test"),
		errors.New("test"),
	}
	assert.Equal(t, errs, errch.Errors())
	assert.Equal(t, errs, errch.Errors())
}

func TestEmptyErrChan(t *testing.T) {
	errch, _ := NewChanGroup(10)
	errch.Wait()
	assert.Len(t, errch.Errors(), 0)
	assert.NotPanics(t, func() { errch.Errors() })
}

func TestErrChanJSON(t *testing.T) {
	errch, ch := NewChan(10)
	ch <- errors.New("test")
	ch <- errors.New("test")
	assert.Equal(t, "[\"test\",\"test\"]", string(errch.JSON()))
	assert.Equal(t, "[\"test\",\"test\"]", string(errch.JSON()))
}
