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

const defaultSize = 256

func normalizeSize(size int) int {
	if size <= 0 {
		return defaultSize
	}
	return size
}

func genPNG(content string, size int) ([]byte, error) {
	return qrcode.Encode(content, qrcode.Medium, normalizeSize(size))
}

// Generate 生成二维码 PNG，返回 Base64（无 data: 前缀）。
func generate(content string, size int) (string, error) {
	data, err := genPNG(content, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// GenerateFile 生成二维码 PNG 并写入路径。
func genToFile(content, path string, size int) error {
	data, err := genPNG(content, size)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func decodeImage(img image.Image) (string, error) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	result, err := gozxingqr.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}

func parseBytes(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return decodeImage(img)
}

// ParseFile 从图片文件解析二维码文本（支持 PNG/JPEG 等）。
func parseFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return parseBytes(data)
}

// ParseBase64 从 PNG（等）图片的 Base64 解析二维码文本。
func parseBase64(pngBase64 string) (string, error) {
	text := strings.TrimSpace(pngBase64)
	if i := strings.Index(text, ","); i >= 0 && strings.Contains(strings.ToLower(text[:i]), "base64") {
		text = text[i+1:]
	}
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	return parseBytes(data)
}
