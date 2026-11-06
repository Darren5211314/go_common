package dhttp

// HTTP类 HTTP 命名空间：默认可直接 Send；也可用 D新建 独立客户端。
type HTTP类 struct{}

// D新建 创建独立 HTTP 客户端。
func (HTTP类) D新建() *D客户端 {
	return newClient()
}

// Send 使用默认客户端发送；选项可省略。
func (HTTP类) Send(method, url string, opts ...D发送选项) (*D响应, error) {
	return defaultClient().Send(method, url, opts...)
}

// D设置代理池 配置默认客户端的代理池。
func (HTTP类) D设置代理池(pool proxyPool) {
	defaultClient().D设置代理池(pool)
}

// D设置默认超时 配置默认客户端超时（秒）。
func (HTTP类) D设置默认超时(sec int) {
	defaultClient().D设置默认超时(sec)
}

// D设置默认重试 配置默认客户端最大尝试次数。
func (HTTP类) D设置默认重试(n int) {
	defaultClient().D设置默认重试(n)
}

// D新建会话 基于默认客户端创建会话。
func (HTTP类) D新建会话() *D会话 {
	return defaultClient().D新建会话()
}
