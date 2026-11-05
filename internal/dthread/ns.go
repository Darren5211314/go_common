package dthread

// Thread类 线程池工厂命名空间。
type Thread类 struct{}

// D新建 创建线程池；最大线程数 <=0 时为 1。
func (Thread类) D新建(最大线程数 int) *D线程池 {
	return 新建(最大线程数)
}
