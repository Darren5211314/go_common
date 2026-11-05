package dproxy

import (
	"fmt"
	"net/url"
	"strings"
)

func 规范化协议(协议 string) (string, error) {
	协议 = strings.ToLower(strings.TrimSpace(协议))
	switch 协议 {
	case "http", "https", "socks5", "socks5h":
		if 协议 == "socks5h" {
			协议 = "socks5"
		}
		return 协议, nil
	default:
		return "", fmt.Errorf("不支持的协议: %s", 协议)
	}
}

// 解析代理行 将一行文本规范为 scheme://[user:pass@]host:port。
func 解析代理行(行, 默认协议 string) (string, error) {
	行 = strings.TrimSpace(行)
	if 行 == "" || strings.HasPrefix(行, "#") {
		return "", fmt.Errorf("空行")
	}
	默认协议, 错误 := 规范化协议(默认协议)
	if 错误 != nil {
		默认协议 = "http"
	}
	if !strings.Contains(行, "://") {
		行 = 默认协议 + "://" + 行
	}
	u, 错误 := url.Parse(行)
	if 错误 != nil {
		return "", 错误
	}
	协议, 错误 := 规范化协议(u.Scheme)
	if 错误 != nil {
		return "", 错误
	}
	if u.Host == "" {
		return "", fmt.Errorf("缺少主机")
	}
	// 强制带端口时 Host 已含端口；无端口也允许（少见）
	u.Scheme = 协议
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func 代理键(代理URL string) string {
	u, 错误 := url.Parse(代理URL)
	if 错误 != nil {
		return 代理URL
	}
	用户 := ""
	if u.User != nil {
		用户 = u.User.String() + "@"
	}
	return strings.ToLower(u.Scheme) + "://" + 用户 + u.Host
}
