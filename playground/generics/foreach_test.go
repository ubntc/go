package generics

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestForEach(t *testing.T) {
	ch := make(chan any, 3)
	exp := []int{1, 2, 3}
	for _, v := range exp {
		ch <- v
	}
	close(ch)
	res := make([]int, 3)
	err := ForEachChan(ch, func(index int, item any) error { res[index] = item.(int); return nil })
	assert.NoError(t, err)
	assert.Equal(t, exp, res)

	err = ForEach(exp, func(index int, item int) error { res[index] = item; return nil })
	assert.NoError(t, err)
	assert.Equal(t, exp, res)
}

func payload(i int) {
	// spend some time to simulate work
	type Evt struct {
		ID        string    `json:"id"`
		Timestamp time.Time `json:"timestamp"`
		Data      string    `json:"data"`
	}
	ts := time.Now().Add(time.Duration(i) * time.Millisecond)
	ev := &Evt{
		ID:        fmt.Sprintf("event-%d", i),
		Timestamp: ts,
		Data:      fmt.Sprintf("data-%d", i),
	}
	data, err := json.Marshal(ev)
	if err != nil {
		panic(err)
	}
	ev = &Evt{}
	err = json.Unmarshal(data, ev)
	if err != nil {
		panic(err)
	}
}

func generateData(l int, b *testing.B) (chan any, []int) {
	b.StartTimer()
	defer b.StartTimer()
	data := make([]int, l)
	ch := make(chan any, l)
	for i := 0; i < l; i++ {
		data[i] = int(rand.Int31n(1000))
		ch <- data[i]
	}
	close(ch)
	return ch, data
}

func benchmarkForEachChan(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		// need to generate a new channel here
		ch, _ := generateData(l, b)
		err = ForEachChan(ch, func(index int, item any) error {
			payload(item.(int))
			return nil
		})
	}
	assert.NoError(b, err)
}

func benchmarkForEachGeneric(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		_, data := generateData(l, b)
		err = ForEach(data, func(index int, item int) error {
			payload(item)
			return nil
		})
	}
	assert.NoError(b, err)
}

func benchmarkLoopSequentiual(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		// for fairness we regen all data here as required by the concurrent channel version
		_, data := generateData(l, b)
		for i := range data {
			payload(data[i])
		}
	}
	assert.NoError(b, err)
}

func benchmarkForLoopConcurrent(l int, b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		_, data := generateData(l, b)
		errg := new(errgroup.Group)
		for i := range data {
			errg.Go(func() error {
				payload(data[i])
				return nil
			})
		}
		errg.Wait()
	}
	assert.NoError(b, err)
}

func BenchmarkAll(b *testing.B) {
	for _, n := range []int{10, 1000} {
		b.Run("foreach-generic-"+fmt.Sprint(n), func(b *testing.B) { benchmarkForEachGeneric(n, b) })
		b.Run("foreach-chan-"+fmt.Sprint(n), func(b *testing.B) { benchmarkForEachChan(n, b) })
		b.Run("loop-sequential-"+fmt.Sprint(n), func(b *testing.B) { benchmarkLoopSequentiual(n, b) })
		b.Run("loop-concurrent-"+fmt.Sprint(n), func(b *testing.B) { benchmarkForLoopConcurrent(n, b) })
	}
}
