package dcodec

import "encoding/hex"

// 十六进制编码 将文本按 UTF-8 字节编为小写 hex。
func hexEncode(text string) string {
	return hex.EncodeToString([]byte(text))
}

// 十六进制解码 解码 hex 为字符串；失败返回空串。
func hexDecode(text string) string {
	data, err := hex.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(data)
}
