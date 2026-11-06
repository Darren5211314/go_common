package dfile

import (
	"errors"
	"os"
	"strings"
)

// D取文件大小 按单位返回文件大小（1024 进制）。
// 单位：B/KB/MB/GB/TB（大小写不敏感）；空或非法按 B。
func getFileSize(path, unit string) (float64, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return 0, errors.New("路径为空")
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 0, errors.New("路径是目录")
	}
	bytes := float64(info.Size())
	div := unitDivisor(unit)
	return bytes / div, nil
}

func unitDivisor(unit string) float64 {
	u := strings.ToUpper(strings.TrimSpace(unit))
	u = strings.TrimSuffix(u, "YTE")
	u = strings.TrimSuffix(u, "IB")
	switch u {
	case "K", "KB":
		return 1024
	case "M", "MB":
		return 1024 * 1024
	case "G", "GB":
		return 1024 * 1024 * 1024
	case "T", "TB":
		return 1024 * 1024 * 1024 * 1024
	default:
		return 1 // B 或非法
	}
}
