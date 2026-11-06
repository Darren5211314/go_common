package dcodec

import "encoding/base64"

// Base64编码 标准 Base64 编码。
func base64Encode(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// Base64解码 标准 Base64 解码；失败返回空串。
func base64Decode(text string) string {
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(data)
}

// Base64网址编码 URL 安全、无填充的 Base64（RawURLEncoding）。
func base64URLEncode(text string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(text))
}

// Base64网址解码 解码 RawURLEncoding；失败返回空串。
func base64URLDecode(text string) string {
	data, err := base64.RawURLEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(data)
}
