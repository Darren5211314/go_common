package dfile

import (
	"errors"
	"os"
	"strings"
)

// D取文件大小 按单位返回文件大小（1024 进制）。
// 单位：B/KB/MB/GB/TB（大小写不敏感）；空或非法按 B。
func 取文件大小(路径, 单位 string) (float64, error) {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return 0, errors.New("路径为空")
	}
	信息, 错误 := os.Stat(路径)
	if 错误 != nil {
		return 0, 错误
	}
	if 信息.IsDir() {
		return 0, errors.New("路径是目录")
	}
	字节 := float64(信息.Size())
	除数 := 单位除数(单位)
	return 字节 / 除数, nil
}

func 单位除数(单位 string) float64 {
	开关 := strings.ToUpper(strings.TrimSpace(单位))
	开关 = strings.TrimSuffix(开关, "YTE")
	开关 = strings.TrimSuffix(开关, "IB")
	switch 开关 {
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
