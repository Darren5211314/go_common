package dqr

// QR类 对外命名空间类型（根包通过类型别名导出）。
type QR类 struct{}

func (QR类) D生成(content string, size int) (string, error) {
	return generate(content, size)
}

func (QR类) D生成到文件(content, path string, size int) error {
	return genToFile(content, path, size)
}

func (QR类) D解析文件(path string) (string, error) {
	return parseFile(path)
}

func (QR类) D解析Base64(pngBase64 string) (string, error) {
	return parseBase64(pngBase64)
}

