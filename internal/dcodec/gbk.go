package dcodec

import (
	"bytes"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 到GBK 将 UTF-8 文本转为 GBK 字节。
func 到GBK(文本 string) ([]byte, error) {
	阅读器 := transform.NewReader(bytes.NewReader([]byte(文本)), simplifiedchinese.GBK.NewEncoder())
	缓冲 := &bytes.Buffer{}
	if _, 错误 := 缓冲.ReadFrom(阅读器); 错误 != nil {
		return nil, 错误
	}
	返回 := 缓冲.Bytes()
	结果 := make([]byte, len(返回))
	copy(结果, 返回)
	return 结果, nil
}

// 从GBK 将 GBK 字节转为 UTF-8 文本。
func 从GBK(数据 []byte) (string, error) {
	阅读器 := transform.NewReader(bytes.NewReader(数据), simplifiedchinese.GBK.NewDecoder())
	缓冲 := &bytes.Buffer{}
	if _, 错误 := 缓冲.ReadFrom(阅读器); 错误 != nil {
		return "", 错误
	}
	return 缓冲.String(), nil
}
