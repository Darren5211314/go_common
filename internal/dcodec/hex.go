package dcodec

import "encoding/hex"

// 十六进制编码 将文本按 UTF-8 字节编为小写 hex。
func 十六进制编码(文本 string) string {
	return hex.EncodeToString([]byte(文本))
}

// 十六进制解码 解码 hex 为字符串；失败返回空串。
func 十六进制解码(文本 string) string {
	数据, 错误 := hex.DecodeString(文本)
	if 错误 != nil {
		return ""
	}
	return string(数据)
}
