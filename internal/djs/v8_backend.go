//go:build cgo && !windows

package djs

import (
	"fmt"

	"rogchap.com/v8go"
)

type v8后端 struct {
	隔离  *v8go.Isolate
	上下文 *v8go.Context
}

func 新建V8() (后端, error) {
	隔离 := v8go.NewIsolate()
	上下文 := v8go.NewContext(隔离)
	return &v8后端{隔离: 隔离, 上下文: 上下文}, nil
}

func (v *v8后端) 执行(代码 string) (string, error) {
	值, 错误 := v.上下文.RunScript(代码, "djs.js")
	if 错误 != nil {
		return "", 错误
	}
	if 值 == nil || 值.IsUndefined() || 值.IsNull() {
		return "", nil
	}
	return 值.String(), nil
}

func (v *v8后端) 设置变量(名 string, 值 any) error {
	全局 := v.上下文.Global()
	js值, 错误 := 转V8值(v.隔离, 值)
	if 错误 != nil {
		return 错误
	}
	return 全局.Set(名, js值)
}

func (v *v8后端) 取原始(名 string) (any, bool, error) {
	全局 := v.上下文.Global()
	值, 错误 := 全局.Get(名)
	if 错误 != nil {
		return nil, false, 错误
	}
	if 值 == nil || 值.IsUndefined() {
		return nil, false, nil
	}
	if 值.IsBoolean() {
		return 值.Boolean(), true, nil
	}
	if 值.IsNumber() {
		n := 值.Number()
		if n == float64(int64(n)) {
			return int64(n), true, nil
		}
		return n, true, nil
	}
	if 值.IsString() {
		return 值.String(), true, nil
	}
	if 值.IsNull() {
		return nil, true, nil
	}
	return 值.String(), true, nil
}

func (v *v8后端) 关闭() {
	if v.上下文 != nil {
		v.上下文.Close()
		v.上下文 = nil
	}
	if v.隔离 != nil {
		v.隔离.Dispose()
		v.隔离 = nil
	}
}

func 转V8值(隔离 *v8go.Isolate, 值 any) (*v8go.Value, error) {
	if 值 == nil {
		return v8go.Null(隔离), nil
	}
	switch x := 值.(type) {
	case string:
		return v8go.NewValue(隔离, x)
	case bool:
		return v8go.NewValue(隔离, x)
	case int64:
		return v8go.NewValue(隔离, x)
	case float64:
		return v8go.NewValue(隔离, x)
	default:
		return v8go.NewValue(隔离, fmt.Sprint(x))
	}
}
