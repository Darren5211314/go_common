package dfile

// D文件类 对外命名空间类型（根包通过类型别名导出）。
type D文件类 struct{}

func (D文件类) D删除文件(路径 string) error {
	return 删除文件(路径)
}

func (D文件类) D移动文件(源, 目标 string) error {
	return 移动文件(源, 目标)
}

func (D文件类) D复制文件(源, 目标 string) error {
	return 复制文件(源, 目标)
}

func (D文件类) D重命名(旧路径, 新路径 string) error {
	return 重命名(旧路径, 新路径)
}

func (D文件类) D取扩展名(路径 string) string {
	return 取扩展名(路径)
}

func (D文件类) D文件是否存在(路径 string) bool {
	return 文件是否存在(路径)
}

func (D文件类) D目录是否存在(路径 string) bool {
	return 目录是否存在(路径)
}

func (D文件类) D取文件名(路径 string, 带扩展名 bool) string {
	return 取文件名(路径, 带扩展名)
}

func (D文件类) D保存文本(路径, 分隔符 string, 内容 ...string) error {
	return 保存文本(路径, 分隔符, 内容...)
}

func (D文件类) D取文件大小(路径, 单位 string) (float64, error) {
	return 取文件大小(路径, 单位)
}

func (D文件类) D删除到回收站(路径 string) error {
	return 删除到回收站(路径)
}

