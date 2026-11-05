package dsync

import "sync"

// D读写锁 有状态读写锁。
type D读写锁 struct {
	锁 sync.RWMutex
}

func 新建读写锁() *D读写锁 {
	return &D读写锁{}
}

// D读加锁 获取读锁。
func (m *D读写锁) D读加锁() {
	m.锁.RLock()
}

// D读解锁 释放读锁。
func (m *D读写锁) D读解锁() {
	m.锁.RUnlock()
}

// D写加锁 获取写锁。
func (m *D读写锁) D写加锁() {
	m.锁.Lock()
}

// D写解锁 释放写锁。
func (m *D读写锁) D写解锁() {
	m.锁.Unlock()
}

// D尝试读加锁 非阻塞尝试获取读锁。
func (m *D读写锁) D尝试读加锁() bool {
	return m.锁.TryRLock()
}

// D尝试写加锁 非阻塞尝试获取写锁。
func (m *D读写锁) D尝试写加锁() bool {
	return m.锁.TryLock()
}
