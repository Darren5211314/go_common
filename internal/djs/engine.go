// Package djs 嵌入式 JS 引擎：默认 goja（免 CGO），可选 V8（需 CGO）。
package djs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	错误已关闭     = errors.New("引擎已关闭")
	错误空代码     = errors.New("代码为空")
	错误变量名为空 = errors.New("变量名为空")
	错误变量不存在 = errors.New("变量不存在")
	错误类型不符   = errors.New("变量类型不符")
	错误不支持引擎 = errors.New("不支持的引擎")
)

// 后端 具体 JS 实现。
type 后端 interface {
	执行(代码 string) (string, error)
	设置变量(名 string, 值 any) error
	取原始(名 string) (any, bool, error)
	关闭()
}

// D引擎 有状态 JS 引擎（内部加锁，多 goroutine 共用会串行，不崩溃）。
type D引擎 struct {
	锁   sync.Mutex
	后端  后端
	已关闭 bool
	名称  string
}

// 新建引擎 按名称创建后端；空或 goja 用 goja；v8 需 CGO。
func 新建引擎(引擎名 string) (*D引擎, error) {
	名 := strings.ToLower(strings.TrimSpace(引擎名))
	if 名 == "" {
		名 = "goja"
	}
	var (
		后端实例 后端
		错误   error
	)
	switch 名 {
	case "goja":
		后端实例, 错误 = 新建Goja()
	case "v8":
		后端实例, 错误 = 新建V8()
	default:
		return nil, fmt.Errorf("%w: %s（支持 goja / v8）", 错误不支持引擎, 引擎名)
	}
	if 错误 != nil {
		return nil, 错误
	}
	return &D引擎{后端: 后端实例, 名称: 名}, nil
}

// D关闭 释放引擎资源。
func (e *D引擎) D关闭() {
	if e == nil {
		return
	}
	e.锁.Lock()
	defer e.锁.Unlock()
	if e.已关闭 {
		return
	}
	e.已关闭 = true
	if e.后端 != nil {
		e.后端.关闭()
		e.后端 = nil
	}
}

// D引擎名 返回 goja 或 v8。
func (e *D引擎) D引擎名() string {
	if e == nil {
		return ""
	}
	return e.名称
}

// D执行 执行代码，返回最后表达式的字符串形式。
func (e *D引擎) D执行(代码 string) (string, error) {
	代码 = strings.TrimSpace(代码)
	if 代码 == "" {
		return "", 错误空代码
	}
	e.锁.Lock()
	defer e.锁.Unlock()
	if e.已关闭 || e.后端 == nil {
		return "", 错误已关闭
	}
	return e.后端.执行(代码)
}

// D执行文件 读取文件并执行。
func (e *D引擎) D执行文件(路径 string) (string, error) {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return "", errors.New("路径为空")
	}
	数据, 错误 := os.ReadFile(路径)
	if 错误 != nil {
		return "", 错误
	}
	return e.D执行(string(数据))
}

// D设置变量 设置全局变量。
func (e *D引擎) D设置变量(名 string, 值 any) error {
	名 = strings.TrimSpace(名)
	if 名 == "" {
		return 错误变量名为空
	}
	e.锁.Lock()
	defer e.锁.Unlock()
	if e.已关闭 || e.后端 == nil {
		return 错误已关闭
	}
	return e.后端.设置变量(名, 规范化值(值))
}

func (e *D引擎) 取原始锁定(名 string) (any, error) {
	名 = strings.TrimSpace(名)
	if 名 == "" {
		return nil, 错误变量名为空
	}
	e.锁.Lock()
	defer e.锁.Unlock()
	if e.已关闭 || e.后端 == nil {
		return nil, 错误已关闭
	}
	值, 存在, 错误 := e.后端.取原始(名)
	if 错误 != nil {
		return nil, 错误
	}
	if !存在 {
		return nil, 错误变量不存在
	}
	return 值, nil
}

// D取文本 读取全局变量为字符串。
func (e *D引擎) D取文本(名 string) (string, error) {
	值, 错误 := e.取原始锁定(名)
	if 错误 != nil {
		return "", 错误
	}
	switch v := 值.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return fmt.Sprint(v), nil
	}
}

// D取整数 读取全局变量为 int64。
func (e *D引擎) D取整数(名 string) (int64, error) {
	值, 错误 := e.取原始锁定(名)
	if 错误 != nil {
		return 0, 错误
	}
	switch v := 值.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, 错误类型不符
	}
}

// D取小数 读取全局变量为 float64。
func (e *D引擎) D取小数(名 string) (float64, error) {
	值, 错误 := e.取原始锁定(名)
	if 错误 != nil {
		return 0, 错误
	}
	switch v := 值.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	default:
		return 0, 错误类型不符
	}
}

// D取逻辑 读取全局变量为 bool。
func (e *D引擎) D取逻辑(名 string) (bool, error) {
	值, 错误 := e.取原始锁定(名)
	if 错误 != nil {
		return false, 错误
	}
	b, ok := 值.(bool)
	if !ok {
		return false, 错误类型不符
	}
	return b, nil
}

func 规范化值(值 any) any {
	if 值 == nil {
		return nil
	}
	switch v := 值.(type) {
	case string, bool, int64, float64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case float32:
		return float64(v)
	default:
		return fmt.Sprint(v)
	}
}
