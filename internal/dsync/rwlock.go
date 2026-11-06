package dsync

import "sync"

// RWMutex 有状态读写锁。
type RWMutex struct {
	mu sync.RWMutex
}

func newRWMutex() *RWMutex {
	return &RWMutex{}
}

// RLock 获取读锁。
func (m *RWMutex) RLock() {
	m.mu.RLock()
}

// RUnlock 释放读锁。
func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
}

// Lock 获取写锁。
func (m *RWMutex) Lock() {
	m.mu.Lock()
}

// Unlock 释放写锁。
func (m *RWMutex) Unlock() {
	m.mu.Unlock()
}

// TryRLock 非阻塞尝试获取读锁。
func (m *RWMutex) TryRLock() bool {
	return m.mu.TryRLock()
}

// TryLock 非阻塞尝试获取写锁。
func (m *RWMutex) TryLock() bool {
	return m.mu.TryLock()
}
