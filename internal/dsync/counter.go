package dsync

import "sync/atomic"

// D原子计数 有状态原子 int64 计数器。
type D原子计数 struct {
	值 atomic.Int64
}

func 新建原子计数() *D原子计数 {
	return &D原子计数{}
}

func 新建原子计数初值(初值 int64) *D原子计数 {
	c := &D原子计数{}
	c.值.Store(初值)
	return c
}

// D加 原子加上增量，返回新值。
func (c *D原子计数) D加(增量 int64) int64 {
	return c.值.Add(增量)
}

// D减 原子减去减量，返回新值。
func (c *D原子计数) D减(减量 int64) int64 {
	return c.值.Add(-减量)
}

// D取 读取当前值。
func (c *D原子计数) D取() int64 {
	return c.值.Load()
}

// D设置 直接存入新值。
func (c *D原子计数) D设置(值 int64) {
	c.值.Store(值)
}

// D比较交换 仅当当前值为旧值时改为新值。
func (c *D原子计数) D比较交换(旧值, 新值 int64) bool {
	return c.值.CompareAndSwap(旧值, 新值)
}
