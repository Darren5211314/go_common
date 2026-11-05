//go:build !windows && !linux

package dfile

import "errors"

// D删除到回收站 当前平台不支持。
func 删除到回收站(路径 string) error {
	_ = 路径
	return errors.New("当前平台不支持删除到回收站")
}
