package dsync

// SyncNS 并发同步原语工厂命名空间。
type SyncNS struct{}

// NewMutex 创建互斥锁。
func (SyncNS) NewMutex() *Mutex {
	return newMutex()
}

// NewRWMutex 创建读写锁。
func (SyncNS) NewRWMutex() *RWMutex {
	return newRWMutex()
}

// NewAtomicInt 创建初值为 0 的原子计数器。
func (SyncNS) NewAtomicInt() *AtomicInt {
	return newAtomicCounter()
}

// NewAtomicIntValue 创建指定初值的原子计数器。
func (SyncNS) NewAtomicIntValue(initial int64) *AtomicInt {
	return newAtomicCounterAt(initial)
}
