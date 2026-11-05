// Package dfile 实现文件与目录常用操作。
package dfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// D删除文件 永久删除文件或目录（目录非空时递归删除）。
func 删除文件(路径 string) error {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return errors.New("路径为空")
	}
	信息, 错误 := os.Lstat(路径)
	if 错误 != nil {
		return 错误
	}
	if 信息.IsDir() {
		return os.RemoveAll(路径)
	}
	return os.Remove(路径)
}

// D移动文件 移动文件或目录到目标路径；目标父目录自动创建。
// 同盘用 Rename；跨盘失败则复制后删除源（仅普通文件可跨盘回退）。
func 移动文件(源, 目标 string) error {
	源 = strings.TrimSpace(源)
	目标 = strings.TrimSpace(目标)
	if 源 == "" || 目标 == "" {
		return errors.New("路径为空")
	}
	if 错误 := os.MkdirAll(filepath.Dir(目标), 0o755); 错误 != nil {
		return 错误
	}
	if 错误 := os.Rename(源, 目标); 错误 == nil {
		return nil
	}
	信息, 错误 := os.Stat(源)
	if 错误 != nil {
		return 错误
	}
	if 信息.IsDir() {
		return errors.New("跨盘移动目录失败，请先复制再删除")
	}
	if 错误 := 复制文件(源, 目标); 错误 != nil {
		return 错误
	}
	return os.Remove(源)
}

// D复制文件 复制普通文件到目标；目标父目录自动创建；覆盖已存在目标。
func 复制文件(源, 目标 string) error {
	源 = strings.TrimSpace(源)
	目标 = strings.TrimSpace(目标)
	if 源 == "" || 目标 == "" {
		return errors.New("路径为空")
	}
	源信息, 错误 := os.Stat(源)
	if 错误 != nil {
		return 错误
	}
	if 源信息.IsDir() {
		return errors.New("不支持复制目录")
	}
	if 错误 := os.MkdirAll(filepath.Dir(目标), 0o755); 错误 != nil {
		return 错误
	}
	读, 错误 := os.Open(源)
	if 错误 != nil {
		return 错误
	}
	defer 读.Close()

	写, 错误 := os.OpenFile(目标, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 源信息.Mode().Perm())
	if 错误 != nil {
		return 错误
	}
	defer 写.Close()

	if _, 错误 = io.Copy(写, 读); 错误 != nil {
		return 错误
	}
	return 写.Sync()
}

// D重命名 重命名或移动路径（同 os.Rename）。
func 重命名(旧路径, 新路径 string) error {
	旧路径 = strings.TrimSpace(旧路径)
	新路径 = strings.TrimSpace(新路径)
	if 旧路径 == "" || 新路径 == "" {
		return errors.New("路径为空")
	}
	if 错误 := os.MkdirAll(filepath.Dir(新路径), 0o755); 错误 != nil {
		return 错误
	}
	return os.Rename(旧路径, 新路径)
}

// D取扩展名 返回含点扩展名，如 ".txt"；无则 ""。
func 取扩展名(路径 string) string {
	return filepath.Ext(路径)
}

// D文件是否存在 路径存在且为普通文件（非目录）时为 true。
func 文件是否存在(路径 string) bool {
	信息, 错误 := os.Stat(路径)
	if 错误 != nil {
		return false
	}
	return !信息.IsDir()
}

// D目录是否存在 路径存在且为目录时为 true。
func 目录是否存在(路径 string) bool {
	信息, 错误 := os.Stat(路径)
	if 错误 != nil {
		return false
	}
	return 信息.IsDir()
}

// D取文件名 取路径最后一段；带扩展名=false 时去掉扩展名。
func 取文件名(路径 string, 带扩展名 bool) string {
	名 := filepath.Base(路径)
	if 带扩展名 {
		return 名
	}
	扩展 := filepath.Ext(名)
	if 扩展 == "" {
		return 名
	}
	return strings.TrimSuffix(名, 扩展)
}
