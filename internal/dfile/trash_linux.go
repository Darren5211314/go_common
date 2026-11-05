//go:build linux

package dfile

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// D删除到回收站 按 XDG Trash 规范移入 ~/.local/share/Trash。
func 删除到回收站(路径 string) error {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return errors.New("路径为空")
	}
	绝对, 错误 := filepath.Abs(路径)
	if 错误 != nil {
		return 错误
	}
	信息, 错误 := os.Lstat(绝对)
	if 错误 != nil {
		return 错误
	}

	回收根, 错误 := 取回收站根()
	if 错误 != nil {
		return 错误
	}
	文件目录 := filepath.Join(回收根, "files")
	信息目录 := filepath.Join(回收根, "info")
	if 错误 := os.MkdirAll(文件目录, 0o700); 错误 != nil {
		return 错误
	}
	if 错误 := os.MkdirAll(信息目录, 0o700); 错误 != nil {
		return 错误
	}

	基名 := filepath.Base(绝对)
	目标 := filepath.Join(文件目录, 基名)
	信息文件 := filepath.Join(信息目录, 基名+".trashinfo")
	目标, 信息文件 = 避免冲突(目标, 信息文件)

	正文 := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		转义Trash路径(绝对),
		time.Now().Format("2006-01-02T15:04:05"),
	)
	if 错误 := os.WriteFile(信息文件, []byte(正文), 0o600); 错误 != nil {
		return 错误
	}
	if 错误 := os.Rename(绝对, 目标); 错误 != nil {
		_ = os.Remove(信息文件)
		// 跨盘：复制/移动回退
		if 信息.IsDir() {
			return 错误
		}
		读数据, 读错 := os.ReadFile(绝对)
		if 读错 != nil {
			return 读错
		}
		if 写错 := os.WriteFile(目标, 读数据, 信息.Mode().Perm()); 写错 != nil {
			_ = os.Remove(信息文件)
			return 写错
		}
		if 删错 := os.Remove(绝对); 删错 != nil {
			_ = os.Remove(目标)
			_ = os.Remove(信息文件)
			return 删错
		}
	}
	return nil
}

func 取回收站根() (string, error) {
	if 家 := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); 家 != "" {
		return filepath.Join(家, "Trash"), nil
	}
	当前, 错误 := user.Current()
	if 错误 != nil {
		家 := os.Getenv("HOME")
		if 家 == "" {
			return "", 错误
		}
		return filepath.Join(家, ".local", "share", "Trash"), nil
	}
	return filepath.Join(当前.HomeDir, ".local", "share", "Trash"), nil
}

func 避免冲突(目标, 信息文件 string) (string, string) {
	可用 := func(文件路径, 信息路径 string) bool {
		_, 错1 := os.Lstat(文件路径)
		_, 错2 := os.Lstat(信息路径)
		return os.IsNotExist(错1) && os.IsNotExist(错2)
	}
	if 可用(目标, 信息文件) {
		return 目标, 信息文件
	}
	扩展 := filepath.Ext(目标)
	主干 := strings.TrimSuffix(目标, 扩展)
	信息目录 := filepath.Dir(信息文件)
	for i := 1; ; i++ {
		候选 := fmt.Sprintf("%s_%d%s", 主干, i, 扩展)
		候选信息 := filepath.Join(信息目录, filepath.Base(候选)+".trashinfo")
		if 可用(候选, 候选信息) {
			return 候选, 候选信息
		}
	}
}

func 转义Trash路径(路径 string) string {
	替换 := strings.NewReplacer(
		`\`, `%5C`,
		` `, `%20`,
		`#`, `%23`,
		`?`, `%3F`,
		`&`, `%26`,
		`=`, `%3D`,
		`%`, `%25`,
	)
	return 替换.Replace(路径)
}
