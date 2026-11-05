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
	字符集小写 = "abcdefghijklmnopqrstuvwxyz"
	字符集大写 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	字符集数字 = "0123456789"
)

// D取文本左边 返回内容中首次出现标记左侧的文本。
// 标记为空或未找到时返回空字符串。
func 取文本左边(内容, 标记 string) string {
	if 标记 == "" {
		return ""
	}
	位置 := strings.Index(内容, 标记)
	if 位置 < 0 {
		return ""
	}
	return 内容[:位置]
}

// D取文本右边 返回内容中首次出现标记右侧的文本。
// 标记为空或未找到时返回空字符串。
func 取文本右边(内容, 标记 string) string {
	if 标记 == "" {
		return ""
	}
	位置 := strings.Index(内容, 标记)
	if 位置 < 0 {
		return ""
	}
	return 内容[位置+len(标记):]
}

// D取左边 按字符数返回内容左侧文本。
// 长度小于等于 0 时返回空字符串；长度超出内容时返回全部内容。
func 取左边(内容 string, 长度 int) string {
	if 长度 <= 0 {
		return ""
	}
	字符 := []rune(内容)
	if 长度 >= len(字符) {
		return 内容
	}
	return string(字符[:长度])
}

// D取右边 按字符数返回内容右侧文本。
// 长度小于等于 0 时返回空字符串；长度超出内容时返回全部内容。
func 取右边(内容 string, 长度 int) string {
	if 长度 <= 0 {
		return ""
	}
	字符 := []rune(内容)
	总数 := len(字符)
	if 长度 >= 总数 {
		return 内容
	}
	return string(字符[总数-长度:])
}

// D取文本长度 返回内容的字符数（按 UTF-8 rune，不按字节）。
func 取文本长度(内容 string) int {
	return len([]rune(内容))
}

// D寻找文本 返回子文本在内容中首次出现的字符位置（从 1 起算）。
// 找不到或子文本为空时返回 0。
func 寻找文本(内容, 子文本 string) int {
	if 子文本 == "" {
		return 0
	}
	字节位置 := strings.Index(内容, 子文本)
	if 字节位置 < 0 {
		return 0
	}
	return len([]rune(内容[:字节位置])) + 1
}

// D删首空 去掉内容左侧的空白字符（含空格、制表符、换行、全角空格等）。
func 删首空(内容 string) string {
	return strings.TrimLeftFunc(内容, unicode.IsSpace)
}

// D删尾空 去掉内容右侧的空白字符（含空格、制表符、换行、全角空格等）。
func 删尾空(内容 string) string {
	return strings.TrimRightFunc(内容, unicode.IsSpace)
}

// D删首尾空 去掉内容两侧的空白字符（含空格、制表符、换行、全角空格等）；中间空白保留。
func 删首尾空(内容 string) string {
	return strings.TrimSpace(内容)
}

// D删全部空 去掉内容中全部空白字符（含中间）。
func 删全部空(内容 string) string {
	var 结果 strings.Builder
	结果.Grow(len(内容))
	for _, 字 := range 内容 {
		if !unicode.IsSpace(字) {
			结果.WriteRune(字)
		}
	}
	return 结果.String()
}

// D子文本替换 将内容中全部原文本替换为新文本。
// 原文本为空时原样返回内容。
func 子文本替换(内容, 原文本, 新文本 string) string {
	if 原文本 == "" {
		return 内容
	}
	return strings.ReplaceAll(内容, 原文本, 新文本)
}

// D取文本中间 从起始位置起按字符截取指定长度（起始位置从 1 起算）。
// 起始位置小于 1 或长度小于等于 0 时返回空字符串；越界则截到末尾。
func 取文本中间(内容 string, 起始位置, 长度 int) string {
	if 起始位置 < 1 || 长度 <= 0 {
		return ""
	}
	字符 := []rune(内容)
	总数 := len(字符)
	起始下标 := 起始位置 - 1
	if 起始下标 >= 总数 {
		return ""
	}
	结束下标 := 起始下标 + 长度
	if 结束下标 > 总数 {
		结束下标 = 总数
	}
	return string(字符[起始下标:结束下标])
}

// D取文本中间文本 返回左边文本与右边文本之间的内容（取各自首次出现）。
// 左边或右边为空、任一侧未找到时返回空字符串。
func 取文本中间文本(内容, 左边文本, 右边文本 string) string {
	if 左边文本 == "" || 右边文本 == "" {
		return ""
	}
	左位置 := strings.Index(内容, 左边文本)
	if 左位置 < 0 {
		return ""
	}
	起始 := 左位置 + len(左边文本)
	剩余 := 内容[起始:]
	右位置 := strings.Index(剩余, 右边文本)
	if 右位置 < 0 {
		return ""
	}
	return 剩余[:右位置]
}

func 是字母(字 rune) bool {
	return (字 >= 'a' && 字 <= 'z') || (字 >= 'A' && 字 <= 'Z')
}

func 是数字(字 rune) bool {
	return 字 >= '0' && 字 <= '9'
}

func 是汉字(字 rune) bool {
	return unicode.Is(unicode.Han, 字)
}

func 整段符合(内容 string, 判定 func(rune) bool) bool {
	if 内容 == "" {
		return false
	}
	for _, 字 := range 内容 {
		if !判定(字) {
			return false
		}
	}
	return true
}

func 只取符合(内容 string, 判定 func(rune) bool) string {
	var 结果 strings.Builder
	结果.Grow(len(内容))
	for _, 字 := range 内容 {
		if 判定(字) {
			结果.WriteRune(字)
		}
	}
	return 结果.String()
}

// D是否为字母 判断内容是否非空且每个字符都是英文字母（a-z / A-Z）。
func 是否为字母(内容 string) bool {
	return 整段符合(内容, 是字母)
}

// D是否为数字 判断内容是否非空且每个字符都是数字（0-9）。
func 是否为数字(内容 string) bool {
	return 整段符合(内容, 是数字)
}

// D是否为汉字 判断内容是否非空且每个字符都是汉字（Unicode Han）。
func 是否为汉字(内容 string) bool {
	return 整段符合(内容, 是汉字)
}

// D只取字母 保留内容中的英文字母，其余丢弃。
func 只取字母(内容 string) string {
	return 只取符合(内容, 是字母)
}

// D只取数字 保留内容中的数字字符，其余丢弃。
func 只取数字(内容 string) string {
	return 只取符合(内容, 是数字)
}

// D只取汉字 保留内容中的汉字，其余丢弃。
func 只取汉字(内容 string) string {
	return 只取符合(内容, 是汉字)
}

// D取随机字符 生成指定位数的随机串。
// 小写字母、大写字母、数字 为字符集开关；三者都为 false 时默认「小写字母 + 数字」。
// 位数小于等于 0 时返回空字符串。使用 crypto/rand。
func 取随机字符(位数 int, 小写字母, 大写字母, 数字 bool) string {
	if 位数 <= 0 {
		return ""
	}
	if !小写字母 && !大写字母 && !数字 {
		小写字母 = true
		数字 = true
	}
	字符集 := ""
	if 小写字母 {
		字符集 += 字符集小写
	}
	if 大写字母 {
		字符集 += 字符集大写
	}
	if 数字 {
		字符集 += 字符集数字
	}
	上限 := big.NewInt(int64(len(字符集)))
	结果 := make([]byte, 位数)
	for i := 0; i < 位数; i++ {
		序号, 错误 := rand.Int(rand.Reader, 上限)
		if 错误 != nil {
			return ""
		}
		结果[i] = 字符集[序号.Int64()]
	}
	return string(结果)
}
