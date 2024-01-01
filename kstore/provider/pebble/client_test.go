package pebble_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/kstore/provider/pebble"
)

func Setup(t *testing.T) *pebble.Client {
	c := pebble.NewClient(t.TempDir())
	errs, err := c.CreateTopics(context.Background(), "test")
	assert.NoError(t, err)
	assert.Len(t, errs, 0)
	return c
}

func TestWrite(t *testing.T) {
	c := Setup(t)
	defer c.Close()
	w := c.NewWriter().(*pebble.Writer)

	ctx := context.Background()

	// start writer
	msgIn := Msg("test", 0, k, v)
	err := w.Write(ctx, "test", &msgIn)
	assert.NoError(t, err)

	m, err := c.Get("test", pebble.StorageKey(&msgIn))
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestReadAndWrite(t *testing.T) {
	c := Setup(t)
	defer c.Close()
	r := c.NewReader("test").(*pebble.Reader)
	w := c.NewWriter().(*pebble.Writer)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// start reader
	go func() {
		defer cancel()
		msgOut, err := r.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, msgOut)
	}()

	// start writer
	time.Sleep(time.Millisecond * 10)
	msgIn := Msg("test", 0, k, v)
	err := w.Write(ctx, "test", &msgIn)
	assert.NoError(t, err)

	m, err := r.Get(pebble.StorageKey(&msgIn))
	assert.NoError(t, err)
	assert.NotNil(t, m)

	// wait for reader to finish
	<-ctx.Done()
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
}

func TestReadLast(t *testing.T) {
	c := Setup(t)
	defer c.Close()
	r := pebble.NewReader(c, "test", pebble.StartOffsetLast)
	w := c.NewWriter().(*pebble.Writer)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	samples := 10

	// start writer
	offset := 0
	write := func(n int) {
		for i := 0; i < n; i++ {
			msgIn := Msg("test", uint64(offset), k, v)
			err := w.Write(ctx, "test", &msgIn)
			offset++
			assert.NoError(t, err)
		}
	}

	write(samples)

	received := 0

	// start reader
	go func() {
		defer cancel()

		for received < samples {
			msgOut, err := r.Read(ctx)
			assert.NoError(t, err)
			assert.NotNil(t, msgOut)
			if err != nil {
				return
			}
			if msgOut == nil {
				log.Println("msgOut is nil")
				return
			}
			assert.GreaterOrEqual(t, msgOut.Offset(), uint64(samples))
			received++
			r.Commit(ctx, msgOut)
		}
	}()

	time.Sleep(time.Millisecond * 10)

	write(samples)

	// wait for reader to finish
	<-ctx.Done()
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
	assert.Equal(t, samples, received)

	// check shared metrics for sane value
	reads := pebble.Metrics.GetReads("test")
	writes := pebble.Metrics.GetWrites("test")
	assert.GreaterOrEqual(t, reads[pebble.OffsetStatusNewer], samples)
	assert.GreaterOrEqual(t, writes, samples*2)
}

func TestClient(t *testing.T) {
	c := Setup(t)
	defer c.Close()
	defer c.DeleteDB("test")

	r := c.NewReader("test").(*pebble.Reader)
	assert.NotNil(t, r)
	defer r.Close()

	w := c.NewWriter().(*pebble.Writer)
	assert.NotNil(t, w)
	defer w.Close()

	ctx := context.Background()
	errs, err := c.CreateTopics(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, errs, 0)

	rw := func(offset uint64) {
		msgIn := Msg("test", offset, k, v)
		err = c.Write(ctx, "test", &msgIn)
		assert.NoError(t, err)

		msgOut, err := c.Get("test", pebble.StorageKey(&msgIn))
		assert.NoError(t, err)
		assert.NotNil(t, msgOut)
		assert.Equal(t, msgIn, msgOut.(*pebble.Message).Message)
	}

	rw(0)
	rw(10)
	rw(100)
}
