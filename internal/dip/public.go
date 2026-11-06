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
	publicTimeout = 3 * time.Second
	publicUA = "Mozilla/5.0 (compatible; darren-ip/1.0)"
)

type publicSource struct {
	name string
	addr string
	field string // 空=整段正文当 IP；非空=JSON 点号路径
	desc string // 健康/特性注释
}

// 外网源表 接口编号 1 起连续；仅收录实测可用源。
var publicSources = map[int]publicSource{
	1:  {name: "3322", addr: "http://members.3322.org/dyndns/getip", desc: "纯文本，偏 IPv4，较快"},
	2:  {name: "gdt", addr: "https://ipv4.gdt.qq.com/get_client_ip", desc: "纯文本 IPv4"},
	3:  {name: "qq-checktime", addr: "https://vv.video.qq.com/checktime?otype=ojson", field: "ip", desc: "JSON ip，偏 IPv4"},
	4:  {name: "pconline", addr: "https://whois.pconline.com.cn/ipJson.jsp?json=true", field: "ip", desc: "JSON ip，偏 IPv4"},
	5:  {name: "ipinfo", addr: "https://ipinfo.io/ip", desc: "纯文本，偏 IPv4"},
	6:  {name: "icanhazip-v4", addr: "https://ipv4.icanhazip.com", desc: "纯文本强制 IPv4"},
	7:  {name: "126-locate", addr: "https://ipservice.ws.126.net/locate/api/getLocByIp", field: "result.ip", desc: "JSON result.ip"},
	8:  {name: "httpbin", addr: "http://httpbin.org/ip", field: "origin", desc: "JSON origin，偏慢"},
	9:  {name: "ip-api-demo", addr: "http://demo.ip-api.com/json/?lang=zh-CN", field: "query", desc: "演示接口，较慢"},
	10: {name: "ipip", addr: "https://myip.ipip.net/json", field: "data.ip", desc: "常返回 IPv6"},
	11: {name: "qq-ip2city", addr: "https://r.inews.qq.com/api/ip2city", field: "ip", desc: "常返回 IPv6"},
	12: {name: "icanhazip", addr: "https://icanhazip.com", desc: "纯文本，可能 IPv6"},
	13: {name: "api64-ipify", addr: "https://api64.ipify.org", desc: "纯文本，可能 IPv6"},
}

func isPublicIP(text string) bool {
	ip := parseIP(strings.TrimSpace(text))
	if ip == nil {
		return false
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		return isUsablePublicIPv4(ip)
	}
	return true
}

func sliceJSONObject(body string) string {
	from := strings.Index(body, "{")
	if from < 0 {
		return strings.TrimSpace(body)
	}
	return body[from:]
}

func jsonTextAtPath(body, path string) string {
	body = sliceJSONObject(body)
	var root any
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return ""
	}
	cur := root
	for _, seg := range strings.Split(path, ".") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		cur, ok = obj[seg]
		if !ok {
			return ""
		}
	}
	value, ok := cur.(string)
	if !ok {
		return ""
	}
	if i := strings.IndexByte(value, ','); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}

func ipFromSource(src publicSource, body string) (string, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return "", fmt.Errorf("外网IP源 %s 响应为空", src.name)
	}
	candidate := body
	if src.field != "" {
		candidate = jsonTextAtPath(body, src.field)
		if candidate == "" {
			return "", fmt.Errorf("外网IP源 %s 无法解析字段 %s", src.name, src.field)
		}
	} else {
		candidate = strings.TrimSpace(strings.Split(body, "\n")[0])
	}
	candidate = strings.TrimSpace(candidate)
	if !isPublicIP(candidate) {
		return "", fmt.Errorf("外网IP源 %s 返回非公网地址: %q", src.name, candidate)
	}
	return net.ParseIP(candidate).String(), nil
}

func fetchPublicBody(ctx context.Context, src publicSource) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.addr, nil)
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 构造请求失败: %w", src.name, err)
	}
	req.Header.Set("User-Agent", publicUA)
	req.Header.Set("Accept", "*/*")
	client := &http.Client{Timeout: publicTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 请求失败: %w", src.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("外网IP源 %s HTTP %d", src.name, resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("外网IP源 %s 读响应失败: %w", src.name, err)
	}
	return string(data), nil
}

// 取外网IP 按接口编号单源请求；未知编号或失败返回 error。
func publicIP(iface int) (string, error) {
	src, ok := publicSources[iface]
	if !ok {
		return "", fmt.Errorf("未知外网IP接口: %d", iface)
	}
	ctx, cancel := context.WithTimeout(context.Background(), publicTimeout)
	defer cancel()
	body, err := fetchPublicBody(ctx, src)
	if err != nil {
		return "", err
	}
	return ipFromSource(src, body)
}
