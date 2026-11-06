package dcodec

import "unicode/utf8"

// 是否UTF8 判断文本是否为合法 UTF-8。
func isUTF8(text string) bool {
	return utf8.ValidString(text)
}
