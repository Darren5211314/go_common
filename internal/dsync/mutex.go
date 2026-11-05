package dsync

import "sync"

// D互斥锁 有状态互斥锁。
type D互斥锁 struct {
	锁 sync.Mutex
}

func 新建互斥锁() *D互斥锁 {
	return &D互斥锁{}
}

// D加锁 阻塞获取锁。
func (m *D互斥锁) D加锁() {
	m.锁.Lock()
}

// D解锁 释放锁。
func (m *D互斥锁) D解锁() {
	m.锁.Unlock()
}

// D尝试加锁 非阻塞尝试获取锁。
func (m *D互斥锁) D尝试加锁() bool {
	return m.锁.TryLock()
}
