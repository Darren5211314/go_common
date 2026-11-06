package dhttp

import (
	"net/http"
	"strings"
)

// D响应 HTTP 响应封装。
type D响应 struct {
	statusCode  int
	body        []byte
	redirectURL string
	setCookie   string
	allCookie   string
	headers     http.Header
	useProxy    string
	elapsedMs   int64
}

func (r *D响应) GetStatusCode() int {
	if r == nil {
		return 0
	}
	return r.statusCode
}

func (r *D响应) GetBody() string {
	if r == nil {
		return ""
	}
	return string(r.body)
}

func (r *D响应) GetBytes() []byte {
	if r == nil {
		return nil
	}
	return r.body
}

func (r *D响应) GetLocation() string {
	if r == nil {
		return ""
	}
	return r.redirectURL
}

func (r *D响应) GetSetCookie() string {
	if r == nil {
		return ""
	}
	return r.setCookie
}

func (r *D响应) GetCookie() string {
	if r == nil {
		return ""
	}
	return r.allCookie
}

func (r *D响应) GetHeaders() map[string][]string {
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

func (r *D响应) GetHeader(name string) string {
	if r == nil || r.headers == nil {
		return ""
	}
	values := r.headers.Values(name)
	if len(values) == 0 {
		// 兼容大小写：遍历
		keyLower := strings.ToLower(name)
		for curKey, cur := range r.headers {
			if strings.ToLower(curKey) == keyLower {
				return strings.Join(cur, ", ")
			}
		}
		return ""
	}
	return strings.Join(values, ", ")
}

func (r *D响应) GetUseProxy() string {
	if r == nil {
		return ""
	}
	return r.useProxy
}

func (r *D响应) GetElapsedMs() int64 {
	if r == nil {
		return 0
	}
	return r.elapsedMs
}
