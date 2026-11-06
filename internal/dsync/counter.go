package dsync

import "sync/atomic"

// D原子计数 有状态原子 int64 计数器。
type D原子计数 struct {
	value atomic.Int64
}

func newAtomicCounter() *D原子计数 {
	return &D原子计数{}
}

func newAtomicCounterAt(initial int64) *D原子计数 {
	c := &D原子计数{}
	c.value.Store(initial)
	return c
}

// D加 原子加上增量，返回新值。
func (c *D原子计数) D加(delta int64) int64 {
	return c.value.Add(delta)
}

// D减 原子减去减量，返回新值。
func (c *D原子计数) D减(delta int64) int64 {
	return c.value.Add(-delta)
}

// D取 读取当前值。
func (c *D原子计数) D取() int64 {
	return c.value.Load()
}

// D设置 直接存入新值。
func (c *D原子计数) D设置(value int64) {
	c.value.Store(value)
}

// D比较交换 仅当当前值为旧值时改为新值。
func (c *D原子计数) D比较交换(oldVal, newVal int64) bool {
	return c.value.CompareAndSwap(oldVal, newVal)
}
