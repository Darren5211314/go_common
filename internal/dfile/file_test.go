package dfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test保存文本与读写辅助(t *testing.T) {
	目录 := t.TempDir()
	路径 := filepath.Join(目录, "子", "a.txt")

	if 错误 := 保存文本(路径, "", "甲", "乙", "丙"); 错误 != nil {
		t.Fatalf("保存失败: %v", 错误)
	}
	数据, 错误 := os.ReadFile(路径)
	if 错误 != nil {
		t.Fatal(错误)
	}
	期望 := "甲----乙----丙"
	if string(数据) != 期望 {
		t.Errorf("内容不符: %q", 数据)
	}

	if 错误 := 保存文本(路径, "|", "1", "2"); 错误 != nil {
		t.Fatal(错误)
	}
	数据, _ = os.ReadFile(路径)
	if string(数据) != "1|2" {
		t.Errorf("自定义分隔符不符: %q", 数据)
	}
}

func Test复制移动重命名删除(t *testing.T) {
	目录 := t.TempDir()
	源 := filepath.Join(目录, "源.txt")
	if 错误 := os.WriteFile(源, []byte("hello"), 0o644); 错误 != nil {
		t.Fatal(错误)
	}

	复制目标 := filepath.Join(目录, "子", "复制.txt")
	if 错误 := 复制文件(源, 复制目标); 错误 != nil {
		t.Fatalf("复制失败: %v", 错误)
	}
	if !文件是否存在(复制目标) {
		t.Fatal("复制目标应存在")
	}

	移动目标 := filepath.Join(目录, "移动.txt")
	if 错误 := 移动文件(复制目标, 移动目标); 错误 != nil {
		t.Fatalf("移动失败: %v", 错误)
	}
	if 文件是否存在(复制目标) || !文件是否存在(移动目标) {
		t.Fatal("移动后路径不符")
	}

	新名 := filepath.Join(目录, "新名.txt")
	if 错误 := 重命名(移动目标, 新名); 错误 != nil {
		t.Fatalf("重命名失败: %v", 错误)
	}
	if !文件是否存在(新名) {
		t.Fatal("重命名后应存在")
	}

	if 错误 := 删除文件(新名); 错误 != nil {
		t.Fatalf("删除失败: %v", 错误)
	}
	if 文件是否存在(新名) {
		t.Fatal("删除后不应存在")
	}
}

func Test扩展名文件名存在大小(t *testing.T) {
	目录 := t.TempDir()
	路径 := filepath.Join(目录, "demo.TXT")
	内容 := strings.Repeat("a", 2048)
	if 错误 := os.WriteFile(路径, []byte(内容), 0o644); 错误 != nil {
		t.Fatal(错误)
	}

	if 取扩展名(路径) != ".TXT" {
		t.Errorf("扩展名不符: %q", 取扩展名(路径))
	}
	if 取文件名(路径, true) != "demo.TXT" {
		t.Errorf("带扩展文件名不符")
	}
	if 取文件名(路径, false) != "demo" {
		t.Errorf("不带扩展文件名不符: %q", 取文件名(路径, false))
	}
	if !文件是否存在(路径) || 目录是否存在(路径) {
		t.Errorf("文件存在判断不符")
	}
	if !目录是否存在(目录) || 文件是否存在(目录) {
		t.Errorf("目录存在判断不符")
	}

	字节, 错误 := 取文件大小(路径, "B")
	if 错误 != nil || 字节 != 2048 {
		t.Errorf("B 大小不符: %v %v", 字节, 错误)
	}
	kb, 错误 := 取文件大小(路径, "kb")
	if 错误 != nil || kb != 2 {
		t.Errorf("KB 大小不符: %v %v", kb, 错误)
	}
	dlog.Info("文件大小", "B", 字节, "KB", kb)
}

func Test删除到回收站(t *testing.T) {
	目录 := t.TempDir()
	路径 := filepath.Join(目录, "trash_me.txt")
	if 错误 := os.WriteFile(路径, []byte("x"), 0o644); 错误 != nil {
		t.Fatal(错误)
	}
	if 错误 := 删除到回收站(路径); 错误 != nil {
		t.Fatalf("回收站失败: %v", 错误)
	}
	if 文件是否存在(路径) {
		t.Fatal("回收站后原路径不应存在")
	}
}
