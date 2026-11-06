package dcodec

import (
	"bytes"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 到GBK 将 UTF-8 文本转为 GBK 字节。
func toGBK(text string) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(text)), simplifiedchinese.GBK.NewEncoder())
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}
	ret := buf.Bytes()
	result := make([]byte, len(ret))
	copy(result, ret)
	return result, nil
}

// 从GBK 将 GBK 字节转为 UTF-8 文本。
func fromGBK(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(reader); err != nil {
		return "", err
	}
	return buf.String(), nil
}
