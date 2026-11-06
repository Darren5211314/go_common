package dproxy

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeScheme(scheme string) (string, error) {
	scheme = strings.ToLower(strings.TrimSpace(scheme))
	switch scheme {
	case "http", "https", "socks5", "socks5h":
		if scheme == "socks5h" {
			scheme = "socks5"
		}
		return scheme, nil
	default:
		return "", fmt.Errorf("不支持的协议: %s", scheme)
	}
}

// 解析代理行 将一行文本规范为 scheme://[user:pass@]host:port。
func parseProxyLine(line, defaultScheme string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", fmt.Errorf("空行")
	}
	defaultScheme, err := normalizeScheme(defaultScheme)
	if err != nil {
		defaultScheme = "http"
	}
	if !strings.Contains(line, "://") {
		line = defaultScheme + "://" + line
	}
	u, err := url.Parse(line)
	if err != nil {
		return "", err
	}
	scheme, err := normalizeScheme(u.Scheme)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return "", fmt.Errorf("缺少主机")
	}
	// 强制带端口时 Host 已含端口；无端口也允许（少见）
	u.Scheme = scheme
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func proxyKey(proxyURL string) string {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return proxyURL
	}
	user := ""
	if u.User != nil {
		user = u.User.String() + "@"
	}
	return strings.ToLower(u.Scheme) + "://" + user + u.Host
}
