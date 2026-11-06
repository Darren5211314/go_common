package text

// TextNS 命名空间类型（根包以未导出别名挂到 var Text）。
type TextNS struct{}

func (TextNS) LeftOf(content, flag string) string {
	return leftText(content, flag)
}

func (TextNS) RightOf(content, flag string) string {
	return rightText(content, flag)
}

func (TextNS) Left(content string, n int) string {
	return left(content, n)
}

func (TextNS) Right(content string, n int) string {
	return right(content, n)
}

func (TextNS) Len(content string) int {
	return textLen(content)
}

func (TextNS) Find(content, sub string) int {
	return findText(content, sub)
}

func (TextNS) TrimLeft(content string) string {
	return trimLeftSpace(content)
}

func (TextNS) TrimRight(content string) string {
	return trimRightSpace(content)
}

func (TextNS) Trim(content string) string {
	return trimSpace(content)
}

func (TextNS) RemoveSpace(content string) string {
	return trimAllSpace(content)
}

func (TextNS) Replace(content, rawText, newText string) string {
	return replaceSub(content, rawText, newText)
}

func (TextNS) Mid(content string, startPos, n int) string {
	return mid(content, startPos, n)
}

func (TextNS) Between(content, leftStr, rightStr string) string {
	return midText(content, leftStr, rightStr)
}

func (TextNS) IsLetter(content string) bool {
	return isLetter(content)
}

func (TextNS) IsDigit(content string) bool {
	return isDigit(content)
}

func (TextNS) IsHan(content string) bool {
	return isHan(content)
}

func (TextNS) OnlyLetters(content string) string {
	return lettersOnly(content)
}

func (TextNS) OnlyDigits(content string) string {
	return digitsOnly(content)
}

func (TextNS) OnlyHan(content string) string {
	return hanOnly(content)
}

func (TextNS) Random(bits int, lowerLetters, upperLetters, digit bool) string {
	return randomChars(bits, lowerLetters, upperLetters, digit)
}

