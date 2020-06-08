// Package generics shows how to NOT use interfaces to emulate generics!
package generics

import (
	"golang.org/x/sync/errgroup"
)

// Operation defines an operation to run with each item.
type Operation func(index int, item interface{}) error

// ForEach is too slow. This is just a demo. DON'T USE!
//
// ForEach concurrently runs an `Operation` on the given items in a channel and waits for them to finish.
// If an error occurs it returns this error after waiting is done.
//
// If Go had generics we could define a generic list type and would not have to rely on a channel.
//
func ForEach(ch <-chan interface{}, fn Operation) error {
	var g errgroup.Group
	i := 0
	for v := range ch {
		index := i
		item := v
		g.Go(func() error {
			return fn(index, item)
		})
		i++
	}
	return g.Wait()
}
