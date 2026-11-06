package dqr

// QRNS 对外命名空间类型（根包通过类型别名导出）。
type QRNS struct{}

func (QRNS) D生成(content string, size int) (string, error) {
	return generate(content, size)
}

func (QRNS) D生成到文件(content, path string, size int) error {
	return genToFile(content, path, size)
}

func (QRNS) D解析文件(path string) (string, error) {
	return parseFile(path)
}

func (QRNS) D解析Base64(pngBase64 string) (string, error) {
	return parseBase64(pngBase64)
}

