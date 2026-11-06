package dthread

// ThreadNS 线程池工厂命名空间。
type ThreadNS struct{}

// D新建 创建线程池；最大线程数 <=0 时为 1。
func (ThreadNS) D新建(maxWorkers int) *D线程池 {
	return newClient(maxWorkers)
}
