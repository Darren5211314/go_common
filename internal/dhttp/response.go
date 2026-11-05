package dhttp

import (
	"net/http"
	"strings"
)

// D响应 HTTP 响应封装。
type D响应 struct {
	状态码   int
	正文     []byte
	重定向地址 string
	返回Cookie string
	全部Cookie string
	协议头   http.Header
	使用代理  string
	耗时毫秒  int64
}

func (r *D响应) D取状态码() int {
	if r == nil {
		return 0
	}
	return r.状态码
}

func (r *D响应) D取正文() string {
	if r == nil {
		return ""
	}
	return string(r.正文)
}

func (r *D响应) D取字节() []byte {
	if r == nil {
		return nil
	}
	return r.正文
}

func (r *D响应) D取重定向地址() string {
	if r == nil {
		return ""
	}
	return r.重定向地址
}

func (r *D响应) D取返回Cookie() string {
	if r == nil {
		return ""
	}
	return r.返回Cookie
}

func (r *D响应) D取Cookie() string {
	if r == nil {
		return ""
	}
	return r.全部Cookie
}

func (r *D响应) D取返回协议头() map[string][]string {
	if r == nil || r.协议头 == nil {
		return map[string][]string{}
	}
	结果 := make(map[string][]string, len(r.协议头))
	for 键, 值 := range r.协议头 {
		复制 := append([]string(nil), 值...)
		结果[键] = 复制
	}
	return 结果
}

func (r *D响应) D取协议头值(键 string) string {
	if r == nil || r.协议头 == nil {
		return ""
	}
	值列表 := r.协议头.Values(键)
	if len(值列表) == 0 {
		// 兼容大小写：遍历
		键小写 := strings.ToLower(键)
		for 现键, 现值 := range r.协议头 {
			if strings.ToLower(现键) == 键小写 {
				return strings.Join(现值, ", ")
			}
		}
		return ""
	}
	return strings.Join(值列表, ", ")
}

func (r *D响应) D取使用代理() string {
	if r == nil {
		return ""
	}
	return r.使用代理
}

func (r *D响应) D取耗时() int64 {
	if r == nil {
		return 0
	}
	return r.耗时毫秒
}
