//go:build !windows && !linux

package dfile

import "errors"

// D删除到回收站 当前平台不支持。
func moveToTrash(path string) error {
	_ = path
	return errors.New("当前平台不支持删除到回收站")
}
