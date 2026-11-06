package dsync

import "sync"

// D互斥锁 有状态互斥锁。
type D互斥锁 struct {
	mu sync.Mutex
}

func newMutex() *D互斥锁 {
	return &D互斥锁{}
}

// D加锁 阻塞获取锁。
func (m *D互斥锁) D加锁() {
	m.mu.Lock()
}

// D解锁 释放锁。
func (m *D互斥锁) D解锁() {
	m.mu.Unlock()
}

// D尝试加锁 非阻塞尝试获取锁。
func (m *D互斥锁) D尝试加锁() bool {
	return m.mu.TryLock()
}
