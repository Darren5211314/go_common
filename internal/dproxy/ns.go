package dproxy

// Proxy类 代理池工厂命名空间。
type Proxy类 struct{}

// D新建 创建未启动的代理池（带默认配置）。
func (Proxy类) D新建() *D代理池 {
	return newClient()
}
