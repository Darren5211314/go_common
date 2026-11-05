package dqr

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test生成解析往返(t *testing.T) {
	原文 := "https://example.com/qr-test"
	b64, err := 生成(原文, 256)
	dlog.Info("QR生成", "长度", len(b64), "错误", err)
	if err != nil || b64 == "" {
		t.Fatalf("生成失败: %v", err)
	}

	解出, err := 解析Base64(b64)
	dlog.Info("QR解析Base64", "结果", 解出, "错误", err)
	if err != nil || 解出 != 原文 {
		t.Errorf("Base64往返不符: got %q", 解出)
	}

	中文 := "你好二维码"
	b64, err = 生成(中文, 0) // 默认边长
	解出, err2 := 解析Base64(b64)
	dlog.Info("QR中文往返", "结果", 解出, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解出 != 中文 {
		t.Errorf("中文往返不符: got %q", 解出)
	}

	路径 := filepath.Join(t.TempDir(), "out.png")
	if err := 生成到文件(原文, 路径, 128); err != nil {
		t.Fatalf("生成到文件失败: %v", err)
	}
	解出, err = 解析文件(路径)
	dlog.Info("QR解析文件", "路径", 路径, "结果", 解出, "错误", err)
	if err != nil || 解出 != 原文 {
		t.Errorf("文件往返不符: got %q", 解出)
	}
	信息, _ := os.Stat(路径)
	dlog.Info("文件大小", "字节", 信息.Size())
}
