package djs

import (
	"fmt"

	"github.com/dop251/goja"
)

type goja后端 struct {
	运行时 *goja.Runtime
}

func 新建Goja() (后端, error) {
	return &goja后端{运行时: goja.New()}, nil
}

func (g *goja后端) 执行(代码 string) (string, error) {
	值, 错误 := g.运行时.RunString(代码)
	if 错误 != nil {
		return "", 错误
	}
	if 值 == nil || goja.IsUndefined(值) || goja.IsNull(值) {
		return "", nil
	}
	return 值.String(), nil
}

func (g *goja后端) 设置变量(名 string, 值 any) error {
	return g.运行时.Set(名, 值)
}

func (g *goja后端) 取原始(名 string) (any, bool, error) {
	值 := g.运行时.Get(名)
	if 值 == nil || goja.IsUndefined(值) || goja.IsNull(值) {
		return nil, false, nil
	}
	导出 := 值.Export()
	switch v := 导出.(type) {
	case int64, float64, bool, string, nil:
		return v, true, nil
	case int:
		return int64(v), true, nil
	default:
		return fmt.Sprint(v), true, nil
	}
}

func (g *goja后端) 关闭() {
	g.运行时 = nil
}
