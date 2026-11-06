package dqr

// QRNS 对外命名空间类型（根包通过类型别名导出）。
type QRNS struct{}

func (QRNS) Generate(content string, size int) (string, error) {
	return generate(content, size)
}

func (QRNS) GenerateFile(content, path string, size int) error {
	return genToFile(content, path, size)
}

func (QRNS) ParseFile(path string) (string, error) {
	return parseFile(path)
}

func (QRNS) ParseBase64(pngBase64 string) (string, error) {
	return parseBase64(pngBase64)
}

