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
func 删除到回收站(路径 string) error {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return errors.New("路径为空")
	}
	绝对, 错误 := filepath.Abs(路径)
	if 错误 != nil {
		return 错误
	}
	if _, 错误 = os.Lstat(绝对); 错误 != nil {
		return 错误
	}

	// SHFileOperation 要求双 NUL 结尾的路径列表
	来源切片, 错误 := syscall.UTF16FromString(绝对)
	if 错误 != nil {
		return 错误
	}
	来源切片 = append(来源切片, 0) // 额外 NUL，形成双终止
	操作 := &shFileOpStruct{
		wFunc:  foDelete,
		pFrom:  &来源切片[0],
		fFlags: fofAllowUndo | fofNoConfirmation | fofSilent | fofNoErrorUI,
	}
	返回, _, _ := procSHFileOperationW.Call(uintptr(unsafe.Pointer(操作)))
	if 返回 != 0 {
		return errors.New("移入回收站失败")
	}
	if 操作.fAnyOperationsAborted != 0 {
		return errors.New("移入回收站被中止")
	}
	return nil
}
