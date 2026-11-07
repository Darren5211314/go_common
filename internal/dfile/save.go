package dfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const defaultTextSep = "----"

// SaveText 将若干内容用分隔符拼接后写入文件（UTF-8）。
// 分隔符为空时使用默认 "----"；自动创建父目录。
// overwrite=false 追加（若文件已有内容且本次 body 非空，先写分隔符再写 body）；
// overwrite=true 覆盖。
func saveText(path, sep string, overwrite bool, content ...string) error {
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
	if overwrite {
		return os.WriteFile(path, []byte(body), 0o644)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 0 && body != "" {
		if _, err = f.WriteString(sep); err != nil {
			return err
		}
	}
	_, err = f.WriteString(body)
	return err
}
