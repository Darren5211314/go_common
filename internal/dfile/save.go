package dfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const 默认文本分隔符 = "----"

// D保存文本 将若干内容用分隔符拼接后覆盖写入文件（UTF-8）。
// 分隔符为空时使用默认 "----"；自动创建父目录。
func 保存文本(路径, 分隔符 string, 内容 ...string) error {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return errors.New("路径为空")
	}
	if 分隔符 == "" {
		分隔符 = 默认文本分隔符
	}
	正文 := strings.Join(内容, 分隔符)
	if 错误 := os.MkdirAll(filepath.Dir(路径), 0o755); 错误 != nil {
		return 错误
	}
	return os.WriteFile(路径, []byte(正文), 0o644)
}
