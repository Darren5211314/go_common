package dcodec

import "encoding/base64"

// Base64编码 标准 Base64 编码。
func Base64编码(文本 string) string {
	return base64.StdEncoding.EncodeToString([]byte(文本))
}

// Base64解码 标准 Base64 解码；失败返回空串。
func Base64解码(文本 string) string {
	数据, 错误 := base64.StdEncoding.DecodeString(文本)
	if 错误 != nil {
		return ""
	}
	return string(数据)
}

// Base64网址编码 URL 安全、无填充的 Base64（RawURLEncoding）。
func Base64网址编码(文本 string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(文本))
}

// Base64网址解码 解码 RawURLEncoding；失败返回空串。
func Base64网址解码(文本 string) string {
	数据, 错误 := base64.RawURLEncoding.DecodeString(文本)
	if 错误 != nil {
		return ""
	}
	return string(数据)
}
