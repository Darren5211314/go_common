package dproxy

// ProxyNS 代理池工厂命名空间。
type ProxyNS struct{}

// D新建 创建未启动的代理池（带默认配置）。
func (ProxyNS) D新建() *D代理池 {
	return newClient()
}
