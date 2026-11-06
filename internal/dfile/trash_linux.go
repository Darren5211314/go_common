//go:build linux

package dfile

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// D删除到回收站 按 XDG Trash 规范移入 ~/.local/share/Trash。
func moveToTrash(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("路径为空")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, err := os.Lstat(absPath)
	if err != nil {
		return err
	}

	trashBase, err := trashRoot()
	if err != nil {
		return err
	}
	filesDir := filepath.Join(trashBase, "files")
	infoDir := filepath.Join(trashBase, "info")
	if err := os.MkdirAll(filesDir, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(infoDir, 0o700); err != nil {
		return err
	}

	base := filepath.Base(absPath)
	dst := filepath.Join(filesDir, base)
	infoFile := filepath.Join(infoDir, base+".trashinfo")
	dst, infoFile = avoidConflict(dst, infoFile)

	body := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		escapeTrashPath(absPath),
		time.Now().Format("2006-01-02T15:04:05"),
	)
	if err := os.WriteFile(infoFile, []byte(body), 0o600); err != nil {
		return err
	}
	if err := os.Rename(absPath, dst); err != nil {
		_ = os.Remove(infoFile)
		// 跨盘：复制/移动回退
		if info.IsDir() {
			return err
		}
		readData, readErr := os.ReadFile(absPath)
		if readErr != nil {
			return readErr
		}
		if writeErr := os.WriteFile(dst, readData, info.Mode().Perm()); writeErr != nil {
			_ = os.Remove(infoFile)
			return writeErr
		}
		if removeErr := os.Remove(absPath); removeErr != nil {
			_ = os.Remove(dst)
			_ = os.Remove(infoFile)
			return removeErr
		}
	}
	return nil
}

func trashRoot() (string, error) {
	if home := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); home != "" {
		return filepath.Join(home, "Trash"), nil
	}
	cur, err := user.Current()
	if err != nil {
		home := os.Getenv("HOME")
		if home == "" {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "Trash"), nil
	}
	return filepath.Join(cur.HomeDir, ".local", "share", "Trash"), nil
}

func avoidConflict(dst, infoFile string) (string, string) {
	available := func(filePath, infoPath string) bool {
		_, err1 := os.Lstat(filePath)
		_, err2 := os.Lstat(infoPath)
		return os.IsNotExist(err1) && os.IsNotExist(err2)
	}
	if available(dst, infoFile) {
		return dst, infoFile
	}
	ext := filepath.Ext(dst)
	stem := strings.TrimSuffix(dst, ext)
	infoDir := filepath.Dir(infoFile)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", stem, i, ext)
		candidateInfo := filepath.Join(infoDir, filepath.Base(candidate)+".trashinfo")
		if available(candidate, candidateInfo) {
			return candidate, candidateInfo
		}
	}
}

func escapeTrashPath(path string) string {
	replacer := strings.NewReplacer(
		`\`, `%5C`,
		` `, `%20`,
		`#`, `%23`,
		`?`, `%3F`,
		`&`, `%26`,
		`=`, `%3D`,
		`%`, `%25`,
	)
	return replacer.Replace(path)
}
