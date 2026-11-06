package dcodec

// CodecNS 编码/解码命名空间。
type CodecNS struct{}

func (CodecNS) URLEncode(text string) string { return urlEncode(text) }
func (CodecNS) URLDecode(text string) string { return urlDecode(text) }

// QueryParam 取 URL 查询参数；找不到或非法返回空串。
func (CodecNS) QueryParam(rawURL, key string) string {
	return queryParam(rawURL, key)
}
func (CodecNS) Base64Encode(text string) string { return base64Encode(text) }
func (CodecNS) Base64Decode(text string) string { return base64Decode(text) }
func (CodecNS) Base64URLEncode(text string) string {
	return base64URLEncode(text)
}
func (CodecNS) Base64URLDecode(text string) string {
	return base64URLDecode(text)
}
func (CodecNS) HexEncode(text string) string { return hexEncode(text) }
func (CodecNS) HexDecode(text string) string { return hexDecode(text) }
func (CodecNS) IsUTF8(text string) bool     { return isUTF8(text) }
func (CodecNS) ToGBK(text string) ([]byte, error) {
	return toGBK(text)
}
func (CodecNS) FromGBK(data []byte) (string, error) {
	return fromGBK(data)
}
