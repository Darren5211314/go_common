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

// queryParam 从 URL 或查询串取指定参数；找不到或非法返回空串。
func queryParam(rawURL, key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	if u, err := url.Parse(rawURL); err == nil && u.RawQuery != "" {
		return u.Query().Get(key)
	}

	query := rawURL
	if i := strings.Index(rawURL, "?"); i >= 0 {
		query = rawURL[i+1:]
		if j := strings.Index(query, "#"); j >= 0 {
			query = query[:j]
		}
	} else if strings.Contains(rawURL, "=") || strings.Contains(rawURL, "&") {
		// 无 ? 的纯参数串：a=1&b=2
		query = rawURL
	} else {
		return ""
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return ""
	}
	return values.Get(key)
}
