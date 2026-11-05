package dhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	xproxy "golang.org/x/net/proxy"
)

var (
	直连传输一次 sync.Once
	直连传输   *http.Transport
)

// 取直连传输 返回可 KeepAlive 复用的共享 Transport（无代理）。
func 取直连传输() *http.Transport {
	直连传输一次.Do(func() {
		直连传输 = &http.Transport{
			Proxy:                 nil,
			DisableKeepAlives:     false,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	})
	return 直连传输
}

// 新建代理传输 为指定代理新建传输（禁用 KeepAlive，避免串代理粘连）。
func 新建代理传输(代理URL string) (*http.Transport, error) {
	传输 := &http.Transport{
		DisableKeepAlives: true,
		Proxy:             nil,
	}
	代理URL = strings.TrimSpace(代理URL)
	if 代理URL == "" {
		return nil, fmt.Errorf("代理URL为空")
	}
	u, 错误 := url.Parse(代理URL)
	if 错误 != nil {
		return nil, 错误
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		传输.Proxy = http.ProxyURL(u)
	case "socks5", "socks5h":
		拨号, 错误 := 新建SOCKS5拨号(u)
		if 错误 != nil {
			return nil, 错误
		}
		传输.DialContext = 包装SOCKS5拨号(拨号)
	default:
		return nil, fmt.Errorf("不支持的代理协议: %s", u.Scheme)
	}
	return 传输, nil
}

func 包装SOCKS5拨号(拨号 xproxy.Dialer) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		类型结果 := make(chan struct {
			连接 net.Conn
			错  error
		}, 1)
		go func() {
			连接, 错 := 拨号.Dial(network, address)
			类型结果 <- struct {
				连接 net.Conn
				错  error
			}{连接, 错}
		}()
		select {
		case <-ctx.Done():
			go func() {
				结果 := <-类型结果
				if 结果.连接 != nil {
					_ = 结果.连接.Close()
				}
			}()
			return nil, ctx.Err()
		case 结果 := <-类型结果:
			if err := ctx.Err(); err != nil {
				if 结果.连接 != nil {
					_ = 结果.连接.Close()
				}
				return nil, err
			}
			return 结果.连接, 结果.错
		}
	}
}

func 新建SOCKS5拨号(u *url.URL) (xproxy.Dialer, error) {
	var 认证 *xproxy.Auth
	if u.User != nil {
		密码, _ := u.User.Password()
		认证 = &xproxy.Auth{User: u.User.Username(), Password: 密码}
	}
	return xproxy.SOCKS5("tcp", u.Host, 认证, xproxy.Direct)
}
