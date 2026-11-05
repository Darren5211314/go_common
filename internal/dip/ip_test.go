package dip

import (
	"os"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test解析与随机本机(t *testing.T) {
	dlog.Info("是否有效", "1.2.3.4", 是否有效("1.2.3.4"), "坏", 是否有效("abc"))
	if !是否有效("1.2.3.4") || 是否有效("abc") {
		t.Errorf("是否有效 不符")
	}
	if 取版本("8.8.8.8") != "IPv4" {
		t.Errorf("取版本 IPv4 不符")
	}
	if 取版本("::1") != "IPv6" {
		t.Errorf("取版本 IPv6 不符")
	}
	if !是否私有("192.168.1.1") || 是否私有("8.8.8.8") {
		t.Errorf("是否私有 不符")
	}
	if !是否回环("127.0.0.1") || 是否回环("8.8.8.8") {
		t.Errorf("是否回环 不符")
	}

	随机 := 取随机IP()
	dlog.Info("随机IP", "结果", 随机, "有效", 是否有效(随机), "私有", 是否私有(随机))
	if !是否有效(随机) || 取版本(随机) != "IPv4" || 是否私有(随机) || 是否回环(随机) {
		t.Errorf("随机IP 不符: %s", 随机)
	}

	本机 := 取本机IP()
	dlog.Info("本机IP", "结果", 本机)
	if 本机 != "" && (取版本(本机) != "IPv4" || 是否回环(本机)) {
		t.Errorf("本机IP 不符: %s", 本机)
	}
}

func Test加载与查询(t *testing.T) {
	_, err := 查询("8.8.8.8")
	if err == nil {
		t.Errorf("未加载时应返回错误")
	}

	路径 := os.Getenv("IP2REGION_XDB")
	if 路径 == "" {
		if err := 加载默认数据文件(); err != nil {
			t.Skip("无 IP2REGION_XDB 且默认 xdb 不可用: ", err)
		}
	} else if err := 加载数据文件(路径); err != nil {
		t.Fatalf("加载失败: %v", err)
	}
	地区, err := 查询("8.8.8.8")
	dlog.Info("查询", "ip", "8.8.8.8", "地区", 地区, "错误", err)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if 地区 == "" {
		t.Errorf("查询结果为空")
	}
}

func Test加载默认数据文件(t *testing.T) {
	if _, err := os.Stat(默认xdb相对路径); err != nil {
		t.Skip("仓库 data/ip2region_v4.xdb 不存在")
	}
	if err := 加载默认数据文件(); err != nil {
		t.Fatalf("默认加载失败: %v", err)
	}
	地区, err := 查询("1.1.1.1")
	dlog.Info("默认xdb查询", "地区", 地区, "错误", err)
	if err != nil || 地区 == "" {
		t.Fatalf("查询失败: %q %v", 地区, err)
	}
}
