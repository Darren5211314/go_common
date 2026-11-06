package djs

// JSNS JS 引擎工厂命名空间。
type JSNS struct{}

// D新建 创建有状态引擎；引擎名为空或 goja 用 goja；v8 需 CGO。
func (JSNS) D新建(engineName string) (*D引擎, error) {
	return newEngine(engineName)
}
