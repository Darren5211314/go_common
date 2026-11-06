// Package text 实现文本处理，按 UTF-8 字符（rune）操作。
// 导出方法均以 D 开头，满足 Go「首字母大写才可跨包访问」的规则。
package text

import (
	"crypto/rand"
	"math/big"
	"strings"
	"unicode"
)

const (
	charsetLower = "abcdefghijklmnopqrstuvwxyz"
	charsetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigit = "0123456789"
)

// LeftOf 返回内容中首次出现标记左侧的文本。
// 标记为空或未找到时返回空字符串。
func leftText(content, flag string) string {
	if flag == "" {
		return ""
	}
	pos := strings.Index(content, flag)
	if pos < 0 {
		return ""
	}
	return content[:pos]
}

// RightOf 返回内容中首次出现标记右侧的文本。
// 标记为空或未找到时返回空字符串。
func rightText(content, flag string) string {
	if flag == "" {
		return ""
	}
	pos := strings.Index(content, flag)
	if pos < 0 {
		return ""
	}
	return content[pos+len(flag):]
}

// Left 按字符数返回内容左侧文本。
// 长度小于等于 0 时返回空字符串；长度超出内容时返回全部内容。
func left(content string, n int) string {
	if n <= 0 {
		return ""
	}
	char := []rune(content)
	if n >= len(char) {
		return content
	}
	return string(char[:n])
}

// Right 按字符数返回内容右侧文本。
// 长度小于等于 0 时返回空字符串；长度超出内容时返回全部内容。
func right(content string, n int) string {
	if n <= 0 {
		return ""
	}
	char := []rune(content)
	total := len(char)
	if n >= total {
		return content
	}
	return string(char[total-n:])
}

// Len 返回内容的字符数（按 UTF-8 rune，不按字节）。
func textLen(content string) int {
	return len([]rune(content))
}

// Find 返回子文本在内容中首次出现的字符位置（从 1 起算）。
// 找不到或子文本为空时返回 0。
func findText(content, sub string) int {
	if sub == "" {
		return 0
	}
	bytePos := strings.Index(content, sub)
	if bytePos < 0 {
		return 0
	}
	return len([]rune(content[:bytePos])) + 1
}

// TrimLeft 去掉内容左侧的空白字符（含空格、制表符、换行、全角空格等）。
func trimLeftSpace(content string) string {
	return strings.TrimLeftFunc(content, unicode.IsSpace)
}

// TrimRight 去掉内容右侧的空白字符（含空格、制表符、换行、全角空格等）。
func trimRightSpace(content string) string {
	return strings.TrimRightFunc(content, unicode.IsSpace)
}

// Trim 去掉内容两侧的空白字符（含空格、制表符、换行、全角空格等）；中间空白保留。
func trimSpace(content string) string {
	return strings.TrimSpace(content)
}

// RemoveSpace 去掉内容中全部空白字符（含中间）。
func trimAllSpace(content string) string {
	var result strings.Builder
	result.Grow(len(content))
	for _, ch := range content {
		if !unicode.IsSpace(ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// Replace 将内容中全部原文本替换为新文本。
// 原文本为空时原样返回内容。
func replaceSub(content, rawText, newText string) string {
	if rawText == "" {
		return content
	}
	return strings.ReplaceAll(content, rawText, newText)
}

// Mid 从起始位置起按字符截取指定长度（起始位置从 1 起算）。
// 起始位置小于 1 或长度小于等于 0 时返回空字符串；越界则截到末尾。
func mid(content string, startPos, n int) string {
	if startPos < 1 || n <= 0 {
		return ""
	}
	char := []rune(content)
	total := len(char)
	startIndex := startPos - 1
	if startIndex >= total {
		return ""
	}
	endIndex := startIndex + n
	if endIndex > total {
		endIndex = total
	}
	return string(char[startIndex:endIndex])
}

// Between 返回左边文本与右边文本之间的内容（取各自首次出现）。
// 左边或右边为空、任一侧未找到时返回空字符串。
func midText(content, leftStr, rightStr string) string {
	if leftStr == "" || rightStr == "" {
		return ""
	}
	leftPos := strings.Index(content, leftStr)
	if leftPos < 0 {
		return ""
	}
	start := leftPos + len(leftStr)
	remain := content[start:]
	rightPos := strings.Index(remain, rightStr)
	if rightPos < 0 {
		return ""
	}
	return remain[:rightPos]
}

func letterOK(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func digitOK(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func hanOK(ch rune) bool {
	return unicode.Is(unicode.Han, ch)
}

func allMatch(content string, check func(rune) bool) bool {
	if content == "" {
		return false
	}
	for _, ch := range content {
		if !check(ch) {
			return false
		}
	}
	return true
}

func keepMatching(content string, check func(rune) bool) string {
	var result strings.Builder
	result.Grow(len(content))
	for _, ch := range content {
		if check(ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// IsLetter 判断内容是否非空且每个字符都是英文字母（a-z / A-Z）。
func isLetter(content string) bool {
	return allMatch(content, letterOK)
}

// IsDigit 判断内容是否非空且每个字符都是数字（0-9）。
func isDigit(content string) bool {
	return allMatch(content, digitOK)
}

// IsHan 判断内容是否非空且每个字符都是汉字（Unicode Han）。
func isHan(content string) bool {
	return allMatch(content, hanOK)
}

// OnlyLetters 保留内容中的英文字母，其余丢弃。
func lettersOnly(content string) string {
	return keepMatching(content, letterOK)
}

// OnlyDigits 保留内容中的数字字符，其余丢弃。
func digitsOnly(content string) string {
	return keepMatching(content, digitOK)
}

// OnlyHan 保留内容中的汉字，其余丢弃。
func hanOnly(content string) string {
	return keepMatching(content, hanOK)
}

// Random 生成指定位数的随机串。
// 小写字母、大写字母、数字 为字符集开关；三者都为 false 时默认「小写字母 + 数字」。
// 位数小于等于 0 时返回空字符串。使用 crypto/rand。
func randomChars(bits int, lowerLetters, upperLetters, digit bool) string {
	if bits <= 0 {
		return ""
	}
	if !lowerLetters && !upperLetters && !digit {
		lowerLetters = true
		digit = true
	}
	charset := ""
	if lowerLetters {
		charset += charsetLower
	}
	if upperLetters {
		charset += charsetUpper
	}
	if digit {
		charset += charsetDigit
	}
	limit := big.NewInt(int64(len(charset)))
	result := make([]byte, bits)
	for i := 0; i < bits; i++ {
		ordinal, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return ""
		}
		result[i] = charset[ordinal.Int64()]
	}
	return string(result)
}
