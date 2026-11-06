package text

// TextNS 命名空间类型（根包以未导出别名挂到 var Text）。
type TextNS struct{}

func (TextNS) D取文本左边(content, flag string) string {
	return leftText(content, flag)
}

func (TextNS) D取文本右边(content, flag string) string {
	return rightText(content, flag)
}

func (TextNS) D取左边(content string, n int) string {
	return left(content, n)
}

func (TextNS) D取右边(content string, n int) string {
	return right(content, n)
}

func (TextNS) D取文本长度(content string) int {
	return textLen(content)
}

func (TextNS) D寻找文本(content, sub string) int {
	return findText(content, sub)
}

func (TextNS) D删首空(content string) string {
	return trimLeftSpace(content)
}

func (TextNS) D删尾空(content string) string {
	return trimRightSpace(content)
}

func (TextNS) D删首尾空(content string) string {
	return trimSpace(content)
}

func (TextNS) D删全部空(content string) string {
	return trimAllSpace(content)
}

func (TextNS) D子文本替换(content, rawText, newText string) string {
	return replaceSub(content, rawText, newText)
}

func (TextNS) D取文本中间(content string, startPos, n int) string {
	return mid(content, startPos, n)
}

func (TextNS) D取文本中间文本(content, leftStr, rightStr string) string {
	return midText(content, leftStr, rightStr)
}

func (TextNS) D是否为字母(content string) bool {
	return isLetter(content)
}

func (TextNS) D是否为数字(content string) bool {
	return isDigit(content)
}

func (TextNS) D是否为汉字(content string) bool {
	return isHan(content)
}

func (TextNS) D只取字母(content string) string {
	return lettersOnly(content)
}

func (TextNS) D只取数字(content string) string {
	return digitsOnly(content)
}

func (TextNS) D只取汉字(content string) string {
	return hanOnly(content)
}

func (TextNS) D取随机字符(bits int, lowerLetters, upperLetters, digit bool) string {
	return randomChars(bits, lowerLetters, upperLetters, digit)
}

