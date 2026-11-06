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
	directOnce sync.Once
	directTripper   *http.Transport
)

// 取直连传输 返回可 KeepAlive 复用的共享 Transport（无代理）。
func directTransport() *http.Transport {
	directOnce.Do(func() {
		directTripper = &http.Transport{
			Proxy:                 nil,
			DisableKeepAlives:     false,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	})
	return directTripper
}

// 新建代理传输 为指定代理新建传输（禁用 KeepAlive，避免串代理粘连）。
func newProxyTransport(proxyURL string) (*http.Transport, error) {
	transport := &http.Transport{
		DisableKeepAlives: true,
		Proxy:             nil,
	}
	proxyURL = strings.TrimSpace(proxyURL)
	if proxyURL == "" {
		return nil, fmt.Errorf("代理URL为空")
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		transport.Proxy = http.ProxyURL(u)
	case "socks5", "socks5h":
		dial, err := newSOCKS5Dialer(u)
		if err != nil {
			return nil, err
		}
		transport.DialContext = wrapSOCKS5Dialer(dial)
	default:
		return nil, fmt.Errorf("不支持的代理协议: %s", u.Scheme)
	}
	return transport, nil
}

func wrapSOCKS5Dialer(dial xproxy.Dialer) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		typedResult := make(chan struct {
			conn net.Conn
			e  error
		}, 1)
		go func() {
			conn, e := dial.Dial(network, address)
			typedResult <- struct {
				conn net.Conn
				e  error
			}{conn, e}
		}()
		select {
		case <-ctx.Done():
			go func() {
				result := <-typedResult
				if result.conn != nil {
					_ = result.conn.Close()
				}
			}()
			return nil, ctx.Err()
		case result := <-typedResult:
			if err := ctx.Err(); err != nil {
				if result.conn != nil {
					_ = result.conn.Close()
				}
				return nil, err
			}
			return result.conn, result.e
		}
	}
}

func newSOCKS5Dialer(u *url.URL) (xproxy.Dialer, error) {
	var auth *xproxy.Auth
	if u.User != nil {
		password, _ := u.User.Password()
		auth = &xproxy.Auth{User: u.User.Username(), Password: password}
	}
	return xproxy.SOCKS5("tcp", u.Host, auth, xproxy.Direct)
}
