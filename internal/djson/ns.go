package djson

// JSON类 对外命名空间类型（根包通过类型别名导出）。
type JSON类 struct{}

func (JSON类) D是否有效(文本 string) bool {
	return 是否有效(文本)
}

func (JSON类) D是否存在(文本, 路径 string) bool {
	return 是否存在(文本, 路径)
}

func (JSON类) D取成员数(文本, 路径 string) int {
	return 取成员数(文本, 路径)
}

func (JSON类) D取文本(文本, 路径 string) string {
	return 取文本(文本, 路径)
}

func (JSON类) D取整数(文本, 路径 string) int64 {
	return 取整数(文本, 路径)
}

func (JSON类) D取小数(文本, 路径 string) float64 {
	return 取小数(文本, 路径)
}

func (JSON类) D取逻辑(文本, 路径 string) bool {
	return 取逻辑(文本, 路径)
}

func (JSON类) D取JSON(文本, 路径 string) string {
	return 取JSON(文本, 路径)
}

func (JSON类) D设置文本(文本, 路径, 值 string) (string, error) {
	return 设置文本(文本, 路径, 值)
}

func (JSON类) D设置整数(文本, 路径 string, 值 int64) (string, error) {
	return 设置整数(文本, 路径, 值)
}

func (JSON类) D设置小数(文本, 路径 string, 值 float64) (string, error) {
	return 设置小数(文本, 路径, 值)
}

func (JSON类) D设置逻辑(文本, 路径 string, 值 bool) (string, error) {
	return 设置逻辑(文本, 路径, 值)
}

func (JSON类) D设置JSON(文本, 路径, JSON片段 string) (string, error) {
	return 设置JSON(文本, 路径, JSON片段)
}

func (JSON类) D压缩(文本 string) string {
	return 压缩(文本)
}

func (JSON类) D格式化(文本 string) string {
	return 格式化(文本)
}

