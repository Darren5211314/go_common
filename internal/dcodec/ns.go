package dcodec

// CodecNS 编码/解码命名空间。
type CodecNS struct{}

func (CodecNS) D网址编码(text string) string   { return urlEncode(text) }
func (CodecNS) D网址解码(text string) string   { return urlDecode(text) }
func (CodecNS) DBase64编码(text string) string { return base64Encode(text) }
func (CodecNS) DBase64解码(text string) string { return base64Decode(text) }
func (CodecNS) DBase64网址编码(text string) string {
	return base64URLEncode(text)
}
func (CodecNS) DBase64网址解码(text string) string {
	return base64URLDecode(text)
}
func (CodecNS) D十六进制编码(text string) string { return hexEncode(text) }
func (CodecNS) D十六进制解码(text string) string { return hexDecode(text) }
func (CodecNS) D是否UTF8(text string) bool     { return isUTF8(text) }
func (CodecNS) D到GBK(text string) ([]byte, error) {
	return toGBK(text)
}
func (CodecNS) D从GBK(data []byte) (string, error) {
	return fromGBK(data)
}
