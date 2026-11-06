//go:build windows

package dfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	shell32              = syscall.NewLazyDLL("shell32.dll")
	procSHFileOperationW = shell32.NewProc("SHFileOperationW")
)

const (
	foDelete          = 0x0003
	fofAllowUndo      = 0x0040
	fofNoConfirmation = 0x0010
	fofSilent         = 0x0004
	fofNoErrorUI      = 0x0400
)

type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

// D删除到回收站 将文件或目录移入 Windows 回收站。
func moveToTrash(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("路径为空")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err = os.Lstat(absPath); err != nil {
		return err
	}

	// SHFileOperation 要求双 NUL 结尾的路径列表
	fromUTF16, err := syscall.UTF16FromString(absPath)
	if err != nil {
		return err
	}
	fromUTF16 = append(fromUTF16, 0) // 额外 NUL，形成双终止
	op := &shFileOpStruct{
		wFunc:  foDelete,
		pFrom:  &fromUTF16[0],
		fFlags: fofAllowUndo | fofNoConfirmation | fofSilent | fofNoErrorUI,
	}
	ret, _, _ := procSHFileOperationW.Call(uintptr(unsafe.Pointer(op)))
	if ret != 0 {
		return errors.New("移入回收站失败")
	}
	if op.fAnyOperationsAborted != 0 {
		return errors.New("移入回收站被中止")
	}
	return nil
}
