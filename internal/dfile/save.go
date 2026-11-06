package dfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const defaultTextSep = "----"

// D保存文本 将若干内容用分隔符拼接后覆盖写入文件（UTF-8）。
// 分隔符为空时使用默认 "----"；自动创建父目录。
func saveText(path, sep string, content ...string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("路径为空")
	}
	if sep == "" {
		sep = defaultTextSep
	}
	body := strings.Join(content, sep)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}
