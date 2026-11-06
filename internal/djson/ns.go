package djson

// JSONNS 对外命名空间类型（根包通过类型别名导出）。
type JSONNS struct{}

func (JSONNS) D是否有效(text string) bool {
	return isValidIP(text)
}

func (JSONNS) D是否存在(text, path string) bool {
	return existsAt(text, path)
}

func (JSONNS) D取成员数(text, path string) int {
	return memberCount(text, path)
}

func (JSONNS) D取文本(text, path string) string {
	return getText(text, path)
}

func (JSONNS) D取整数(text, path string) int64 {
	return getInt(text, path)
}

func (JSONNS) D取小数(text, path string) float64 {
	return getFloat(text, path)
}

func (JSONNS) D取逻辑(text, path string) bool {
	return getBool(text, path)
}

func (JSONNS) D取JSON(text, path string) string {
	return getJSON(text, path)
}

func (JSONNS) D设置文本(text, path, value string) (string, error) {
	return setText(text, path, value)
}

func (JSONNS) D设置整数(text, path string, value int64) (string, error) {
	return setInt(text, path, value)
}

func (JSONNS) D设置小数(text, path string, value float64) (string, error) {
	return setFloat(text, path, value)
}

func (JSONNS) D设置逻辑(text, path string, value bool) (string, error) {
	return setBool(text, path, value)
}

func (JSONNS) D设置JSON(text, path, jsonFrag string) (string, error) {
	return setJSON(text, path, jsonFrag)
}

func (JSONNS) D压缩(text string) string {
	return compact(text)
}

func (JSONNS) D格式化(text string) string {
	return pretty(text)
}

