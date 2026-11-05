package dcodec

import (
	"net/url"
	"strings"
)

// 网址编码 对字符串做 URL query 组件编码（空格为 %20）。
func 网址编码(文本 string) string {
	return strings.ReplaceAll(url.QueryEscape(文本), "+", "%20")
}

// 网址解码 解码 URL query 组件；失败返回空串。
func 网址解码(文本 string) string {
	结果, 错误 := url.QueryUnescape(文本)
	if 错误 != nil {
		return ""
	}
	return 结果
}
