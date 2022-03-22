package buffertest

import "sync"

type metrics struct {
	sync.RWMutex
	value int
}

var mxNumHandled = &metrics{}

func (mx *metrics) Get() int {
	mx.RLock()
	defer mx.RUnlock()
	return mx.value
}

func (mx *metrics) Inc() int {
	mx.Lock()
	defer mx.Unlock()
	mx.value++
	return mx.value
}
