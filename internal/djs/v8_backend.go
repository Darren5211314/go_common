//go:build cgo && !windows

package djs

import (
	"fmt"

	"rogchap.com/v8go"
)

type v8Backend struct {
	isolate  *v8go.Isolate
	ctx *v8go.Context
}

func newV8() (backend, error) {
	isolate := v8go.NewIsolate()
	ctx := v8go.NewContext(isolate)
	return &v8Backend{isolate: isolate, ctx: ctx}, nil
}

func (v *v8Backend) eval(code string) (string, error) {
	value, err := v.ctx.RunScript(code, "djs.js")
	if err != nil {
		return "", err
	}
	if value == nil || value.IsUndefined() || value.IsNull() {
		return "", nil
	}
	return value.String(), nil
}

func (v *v8Backend) setVar(name string, value any) error {
	global := v.ctx.Global()
	jsVal, err := toV8Value(v.isolate, value)
	if err != nil {
		return err
	}
	return global.Set(name, jsVal)
}

func (v *v8Backend) getRaw(name string) (any, bool, error) {
	global := v.ctx.Global()
	value, err := global.Get(name)
	if err != nil {
		return nil, false, err
	}
	if value == nil || value.IsUndefined() {
		return nil, false, nil
	}
	if value.IsBoolean() {
		return value.Boolean(), true, nil
	}
	if value.IsNumber() {
		n := value.Number()
		if n == float64(int64(n)) {
			return int64(n), true, nil
		}
		return n, true, nil
	}
	if value.IsString() {
		return value.String(), true, nil
	}
	if value.IsNull() {
		return nil, true, nil
	}
	return value.String(), true, nil
}

func (v *v8Backend) close() {
	if v.ctx != nil {
		v.ctx.Close()
		v.ctx = nil
	}
	if v.isolate != nil {
		v.isolate.Dispose()
		v.isolate = nil
	}
}

func toV8Value(isolate *v8go.Isolate, value any) (*v8go.Value, error) {
	if value == nil {
		return v8go.Null(isolate), nil
	}
	switch x := value.(type) {
	case string:
		return v8go.NewValue(isolate, x)
	case bool:
		return v8go.NewValue(isolate, x)
	case int64:
		return v8go.NewValue(isolate, x)
	case float64:
		return v8go.NewValue(isolate, x)
	default:
		return v8go.NewValue(isolate, fmt.Sprint(x))
	}
}
