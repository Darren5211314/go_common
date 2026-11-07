package dfile

// FileNS 命名空间类型（根包以未导出别名挂到 var File）。
type FileNS struct{}

func (FileNS) Delete(path string) error {
	return deleteFile(path)
}

func (FileNS) Move(src, dst string) error {
	return moveFile(src, dst)
}

func (FileNS) Copy(src, dst string) error {
	return copyFile(src, dst)
}

func (FileNS) Rename(oldPath, newPath string) error {
	return renamePath(oldPath, newPath)
}

func (FileNS) Ext(path string) string {
	return getExt(path)
}

func (FileNS) Exists(path string) bool {
	return fileExists(path)
}

func (FileNS) DirExists(path string) bool {
	return dirExists(path)
}

func (FileNS) Mkdir(path string) error {
	return mkdir(path)
}

func (FileNS) ExeDir() string {
	return exeDir()
}

func (FileNS) Name(path string, withExt bool) string {
	return getFileName(path, withExt)
}

func (FileNS) SaveText(path, sep string, overwrite bool, content ...string) error {
	return saveText(path, sep, overwrite, content...)
}

func (FileNS) Size(path, unit string) (float64, error) {
	return getFileSize(path, unit)
}

func (FileNS) Trash(path string) error {
	return moveToTrash(path)
}

