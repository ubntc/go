package generics

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEach(t *testing.T) {
	ch := make(chan interface{}, 3)
	exp := []int{1, 2, 3}
	for _, v := range exp {
		ch <- v
	}
	close(ch)
	res := make([]int, 3)
	err := ForEach(ch, func(index int, item interface{}) error { res[index] = item.(int); return nil })
	assert.NoError(t, err)
	assert.Equal(t, exp, res)
}

func generateData(l int) (chan interface{}, []int) {
	data := make([]int, l)
	ch := make(chan interface{}, l)
	for i := 0; i < l; i++ {
		data[i] = int(rand.Int31n(1000))
		ch <- data[i]
	}
	close(ch)
	return ch, data
}

func benchmarkForEach(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		// need to generate a new channel here
		var ch, _ = generateData(l)
		err = ForEach(ch, func(index int, item interface{}) error {
			if item.(int) < 0 {
				return errors.New("cannot happen")
			}
			return nil
		})
	}
	assert.NoError(b, err)
}

func benchmarkLoop(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		// for fairness we regen all data here as required by the concurrent channel version
		var _, data = generateData(l)
		for _, v := range data {
			if v < 0 {
				err = errors.New("cannot happen")
			}
		}
	}
	assert.NoError(b, err)
}

func BenchmarkForEach10(b *testing.B) { benchmarkForEach(10, b) }
func BenchmarkForEach1k(b *testing.B) { benchmarkForEach(1000, b) }
func BenchmarkForEach1M(b *testing.B) { benchmarkForEach(1000000, b) }

func BenchmarkLoop10(b *testing.B) { benchmarkLoop(10, b) }
func BenchmarkLoop1k(b *testing.B) { benchmarkLoop(1000, b) }
func BenchmarkLoop1M(b *testing.B) { benchmarkLoop(1000000, b) }
