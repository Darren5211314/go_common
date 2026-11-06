package dcodec

import (
	"net/url"
	"strings"
)

// 网址编码 对字符串做 URL query 组件编码（空格为 %20）。
func urlEncode(text string) string {
	return strings.ReplaceAll(url.QueryEscape(text), "+", "%20")
}

// 网址解码 解码 URL query 组件；失败返回空串。
func urlDecode(text string) string {
	result, err := url.QueryUnescape(text)
	if err != nil {
		return ""
	}
	return result
}
