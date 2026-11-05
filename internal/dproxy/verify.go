package dproxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	xproxy "golang.org/x/net/proxy"
)

func 验证代理(ctx context.Context, 代理URL, 验证地址 string, 超时 time.Duration) error {
	u, 错误 := url.Parse(代理URL)
	if 错误 != nil {
		return 错误
	}
	客户端 := &http.Client{Timeout: 超时}
	传输 := &http.Transport{
		DisableKeepAlives: true,
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		传输.Proxy = http.ProxyURL(u)
	case "socks5":
		拨号, 错误 := 新建SOCKS5拨号(u)
		if 错误 != nil {
			return 错误
		}
		传输.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return 拨号.Dial(network, address)
		}
	default:
		return fmt.Errorf("不支持的协议: %s", u.Scheme)
	}
	客户端.Transport = 传输

	请求, 错误 := http.NewRequestWithContext(ctx, http.MethodGet, 验证地址, nil)
	if 错误 != nil {
		return 错误
	}
	响应, 错误 := 客户端.Do(请求)
	if 错误 != nil {
		return 错误
	}
	_ = 响应.Body.Close()
	if 响应.StatusCode >= 400 {
		return fmt.Errorf("验证地址返回状态码 %d", 响应.StatusCode)
	}
	return nil
}

func 新建SOCKS5拨号(u *url.URL) (xproxy.Dialer, error) {
	var 认证 *xproxy.Auth
	if u.User != nil {
		密码, _ := u.User.Password()
		认证 = &xproxy.Auth{User: u.User.Username(), Password: 密码}
	}
	return xproxy.SOCKS5("tcp", u.Host, 认证, xproxy.Direct)
}
