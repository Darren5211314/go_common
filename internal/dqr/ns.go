package dqr

// QR类 对外命名空间类型（根包通过类型别名导出）。
type QR类 struct{}

func (QR类) D生成(内容 string, 边长 int) (string, error) {
	return 生成(内容, 边长)
}

func (QR类) D生成到文件(内容, 路径 string, 边长 int) error {
	return 生成到文件(内容, 路径, 边长)
}

func (QR类) D解析文件(路径 string) (string, error) {
	return 解析文件(路径)
}

func (QR类) D解析Base64(图片Base64 string) (string, error) {
	return 解析Base64(图片Base64)
}

