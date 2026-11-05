package dip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	外网超时 = 3 * time.Second
	外网UA = "Mozilla/5.0 (compatible; darren-ip/1.0)"
)

type 外网源 struct {
	名称 string
	地址 string
	字段 string // 空=整段正文当 IP；非空=JSON 点号路径
	说明 string // 健康/特性注释
}

// 外网源表 接口编号 1 起连续；仅收录实测可用源。
var 外网源表 = map[int]外网源{
	1:  {名称: "3322", 地址: "http://members.3322.org/dyndns/getip", 说明: "纯文本，偏 IPv4，较快"},
	2:  {名称: "gdt", 地址: "https://ipv4.gdt.qq.com/get_client_ip", 说明: "纯文本 IPv4"},
	3:  {名称: "qq-checktime", 地址: "https://vv.video.qq.com/checktime?otype=ojson", 字段: "ip", 说明: "JSON ip，偏 IPv4"},
	4:  {名称: "pconline", 地址: "https://whois.pconline.com.cn/ipJson.jsp?json=true", 字段: "ip", 说明: "JSON ip，偏 IPv4"},
	5:  {名称: "ipinfo", 地址: "https://ipinfo.io/ip", 说明: "纯文本，偏 IPv4"},
	6:  {名称: "icanhazip-v4", 地址: "https://ipv4.icanhazip.com", 说明: "纯文本强制 IPv4"},
	7:  {名称: "126-locate", 地址: "https://ipservice.ws.126.net/locate/api/getLocByIp", 字段: "result.ip", 说明: "JSON result.ip"},
	8:  {名称: "httpbin", 地址: "http://httpbin.org/ip", 字段: "origin", 说明: "JSON origin，偏慢"},
	9:  {名称: "ip-api-demo", 地址: "http://demo.ip-api.com/json/?lang=zh-CN", 字段: "query", 说明: "演示接口，较慢"},
	10: {名称: "ipip", 地址: "https://myip.ipip.net/json", 字段: "data.ip", 说明: "常返回 IPv6"},
	11: {名称: "qq-ip2city", 地址: "https://r.inews.qq.com/api/ip2city", 字段: "ip", 说明: "常返回 IPv6"},
	12: {名称: "icanhazip", 地址: "https://icanhazip.com", 说明: "纯文本，可能 IPv6"},
	13: {名称: "api64-ipify", 地址: "https://api64.ipify.org", 说明: "纯文本，可能 IPv6"},
}

func 是公网IP(文本 string) bool {
	ip := 解析IP(strings.TrimSpace(文本))
	if ip == nil {
		return false
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		return 是可用公网IPv4(ip)
	}
	return true
}

func 截取JSON对象(正文 string) string {
	起 := strings.Index(正文, "{")
	if 起 < 0 {
		return strings.TrimSpace(正文)
	}
	return 正文[起:]
}

func json按路径取文本(正文, 路径 string) string {
	正文 = 截取JSON对象(正文)
	var 根 any
	if err := json.Unmarshal([]byte(正文), &根); err != nil {
		return ""
	}
	当前 := 根
	for _, 段 := range strings.Split(路径, ".") {
		对象, ok := 当前.(map[string]any)
		if !ok {
			return ""
		}
		当前, ok = 对象[段]
		if !ok {
			return ""
		}
	}
	值, ok := 当前.(string)
	if !ok {
		return ""
	}
	if i := strings.IndexByte(值, ','); i >= 0 {
		值 = 值[:i]
	}
	return strings.TrimSpace(值)
}

func 从源解析IP(源 外网源, 正文 string) (string, error) {
	正文 = strings.TrimSpace(正文)
	if 正文 == "" {
		return "", fmt.Errorf("外网IP源 %s 响应为空", 源.名称)
	}
	候选 := 正文
	if 源.字段 != "" {
		候选 = json按路径取文本(正文, 源.字段)
		if 候选 == "" {
			return "", fmt.Errorf("外网IP源 %s 无法解析字段 %s", 源.名称, 源.字段)
		}
	} else {
		候选 = strings.TrimSpace(strings.Split(正文, "\n")[0])
	}
	候选 = strings.TrimSpace(候选)
	if !是公网IP(候选) {
		return "", fmt.Errorf("外网IP源 %s 返回非公网地址: %q", 源.名称, 候选)
	}
	return net.ParseIP(候选).String(), nil
}

func 拉取外网正文(ctx context.Context, 源 外网源) (string, error) {
	请求, err := http.NewRequestWithContext(ctx, http.MethodGet, 源.地址, nil)
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 构造请求失败: %w", 源.名称, err)
	}
	请求.Header.Set("User-Agent", 外网UA)
	请求.Header.Set("Accept", "*/*")
	客户端 := &http.Client{Timeout: 外网超时}
	响应, err := 客户端.Do(请求)
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 请求失败: %w", 源.名称, err)
	}
	defer 响应.Body.Close()
	if 响应.StatusCode < 200 || 响应.StatusCode >= 300 {
		return "", fmt.Errorf("外网IP源 %s HTTP %d", 源.名称, 响应.StatusCode)
	}
	数据, err := io.ReadAll(io.LimitReader(响应.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 读响应失败: %w", 源.名称, err)
	}
	return string(数据), nil
}

// 取外网IP 按接口编号单源请求；未知编号或失败返回 error。
func 取外网IP(接口 int) (string, error) {
	源, ok := 外网源表[接口]
	if !ok {
		return "", fmt.Errorf("未知外网IP接口: %d", 接口)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 外网超时)
	defer cancel()
	正文, err := 拉取外网正文(ctx, 源)
	if err != nil {
		return "", err
	}
	return 从源解析IP(源, 正文)
}
