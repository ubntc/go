package errchan

// Collector collects errors.
type Collector struct {
	*chanGroup
	done chan struct{}
}

// NewGroup creates a new ChanGroup and starts a collector goroutine to collect errors from ChanGroup.C.
func NewGroup() (*Collector, chan<- error) {
	g := newChanGroup(10)
	c := &Collector{chanGroup: g, done: make(chan struct{})}
	go c.collect()
	return c, g.errChan
}

// collect reads errors from the channel into a slice.
func (c *Collector) collect() {
	defer close(c.done)
	var errs []error
	for v := range c.errChan {
		errs = append(errs, v)
	}
	c.mu.Lock()
	c.errList = errs
	defer c.mu.Unlock()
}

// Wait returns the error channel.
func (c *Collector) Wait() {
	c.WaitGroup.Wait()
	close(c.errChan)
	<-c.done
}
