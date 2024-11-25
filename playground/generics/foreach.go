// Package generics shows how to NOT use interfaces to emulate generics!
package generics

import (
	"golang.org/x/sync/errgroup"
)

// OperationAny defines an operation to run with each item.
type OperationAny func(index int, item any) error

// DEPRECATED: ForEachChan is deprecated. We now have generics in Go.
//
// ForEach concurrently runs an `Operation` on the given items in a channel and waits for them to finish.
// If an error occurs it returns this error after waiting is done.
//
// If Go had generics we could define a generic list type and would not have to rely on a channel.
func ForEachChan(ch <-chan any, fn OperationAny) error {
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

// Operation defines an operation to run with each item.
type Operation[T any] func(index int, item T) error

// ForEach concurrently runs an `Operation` on the given items and waits for them to finish.
func ForEach[T any](values []T, fn Operation[T]) error {
	var g errgroup.Group
	for i, v := range values {
		index := i
		item := v
		g.Go(func() error {
			return fn(index, item)
		})
	}
	return g.Wait()
}
