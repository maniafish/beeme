package counter

import "sync/atomic"

// Counter atomic couter
type Counter uint64

var c Counter

// Incr incr one
func Incr() uint64 {
	return c.Incr()
}

// Decr decr one
func Decr() uint64 {
	return c.Decr()
}

// Get return current value
func Get() uint64 {
	return c.Get()
}

// Reset reset counter
func Reset() {
	c.Reset()
}

// New return new counter
func New() *Counter {
	c := Counter(0)
	return &c
}

// Get return current value
func (c *Counter) Get() uint64 {
	return atomic.LoadUint64((*uint64)(c))
}

// Reset reset counter
func (c *Counter) Reset() {
	atomic.StoreUint64((*uint64)(c), 0)
}

// Incr incr one
func (c *Counter) Incr() uint64 {
	return atomic.AddUint64((*uint64)(c), 1)
}

// Decr decr one
func (c *Counter) Decr() uint64 {
	return atomic.AddUint64((*uint64)(c), ^uint64(0))
}
