package dip

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test从源解析IP(t *testing.T) {
	ip, err := 从源解析IP(外网源{名称: "t"}, "  1.2.3.4\n")
	if err != nil || ip != "1.2.3.4" {
		t.Fatalf("纯文本解析失败: %q %v", ip, err)
	}
	ip, err = 从源解析IP(外网源{名称: "t", 字段: "ip"}, `{"ip":"8.8.8.8"}`)
	if err != nil || ip != "8.8.8.8" {
		t.Fatalf("json ip 失败: %q %v", ip, err)
	}
	ip, err = 从源解析IP(外网源{名称: "t", 字段: "data.ip"}, `{"data":{"ip":"1.1.1.1"}}`)
	if err != nil || ip != "1.1.1.1" {
		t.Fatalf("嵌套路径失败: %q %v", ip, err)
	}
	ip, err = 从源解析IP(外网源{名称: "t", 字段: "origin"}, `{"origin":"9.9.9.9, 8.8.8.8"}`)
	if err != nil || ip != "9.9.9.9" {
		t.Fatalf("origin 逗号分割失败: %q %v", ip, err)
	}
	if _, err := 从源解析IP(外网源{名称: "t"}, "127.0.0.1"); err == nil {
		t.Fatal("回环应失败")
	}
	if _, err := 从源解析IP(外网源{名称: "t"}, "192.168.1.1"); err == nil {
		t.Fatal("私有应失败")
	}
}

func Test取外网IP_指定接口(t *testing.T) {
	服务 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ip":"8.8.4.4"}`))
	}))
	defer 服务.Close()

	旧 := 外网源表
	外网源表 = map[int]外网源{
		99: {名称: "mock", 地址: 服务.URL, 字段: "ip"},
	}
	defer func() { 外网源表 = 旧 }()

	if _, err := 取外网IP(0); err == nil || !strings.Contains(err.Error(), "未知外网IP接口") {
		t.Fatalf("未知编号应报错: %v", err)
	}
	IP文本, err := 取外网IP(99)
	dlog.Info("指定接口", "ip", IP文本, "错误", err)
	if err != nil || IP文本 != "8.8.4.4" {
		t.Fatalf("得到 %q %v", IP文本, err)
	}
}

func Test取外网IP_Live(t *testing.T) {
	起 := time.Now()
	IP文本, err := 取外网IP(1)
	耗时 := time.Since(起)
	dlog.Info("外网IP联调", "接口", 1, "ip", IP文本, "耗时", 耗时.String(), "错误", err)
	if err != nil || IP文本 == "" {
		t.Skip("当前网络无法取外网 IP: ", err)
	}
	if !是公网IP(IP文本) {
		t.Fatalf("非法公网 IP: %s", IP文本)
	}
}
