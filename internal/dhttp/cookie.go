package dhttp

import (
	"net/http"
	"sort"
	"strings"
)

func parseCookieString(text string) map[string]string {
	result := make(map[string]string)
	text = strings.TrimSpace(text)
	if text == "" {
		return result
	}
	// 支持 a=1; b=2 与 a=1&b=2
	text = strings.ReplaceAll(text, "&", ";")
	for _, part := range strings.Split(text, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		k, value, ok := strings.Cut(part, "=")
		k = strings.TrimSpace(k)
		value = strings.TrimSpace(value)
		if k == "" || !ok || value == "" {
			continue // 空值（如 P_ck_ctl=）不写入
		}
		result[k] = value
	}
	return result
}

func mergeCookieMaps(parts ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, one := range parts {
		for k, value := range one {
			if k == "" || strings.TrimSpace(value) == "" {
				continue // 空值不覆盖已有键
			}
			result[k] = value
		}
	}
	return result
}

func cookieMapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k, value := range m {
		if k == "" || strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	chunk := make([]string, 0, len(keys))
	for _, k := range keys {
		chunk = append(chunk, k+"="+m[k])
	}
	return strings.Join(chunk, "; ")
}

func cookiesFromResponse(header http.Header) map[string]string {
	result := make(map[string]string)
	for _, raw := range header.Values("Set-Cookie") {
		if raw == "" {
			continue
		}
		first := strings.SplitN(raw, ";", 2)[0]
		k, value, ok := strings.Cut(first, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		value = strings.TrimSpace(value)
		if k == "" || value == "" {
			continue // 空值 Set-Cookie 不合并
		}
		result[k] = value
	}
	return result
}
