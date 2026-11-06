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

func verifyProxy(ctx context.Context, proxyURL, verifyURL string, timeout time.Duration) error {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: timeout}
	transport := &http.Transport{
		DisableKeepAlives: true,
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		transport.Proxy = http.ProxyURL(u)
	case "socks5":
		dial, err := newSOCKS5Dialer(u)
		if err != nil {
			return err
		}
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return dial.Dial(network, address)
		}
	default:
		return fmt.Errorf("不支持的协议: %s", u.Scheme)
	}
	client.Transport = transport

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, verifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("验证地址返回状态码 %d", resp.StatusCode)
	}
	return nil
}

func newSOCKS5Dialer(u *url.URL) (xproxy.Dialer, error) {
	var auth *xproxy.Auth
	if u.User != nil {
		password, _ := u.User.Password()
		auth = &xproxy.Auth{User: u.User.Username(), Password: password}
	}
	return xproxy.SOCKS5("tcp", u.Host, auth, xproxy.Direct)
}
