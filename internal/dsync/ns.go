package dsync

// Sync类 并发同步原语工厂命名空间。
type Sync类 struct{}

// D新建互斥锁 创建互斥锁。
func (Sync类) D新建互斥锁() *D互斥锁 {
	return 新建互斥锁()
}

// D新建读写锁 创建读写锁。
func (Sync类) D新建读写锁() *D读写锁 {
	return 新建读写锁()
}

// D新建原子计数 创建初值为 0 的原子计数器。
func (Sync类) D新建原子计数() *D原子计数 {
	return 新建原子计数()
}

// D新建原子计数初值 创建指定初值的原子计数器。
func (Sync类) D新建原子计数初值(初值 int64) *D原子计数 {
	return 新建原子计数初值(初值)
}
