package dcodec

// Codec类 编码/解码命名空间。
type Codec类 struct{}

func (Codec类) D网址编码(文本 string) string   { return 网址编码(文本) }
func (Codec类) D网址解码(文本 string) string   { return 网址解码(文本) }
func (Codec类) DBase64编码(文本 string) string { return Base64编码(文本) }
func (Codec类) DBase64解码(文本 string) string { return Base64解码(文本) }
func (Codec类) DBase64网址编码(文本 string) string {
	return Base64网址编码(文本)
}
func (Codec类) DBase64网址解码(文本 string) string {
	return Base64网址解码(文本)
}
func (Codec类) D十六进制编码(文本 string) string { return 十六进制编码(文本) }
func (Codec类) D十六进制解码(文本 string) string { return 十六进制解码(文本) }
func (Codec类) D是否UTF8(文本 string) bool     { return 是否UTF8(文本) }
func (Codec类) D到GBK(文本 string) ([]byte, error) {
	return 到GBK(文本)
}
func (Codec类) D从GBK(数据 []byte) (string, error) {
	return 从GBK(数据)
}
