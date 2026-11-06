package dsync

import "sync"

// Mutex 有状态互斥锁。
type Mutex struct {
	mu sync.Mutex
}

func newMutex() *Mutex {
	return &Mutex{}
}

// Lock 阻塞获取锁。
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock 释放锁。
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock 非阻塞尝试获取锁。
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}
