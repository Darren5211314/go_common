package dthread

// ThreadNS 线程池工厂命名空间。
type ThreadNS struct{}

// Create 创建线程池；最大线程数 <=0 时为 1。
func (ThreadNS) Create(maxWorkers int) *ThreadPool {
	return newClient(maxWorkers)
}
