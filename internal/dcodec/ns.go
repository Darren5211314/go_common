package dcodec

// Codec类 编码/解码命名空间。
type Codec类 struct{}

func (Codec类) D网址编码(text string) string   { return urlEncode(text) }
func (Codec类) D网址解码(text string) string   { return urlDecode(text) }
func (Codec类) DBase64编码(text string) string { return base64Encode(text) }
func (Codec类) DBase64解码(text string) string { return base64Decode(text) }
func (Codec类) DBase64网址编码(text string) string {
	return base64URLEncode(text)
}
func (Codec类) DBase64网址解码(text string) string {
	return base64URLDecode(text)
}
func (Codec类) D十六进制编码(text string) string { return hexEncode(text) }
func (Codec类) D十六进制解码(text string) string { return hexDecode(text) }
func (Codec类) D是否UTF8(text string) bool     { return isUTF8(text) }
func (Codec类) D到GBK(text string) ([]byte, error) {
	return toGBK(text)
}
func (Codec类) D从GBK(data []byte) (string, error) {
	return fromGBK(data)
}
