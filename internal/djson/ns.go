package djson

// JSONNS 对外命名空间类型（根包通过类型别名导出）。
type JSONNS struct{}

func (JSONNS) Valid(text string) bool {
	return isValidIP(text)
}

func (JSONNS) Has(text, path string) bool {
	return existsAt(text, path)
}

func (JSONNS) Count(text, path string) int {
	return memberCount(text, path)
}

func (JSONNS) GetString(text, path string) string {
	return getText(text, path)
}

func (JSONNS) GetInt(text, path string) int64 {
	return getInt(text, path)
}

func (JSONNS) GetFloat(text, path string) float64 {
	return getFloat(text, path)
}

func (JSONNS) GetBool(text, path string) bool {
	return getBool(text, path)
}

func (JSONNS) GetJSON(text, path string) string {
	return getJSON(text, path)
}

func (JSONNS) SetString(text, path, value string) (string, error) {
	return setText(text, path, value)
}

func (JSONNS) SetInt(text, path string, value int64) (string, error) {
	return setInt(text, path, value)
}

func (JSONNS) SetFloat(text, path string, value float64) (string, error) {
	return setFloat(text, path, value)
}

func (JSONNS) SetBool(text, path string, value bool) (string, error) {
	return setBool(text, path, value)
}

func (JSONNS) SetJSON(text, path, jsonFrag string) (string, error) {
	return setJSON(text, path, jsonFrag)
}

func (JSONNS) Compact(text string) string {
	return compact(text)
}

func (JSONNS) Pretty(text string) string {
	return pretty(text)
}

