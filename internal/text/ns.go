package text

// D文本类 对外命名空间类型（根包通过类型别名导出）。
type D文本类 struct{}

func (D文本类) D取文本左边(content, flag string) string {
	return leftText(content, flag)
}

func (D文本类) D取文本右边(content, flag string) string {
	return rightText(content, flag)
}

func (D文本类) D取左边(content string, n int) string {
	return left(content, n)
}

func (D文本类) D取右边(content string, n int) string {
	return right(content, n)
}

func (D文本类) D取文本长度(content string) int {
	return textLen(content)
}

func (D文本类) D寻找文本(content, sub string) int {
	return findText(content, sub)
}

func (D文本类) D删首空(content string) string {
	return trimLeftSpace(content)
}

func (D文本类) D删尾空(content string) string {
	return trimRightSpace(content)
}

func (D文本类) D删首尾空(content string) string {
	return trimSpace(content)
}

func (D文本类) D删全部空(content string) string {
	return trimAllSpace(content)
}

func (D文本类) D子文本替换(content, rawText, newText string) string {
	return replaceSub(content, rawText, newText)
}

func (D文本类) D取文本中间(content string, startPos, n int) string {
	return mid(content, startPos, n)
}

func (D文本类) D取文本中间文本(content, leftStr, rightStr string) string {
	return midText(content, leftStr, rightStr)
}

func (D文本类) D是否为字母(content string) bool {
	return isLetter(content)
}

func (D文本类) D是否为数字(content string) bool {
	return isDigit(content)
}

func (D文本类) D是否为汉字(content string) bool {
	return isHan(content)
}

func (D文本类) D只取字母(content string) string {
	return lettersOnly(content)
}

func (D文本类) D只取数字(content string) string {
	return digitsOnly(content)
}

func (D文本类) D只取汉字(content string) string {
	return hanOnly(content)
}

func (D文本类) D取随机字符(bits int, lowerLetters, upperLetters, digit bool) string {
	return randomChars(bits, lowerLetters, upperLetters, digit)
}

