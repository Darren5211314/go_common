package dfile

// D文件类 对外命名空间类型（根包通过类型别名导出）。
type D文件类 struct{}

func (D文件类) D删除文件(path string) error {
	return deleteFile(path)
}

func (D文件类) D移动文件(src, dst string) error {
	return moveFile(src, dst)
}

func (D文件类) D复制文件(src, dst string) error {
	return copyFile(src, dst)
}

func (D文件类) D重命名(oldPath, newPath string) error {
	return renamePath(oldPath, newPath)
}

func (D文件类) D取扩展名(path string) string {
	return getExt(path)
}

func (D文件类) D文件是否存在(path string) bool {
	return fileExists(path)
}

func (D文件类) D目录是否存在(path string) bool {
	return dirExists(path)
}

func (D文件类) D取文件名(path string, withExt bool) string {
	return getFileName(path, withExt)
}

func (D文件类) D保存文本(path, sep string, content ...string) error {
	return saveText(path, sep, content...)
}

func (D文件类) D取文件大小(path, unit string) (float64, error) {
	return getFileSize(path, unit)
}

func (D文件类) D删除到回收站(path string) error {
	return moveToTrash(path)
}

