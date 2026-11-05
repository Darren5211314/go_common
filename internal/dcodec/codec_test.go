package dcodec

import (
	"strings"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test网址编解码(t *testing.T) {
	原文 := "a=1&中文 空格"
	编码 := 网址编码(原文)
	解码 := 网址解码(编码)
	dlog.Info("网址", "编码", 编码, "解码", 解码)
	if 解码 != 原文 {
		t.Fatalf("网址往返不符")
	}
	if 网址编码("a b") != "a%20b" {
		t.Fatalf("空格应为 %%20，得到 %q", 网址编码("a b"))
	}
	if !strings.Contains(网址编码("a+b"), "%2B") {
		t.Fatalf("字面 + 应编码为 %%2B，得到 %q", 网址编码("a+b"))
	}
	if 网址解码("a%20b") != "a b" {
		t.Fatalf("%%20 解码应为空格")
	}
	if 网址解码("%ZZ") != "" {
		t.Fatalf("非法网址解码应为空")
	}
}

func TestBase64与网址Base64(t *testing.T) {
	原文 := "hello-中文"
	编码 := Base64编码(原文)
	if Base64解码(编码) != 原文 {
		t.Fatalf("Base64往返不符")
	}
	网址 := Base64网址编码(原文)
	if Base64网址解码(网址) != 原文 {
		t.Fatalf("Base64网址往返不符: %s", 网址)
	}
	if Base64解码("!!!") != "" {
		t.Fatalf("非法 Base64 应为空")
	}
}

func Test十六进制(t *testing.T) {
	原文 := "abc"
	编码 := 十六进制编码(原文)
	if 编码 != "616263" {
		t.Fatalf("hex编码不符: %s", 编码)
	}
	if 十六进制解码(编码) != 原文 {
		t.Fatalf("hex往返不符")
	}
	if 十六进制解码("zz") != "" {
		t.Fatalf("非法 hex 应为空")
	}
}

func TestUTF8与GBK(t *testing.T) {
	if !是否UTF8("你好") {
		t.Fatal("应为 UTF-8")
	}
	if 是否UTF8(string([]byte{0xff, 0xfe})) {
		t.Fatal("非法 UTF-8 应 false")
	}
	gbk, err := 到GBK("你好")
	if err != nil {
		t.Fatal(err)
	}
	utf8, err := 从GBK(gbk)
	dlog.Info("GBK往返", "字节数", len(gbk), "文本", utf8)
	if err != nil || utf8 != "你好" {
		t.Fatalf("GBK往返不符: %q %v", utf8, err)
	}
}
