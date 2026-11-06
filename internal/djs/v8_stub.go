//go:build !cgo

package djs

import "errors"

func newV8() (backend, error) {
	return nil, errors.New("当前构建未启用 CGO，无法使用 V8；Linux 请用 CGO_ENABLED=1，Windows 请改用 goja（v8go 无 Windows 预编译库）")
}
