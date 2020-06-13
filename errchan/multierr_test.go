package errchan

import (
	"errors"
	"math/rand"
	"sync"
	"testing"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

const (
	baseErrorRate = 1000 / 1
	small         = 3
	big           = small * 100
)

func numWorkers(factor int) int {
	return rand.Intn(factor) + factor
}

// multiplies the base error rate with rand value 1-10
func errorRate() int {
	return baseErrorRate*int(rand.Int31n(10)) + 1
}

func benchmarkMultiError(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		var err error
		var wg sync.WaitGroup
		wg.Add(size)
		var mu sync.Mutex
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				if i%errRate == 0 {
					mu.Lock()
					err = multierror.Append(err, errors.New("test"))
					mu.Unlock()
				}
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkMultiErr(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		var err error
		var wg sync.WaitGroup
		wg.Add(size)
		var mu sync.Mutex
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				if i%errRate == 0 {
					mu.Lock()
					err = multierr.Append(err, errors.New("test"))
					mu.Unlock()
				}
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkErrorGroup(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		var grp = errgroup.Group{}
		for i := 0; i < size; i++ {
			i := i
			grp.Go(func() error {
				if i%errRate == 0 {
					return errors.New("test")
				}
				return nil
			})
		}
		grp.Wait()
	}
}

func benchmarkErrChan(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		_, ch := NewChan(size)
		var wg sync.WaitGroup
		wg.Add(size)
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				if i%errRate == 0 {
					ch <- errors.New("test")
				}
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkErrList(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		el := NewList()
		el.Add(size)
		for i := 0; i < size; i++ {
			go func(i int) {
				defer el.Done()
				if i%errRate == 0 {
					el.Append(errors.New("test"))
				}
			}(i)
		}
		el.Wait()
	}
}

func benchmarkErrCollect(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		ec, ch := NewGroup()
		ec.Add(size)
		for i := 0; i < size; i++ {
			go func(i int) {
				defer ec.Done()
				if i%errRate == 0 {
					ch <- errors.New("test")
				}
			}(i)
		}
		ec.Wait()
	}
}

func benchmarkErrPool(b *testing.B, factor int) {
	for n := 0; n < b.N; n++ {
		size := numWorkers(factor)
		errRate := errorRate()
		ec := NewPool()
		ec.Add(size)
		for i := 0; i < size; i++ {
			go func(i int) {
				defer ec.Done()
				if i%errRate == 0 {
					ec.Put(errors.New("test"))
				}
			}(i)
		}
		ec.Wait()
	}
}

func BenchmarkMultiErrorS(b *testing.B) { benchmarkMultiError(b, small) }
func BenchmarkMultiErrorM(b *testing.B) { benchmarkMultiError(b, big) }

func BenchmarkMultiErrA(b *testing.B) { benchmarkMultiErr(b, small) }
func BenchmarkMultiErrB(b *testing.B) { benchmarkMultiErr(b, big) }

func BenchmarkErrorGroupA(b *testing.B) { benchmarkErrorGroup(b, small) }
func BenchmarkErrorGroupB(b *testing.B) { benchmarkErrorGroup(b, big) }

func BenchmarkErrChanA(b *testing.B) { benchmarkErrChan(b, small) }
func BenchmarkErrChanB(b *testing.B) { benchmarkErrChan(b, big) }

func BenchmarkErrCollectA(b *testing.B) { benchmarkErrCollect(b, small) }
func BenchmarkErrCollectB(b *testing.B) { benchmarkErrCollect(b, big) }

func BenchmarkErrListA(b *testing.B) { benchmarkErrList(b, small) }
func BenchmarkErrListB(b *testing.B) { benchmarkErrList(b, big) }

func BenchmarkErrPoolA(b *testing.B) { benchmarkErrPool(b, small) }
func BenchmarkErrPoolB(b *testing.B) { benchmarkErrPool(b, big) }
