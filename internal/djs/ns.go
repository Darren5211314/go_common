package djs

// JS类 JS 引擎工厂命名空间。
type JS类 struct{}

// D新建 创建有状态引擎；引擎名为空或 goja 用 goja；v8 需 CGO。
func (JS类) D新建(引擎名 string) (*D引擎, error) {
	return 新建引擎(引擎名)
}
