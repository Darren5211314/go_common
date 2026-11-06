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
	errClosed     = errors.New("引擎已关闭")
	errEmptyCode     = errors.New("代码为空")
	errEmptyVarName = errors.New("变量名为空")
	errVarNotFound = errors.New("变量不存在")
	errTypeMismatch   = errors.New("变量类型不符")
	errUnsupportedEngine = errors.New("不支持的引擎")
)

// 后端 具体 JS 实现。
type backend interface {
	eval(code string) (string, error)
	setVar(name string, value any) error
	getRaw(name string) (any, bool, error)
	close()
}

// Engine 有状态 JS 引擎（内部加锁，多 goroutine 共用会串行，不崩溃）。
type Engine struct {
	mu   sync.Mutex
	backend  backend
	closed bool
	name  string
}

// 新建引擎 按名称创建后端；空或 goja 用 goja；v8 需 CGO。
func newEngine(engineName string) (*Engine, error) {
	name := strings.ToLower(strings.TrimSpace(engineName))
	if name == "" {
		name = "goja"
	}
	var (
		backendInst backend
		err   error
	)
	switch name {
	case "goja":
		backendInst, err = newGoja()
	case "v8":
		backendInst, err = newV8()
	default:
		return nil, fmt.Errorf("%w: %s（支持 goja / v8）", errUnsupportedEngine, engineName)
	}
	if err != nil {
		return nil, err
	}
	return &Engine{backend: backendInst, name: name}, nil
}

// Close 释放引擎资源。
func (e *Engine) Close() {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return
	}
	e.closed = true
	if e.backend != nil {
		e.backend.close()
		e.backend = nil
	}
}

// Name 返回 goja 或 v8。
func (e *Engine) Name() string {
	if e == nil {
		return ""
	}
	return e.name
}

// Eval 执行代码，返回最后表达式的字符串形式。
func (e *Engine) Eval(code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", errEmptyCode
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed || e.backend == nil {
		return "", errClosed
	}
	return e.backend.eval(code)
}

// EvalFile 读取文件并执行。
func (e *Engine) EvalFile(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("路径为空")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return e.Eval(string(data))
}

// SetVar 设置全局变量。
func (e *Engine) SetVar(name string, value any) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errEmptyVarName
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed || e.backend == nil {
		return errClosed
	}
	return e.backend.setVar(name, normalizeValue(value))
}

func (e *Engine) getRawLocked(name string) (any, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errEmptyVarName
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed || e.backend == nil {
		return nil, errClosed
	}
	value, exists, err := e.backend.getRaw(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errVarNotFound
	}
	return value, nil
}

// GetString 读取全局变量为字符串。
func (e *Engine) GetString(name string) (string, error) {
	value, err := e.getRawLocked(name)
	if err != nil {
		return "", err
	}
	switch v := value.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return fmt.Sprint(v), nil
	}
}

// GetInt 读取全局变量为 int64。
func (e *Engine) GetInt(name string) (int64, error) {
	value, err := e.getRawLocked(name)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
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
		return 0, errTypeMismatch
	}
}

// GetFloat 读取全局变量为 float64。
func (e *Engine) GetFloat(name string) (float64, error) {
	value, err := e.getRawLocked(name)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	default:
		return 0, errTypeMismatch
	}
}

// GetBool 读取全局变量为 bool。
func (e *Engine) GetBool(name string) (bool, error) {
	value, err := e.getRawLocked(name)
	if err != nil {
		return false, err
	}
	b, ok := value.(bool)
	if !ok {
		return false, errTypeMismatch
	}
	return b, nil
}

func normalizeValue(value any) any {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
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
