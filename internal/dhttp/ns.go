package dhttp

// HTTPNS HTTP 命名空间：Send / NewSession / SetProxyPool（共用默认客户端）。
type HTTPNS struct{}

// Send 使用默认客户端发送；选项可省略（超时/重试未填则用内置默认 10s / 3 次）。
func (HTTPNS) Send(method, url string, opts ...SendOptions) (*D响应, error) {
	return defaultClient().Send(method, url, opts...)
}

// SetProxyPool 配置默认客户端的代理池。
func (HTTPNS) SetProxyPool(pool proxyPool) {
	defaultClient().SetProxyPool(pool)
}

// NewSession 基于默认客户端创建会话。
func (HTTPNS) NewSession() *D会话 {
	return defaultClient().NewSession()
}
