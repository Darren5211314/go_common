package dcrypto

import (
	"encoding/hex"
	"errors"
	"strings"
)

var (
	错误密钥长度 = errors.New("密钥长度不正确")
	错误IV长度  = errors.New("IV/Nonce长度不正确")
	错误填充    = errors.New("PKCS7填充无效")
	错误密文    = errors.New("密文无效")
	错误数据    = errors.New("数据无效")
)

func 解码十六进制(文本 string) ([]byte, error) {
	文本 = strings.TrimSpace(文本)
	文本 = strings.ReplaceAll(文本, " ", "")
	数据, 错误 := hex.DecodeString(文本)
	if 错误 != nil {
		return nil, 错误
	}
	return 数据, nil
}

func 编码十六进制(数据 []byte) string {
	return hex.EncodeToString(数据)
}

func PKCS7填充(数据 []byte, 块大小 int) []byte {
	填充长 := 块大小 - len(数据)%块大小
	填充 := make([]byte, 填充长)
	for i := range 填充 {
		填充[i] = byte(填充长)
	}
	return append(数据, 填充...)
}

func PKCS7去填充(数据 []byte, 块大小 int) ([]byte, error) {
	长度 := len(数据)
	if 长度 == 0 || 长度%块大小 != 0 {
		return nil, 错误填充
	}
	填充长 := int(数据[长度-1])
	if 填充长 == 0 || 填充长 > 块大小 || 填充长 > 长度 {
		return nil, 错误填充
	}
	for i := 0; i < 填充长; i++ {
		if 数据[长度-1-i] != byte(填充长) {
			return nil, 错误填充
		}
	}
	return 数据[:长度-填充长], nil
}
