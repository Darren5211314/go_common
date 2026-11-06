// Package dfile 实现文件与目录常用操作。
package dfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Delete 永久删除文件或目录（目录非空时递归删除）。
func deleteFile(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("路径为空")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

// Move 移动文件或目录到目标路径；目标父目录自动创建。
// 同盘用 Rename；跨盘失败则复制后删除源（仅普通文件可跨盘回退）。
func moveFile(src, dst string) error {
	src = strings.TrimSpace(src)
	dst = strings.TrimSpace(dst)
	if src == "" || dst == "" {
		return errors.New("路径为空")
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("跨盘移动目录失败，请先复制再删除")
	}
	if err := copyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// Copy 复制普通文件到目标；目标父目录自动创建；覆盖已存在目标。
func copyFile(src, dst string) error {
	src = strings.TrimSpace(src)
	dst = strings.TrimSpace(dst)
	if src == "" || dst == "" {
		return errors.New("路径为空")
	}
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return errors.New("不支持复制目录")
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, srcInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, inFile); err != nil {
		return err
	}
	return outFile.Sync()
}

// Rename 重命名或移动路径（同 os.Rename）。
func renamePath(oldPath, newPath string) error {
	oldPath = strings.TrimSpace(oldPath)
	newPath = strings.TrimSpace(newPath)
	if oldPath == "" || newPath == "" {
		return errors.New("路径为空")
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

// Ext 返回含点扩展名，如 ".txt"；无则 ""。
func getExt(path string) string {
	return filepath.Ext(path)
}

// Exists 路径存在且为普通文件（非目录）时为 true。
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists 路径存在且为目录时为 true。
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Name 取路径最后一段；带扩展名=false 时去掉扩展名。
func getFileName(path string, withExt bool) string {
	name := filepath.Base(path)
	if withExt {
		return name
	}
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}
