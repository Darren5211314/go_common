//go:build cgo && windows

package djs

import "errors"

// rogchap.com/v8go 仅附带 darwin/linux 预编译 libv8，Windows 无法直接链接。
func newV8() (backend, error) {
	return nil, errors.New("当前 v8go 未提供 Windows 预编译 V8 库；Windows 请用 goja，Linux 可用 CGO_ENABLED=1 启用 v8")
}
