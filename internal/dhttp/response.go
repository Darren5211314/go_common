package dhttp

import (
	"net/http"
	"strings"
)

// D响应 HTTP 响应封装。
type D响应 struct {
	statusCode   int
	body     []byte
	redirectURL string
	setCookie string
	allCookie string
	headers   http.Header
	useProxy  string
	elapsedMs  int64
}

func (r *D响应) D取状态码() int {
	if r == nil {
		return 0
	}
	return r.statusCode
}

func (r *D响应) D取正文() string {
	if r == nil {
		return ""
	}
	return string(r.body)
}

func (r *D响应) D取字节() []byte {
	if r == nil {
		return nil
	}
	return r.body
}

func (r *D响应) D取重定向地址() string {
	if r == nil {
		return ""
	}
	return r.redirectURL
}

func (r *D响应) D取返回Cookie() string {
	if r == nil {
		return ""
	}
	return r.setCookie
}

func (r *D响应) D取Cookie() string {
	if r == nil {
		return ""
	}
	return r.allCookie
}

func (r *D响应) D取返回协议头() map[string][]string {
	if r == nil || r.headers == nil {
		return map[string][]string{}
	}
	result := make(map[string][]string, len(r.headers))
	for k, value := range r.headers {
		clone := append([]string(nil), value...)
		result[k] = clone
	}
	return result
}

func (r *D响应) D取协议头值(k string) string {
	if r == nil || r.headers == nil {
		return ""
	}
	values := r.headers.Values(k)
	if len(values) == 0 {
		// 兼容大小写：遍历
		keyLower := strings.ToLower(k)
		for curKey, cur := range r.headers {
			if strings.ToLower(curKey) == keyLower {
				return strings.Join(cur, ", ")
			}
		}
		return ""
	}
	return strings.Join(values, ", ")
}

func (r *D响应) D取使用代理() string {
	if r == nil {
		return ""
	}
	return r.useProxy
}

func (r *D响应) D取耗时() int64 {
	if r == nil {
		return 0
	}
	return r.elapsedMs
}
