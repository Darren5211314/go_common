package text

// D文本类 对外命名空间类型（根包通过类型别名导出）。
type D文本类 struct{}

func (D文本类) D取文本左边(内容, 标记 string) string {
	return 取文本左边(内容, 标记)
}

func (D文本类) D取文本右边(内容, 标记 string) string {
	return 取文本右边(内容, 标记)
}

func (D文本类) D取左边(内容 string, 长度 int) string {
	return 取左边(内容, 长度)
}

func (D文本类) D取右边(内容 string, 长度 int) string {
	return 取右边(内容, 长度)
}

func (D文本类) D取文本长度(内容 string) int {
	return 取文本长度(内容)
}

func (D文本类) D寻找文本(内容, 子文本 string) int {
	return 寻找文本(内容, 子文本)
}

func (D文本类) D删首空(内容 string) string {
	return 删首空(内容)
}

func (D文本类) D删尾空(内容 string) string {
	return 删尾空(内容)
}

func (D文本类) D删首尾空(内容 string) string {
	return 删首尾空(内容)
}

func (D文本类) D删全部空(内容 string) string {
	return 删全部空(内容)
}

func (D文本类) D子文本替换(内容, 原文本, 新文本 string) string {
	return 子文本替换(内容, 原文本, 新文本)
}

func (D文本类) D取文本中间(内容 string, 起始位置, 长度 int) string {
	return 取文本中间(内容, 起始位置, 长度)
}

func (D文本类) D取文本中间文本(内容, 左边文本, 右边文本 string) string {
	return 取文本中间文本(内容, 左边文本, 右边文本)
}

func (D文本类) D是否为字母(内容 string) bool {
	return 是否为字母(内容)
}

func (D文本类) D是否为数字(内容 string) bool {
	return 是否为数字(内容)
}

func (D文本类) D是否为汉字(内容 string) bool {
	return 是否为汉字(内容)
}

func (D文本类) D只取字母(内容 string) string {
	return 只取字母(内容)
}

func (D文本类) D只取数字(内容 string) string {
	return 只取数字(内容)
}

func (D文本类) D只取汉字(内容 string) string {
	return 只取汉字(内容)
}

func (D文本类) D取随机字符(位数 int, 小写字母, 大写字母, 数字 bool) string {
	return 取随机字符(位数, 小写字母, 大写字母, 数字)
}

