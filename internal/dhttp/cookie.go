package dhttp

import (
	"net/http"
	"sort"
	"strings"
)

func 解析Cookie串(文本 string) map[string]string {
	结果 := make(map[string]string)
	文本 = strings.TrimSpace(文本)
	if 文本 == "" {
		return 结果
	}
	// 支持 a=1; b=2 与 a=1&b=2
	文本 = strings.ReplaceAll(文本, "&", ";")
	for _, 段 := range strings.Split(文本, ";") {
		段 = strings.TrimSpace(段)
		if 段 == "" {
			continue
		}
		键, 值, ok := strings.Cut(段, "=")
		键 = strings.TrimSpace(键)
		if 键 == "" || !ok {
			continue
		}
		结果[键] = strings.TrimSpace(值)
	}
	return 结果
}

func 合并Cookie映射(多份 ...map[string]string) map[string]string {
	结果 := make(map[string]string)
	for _, 一份 := range 多份 {
		for 键, 值 := range 一份 {
			结果[键] = 值
		}
	}
	return 结果
}

func Cookie映射转串(映射 map[string]string) string {
	if len(映射) == 0 {
		return ""
	}
	键列表 := make([]string, 0, len(映射))
	for 键 := range 映射 {
		键列表 = append(键列表, 键)
	}
	sort.Strings(键列表)
	部分 := make([]string, 0, len(键列表))
	for _, 键 := range 键列表 {
		部分 = append(部分, 键+"="+映射[键])
	}
	return strings.Join(部分, "; ")
}

func 从响应提取Cookie(头 http.Header) map[string]string {
	结果 := make(map[string]string)
	for _, 原始 := range 头.Values("Set-Cookie") {
		if 原始 == "" {
			continue
		}
		第一 := strings.SplitN(原始, ";", 2)[0]
		键, 值, 成功 := strings.Cut(第一, "=")
		if !成功 {
			continue
		}
		键 = strings.TrimSpace(键)
		if 键 == "" {
			continue
		}
		结果[键] = strings.TrimSpace(值)
	}
	return 结果
}
