package dqr

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	gozxingqr "github.com/makiuchi-d/gozxing/qrcode"
	qrcode "github.com/skip2/go-qrcode"
)

const 默认边长 = 256

func 规范边长(边长 int) int {
	if 边长 <= 0 {
		return 默认边长
	}
	return 边长
}

func 生成PNG(内容 string, 边长 int) ([]byte, error) {
	return qrcode.Encode(内容, qrcode.Medium, 规范边长(边长))
}

// D生成 生成二维码 PNG，返回 Base64（无 data: 前缀）。
func 生成(内容 string, 边长 int) (string, error) {
	数据, 错误 := 生成PNG(内容, 边长)
	if 错误 != nil {
		return "", 错误
	}
	return base64.StdEncoding.EncodeToString(数据), nil
}

// D生成到文件 生成二维码 PNG 并写入路径。
func 生成到文件(内容, 路径 string, 边长 int) error {
	数据, 错误 := 生成PNG(内容, 边长)
	if 错误 != nil {
		return 错误
	}
	return os.WriteFile(路径, 数据, 0o644)
}

func 解析图片(图 image.Image) (string, error) {
	位图, 错误 := gozxing.NewBinaryBitmapFromImage(图)
	if 错误 != nil {
		return "", 错误
	}
	结果, 错误 := gozxingqr.NewQRCodeReader().Decode(位图, nil)
	if 错误 != nil {
		return "", 错误
	}
	return 结果.GetText(), nil
}

func 解析字节(数据 []byte) (string, error) {
	图, _, 错误 := image.Decode(bytes.NewReader(数据))
	if 错误 != nil {
		return "", 错误
	}
	return 解析图片(图)
}

// D解析文件 从图片文件解析二维码文本（支持 PNG/JPEG 等）。
func 解析文件(路径 string) (string, error) {
	数据, 错误 := os.ReadFile(路径)
	if 错误 != nil {
		return "", 错误
	}
	return 解析字节(数据)
}

// D解析Base64 从 PNG（等）图片的 Base64 解析二维码文本。
func 解析Base64(图片Base64 string) (string, error) {
	文本 := strings.TrimSpace(图片Base64)
	if i := strings.Index(文本, ","); i >= 0 && strings.Contains(strings.ToLower(文本[:i]), "base64") {
		文本 = 文本[i+1:]
	}
	数据, 错误 := base64.StdEncoding.DecodeString(文本)
	if 错误 != nil {
		return "", 错误
	}
	return 解析字节(数据)
}
