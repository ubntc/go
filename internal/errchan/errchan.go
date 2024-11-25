package errchan

import "sync"

type errChan chan error

// Chan provides an error chan for collecting concurrent errors with a few lines of code:
//
//    errch, ch :=  errchan.NewChan(size)
//    var wg sync.WaitGroup
//    wg.Add(2)
//    go func() { defer wg.Done(); ch <- errors.New("error A") }
//    go func() { defer wg.Done(); ch <- errors.New("error B") }
//    wg.Wait()
//    errs := errch.Errors() // closes the chan and returns the channel content as slice
//
// After all errors are written, use the convenience methods to extract them:
//
//     s := errch.String()  // returns errors as string
//     l := errch.Strings() // returns errors as []string
//     e := errch.Collect() // returns errors as []error
//     b := errch.JSON()    // returns an error list as JSON bytes
//
// The fist usage of one of these methods clears the channel and moves the errors to a slice
// for safe reuse. The will then be closed and must not be used again.
type Chan struct {
	errChan
	*errStore
}

// chanGroup base type used by other synchronized error collectors.
type chanGroup struct {
	*Chan
	sync.WaitGroup
}

// ChanGroup is an error channel with builtin sync.WaitGroup and error store.
type ChanGroup struct{ *chanGroup }

// NewChan returns a new Chan and a write error channel.
func NewChan(size int) (*Chan, chan<- error) {
	ch := make(chan error, size)
	return &Chan{
		errChan:  ch,
		errStore: newStore(errChan(ch)),
	}, ch
}

// NewChanGroup returns a new ChanGroup and a write-only error channel.
func NewChanGroup(size int) (*ChanGroup, chan<- error) {
	g := newChanGroup(size)
	return &ChanGroup{g}, g.errChan
}

// newChanGroup returns a new synchronized chanGroup.
func newChanGroup(size int) *chanGroup {
	c, _ := NewChan(size)
	return &chanGroup{c, sync.WaitGroup{}}
}

func (ch errChan) errors() []error {
	close(ch)
	var errs []error
	for err := range ch {
		errs = append(errs, err)
	}
	ch = nil
	return errs
}
