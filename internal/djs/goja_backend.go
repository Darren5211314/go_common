package djs

import (
	"fmt"

	"github.com/dop251/goja"
)

type gojaBackend struct {
	runtime *goja.Runtime
}

func newGoja() (backend, error) {
	return &gojaBackend{runtime: goja.New()}, nil
}

func (g *gojaBackend) eval(code string) (string, error) {
	value, err := g.runtime.RunString(code)
	if err != nil {
		return "", err
	}
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return "", nil
	}
	return value.String(), nil
}

func (g *gojaBackend) setVar(name string, value any) error {
	return g.runtime.Set(name, value)
}

func (g *gojaBackend) getRaw(name string) (any, bool, error) {
	value := g.runtime.Get(name)
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil, false, nil
	}
	export := value.Export()
	switch v := export.(type) {
	case int64, float64, bool, string, nil:
		return v, true, nil
	case int:
		return int64(v), true, nil
	default:
		return fmt.Sprint(v), true, nil
	}
}

func (g *gojaBackend) close() {
	g.runtime = nil
}
