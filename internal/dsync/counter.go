package dsync

import "sync/atomic"

// AtomicInt 有状态原子 int64 计数器。
type AtomicInt struct {
	value atomic.Int64
}

func newAtomicCounter() *AtomicInt {
	return &AtomicInt{}
}

func newAtomicCounterAt(initial int64) *AtomicInt {
	c := &AtomicInt{}
	c.value.Store(initial)
	return c
}

// Add 原子加上增量，返回新值。
func (c *AtomicInt) Add(delta int64) int64 {
	return c.value.Add(delta)
}

// Sub 原子减去减量，返回新值。
func (c *AtomicInt) Sub(delta int64) int64 {
	return c.value.Add(-delta)
}

// Get 读取当前值。
func (c *AtomicInt) Get() int64 {
	return c.value.Load()
}

// Set 直接存入新值。
func (c *AtomicInt) Set(value int64) {
	c.value.Store(value)
}

// CompareAndSwap 仅当当前值为旧值时改为新值。
func (c *AtomicInt) CompareAndSwap(oldVal, newVal int64) bool {
	return c.value.CompareAndSwap(oldVal, newVal)
}
