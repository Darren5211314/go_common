package dfile

// FileNS 命名空间类型（根包以未导出别名挂到 var File）。
type FileNS struct{}

func (FileNS) D删除文件(path string) error {
	return deleteFile(path)
}

func (FileNS) D移动文件(src, dst string) error {
	return moveFile(src, dst)
}

func (FileNS) D复制文件(src, dst string) error {
	return copyFile(src, dst)
}

func (FileNS) D重命名(oldPath, newPath string) error {
	return renamePath(oldPath, newPath)
}

func (FileNS) D取扩展名(path string) string {
	return getExt(path)
}

func (FileNS) D文件是否存在(path string) bool {
	return fileExists(path)
}

func (FileNS) D目录是否存在(path string) bool {
	return dirExists(path)
}

func (FileNS) D取文件名(path string, withExt bool) string {
	return getFileName(path, withExt)
}

func (FileNS) D保存文本(path, sep string, content ...string) error {
	return saveText(path, sep, content...)
}

func (FileNS) D取文件大小(path, unit string) (float64, error) {
	return getFileSize(path, unit)
}

func (FileNS) D删除到回收站(path string) error {
	return moveToTrash(path)
}

