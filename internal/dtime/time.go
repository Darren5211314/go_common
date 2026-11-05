// Package dtime 实现时间戳与格式化，使用本地时区，不保存全局可变状态。
// 导出方法均以 D 开头，满足 Go「首字母大写才可跨包访问」的规则。
package dtime

import (
	"strings"
	"time"
)

const 毫秒阈值 = int64(1_000_000_000_000)

var 解析候选布局 = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05.000000000",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05.000",
	"2006/01/02 15:04:05",
	"2006年01月02日 15时04分05秒",
	"2006-01-02 15:04",
	"2006/01/02 15:04",
	"2006年01月02日",
	"2006-01-02",
	"2006/01/02",
	"20060102150405",
	"20060102",
	"15:04:05",
	"15时04分05秒",
	"15:04",
}

func 绝对值(值 int64) int64 {
	if 值 < 0 {
		return -值
	}
	return 值
}

func 时间戳到时间(戳 int64) time.Time {
	if 绝对值(戳) >= 毫秒阈值 {
		return time.UnixMilli(戳).Local()
	}
	return time.Unix(戳, 0).Local()
}

// 通配符转布局 将 yyyy-mm-dd hh:mm:ss 等通配符转为内部时间布局。
func 通配符转布局(格式 string) string {
	if 格式 == "" {
		return 格式
	}

	结果 := 格式
	短语 := []struct {
		通配符 string
		布局  string
	}{
		{"yyyy-mm-dd hh:mm:ss", "2006-01-02 15:04:05"},
		{"yyyy/mm/dd hh:mm:ss", "2006/01/02 15:04:05"},
		{"yyyy-MM-dd HH:mm:ss", "2006-01-02 15:04:05"},
		{"yyyy/MM/dd HH:mm:ss", "2006/01/02 15:04:05"},
		{"yyyy年mm月dd日 hh时mm分ss秒", "2006年01月02日 15时04分05秒"},
		{"yyyy-mm-dd hh:mm", "2006-01-02 15:04"},
		{"yyyy/mm/dd hh:mm", "2006/01/02 15:04"},
		{"yyyy-MM-dd HH:mm", "2006-01-02 15:04"},
		{"yyyy年mm月dd日", "2006年01月02日"},
		{"yyyy-mm-dd", "2006-01-02"},
		{"yyyy/mm/dd", "2006/01/02"},
		{"yyyy-MM-dd", "2006-01-02"},
		{"yyyy/MM/dd", "2006/01/02"},
		{"hh:mm:ss", "15:04:05"},
		{"HH:mm:ss", "15:04:05"},
		{"hh时mm分ss秒", "15时04分05秒"},
		{"hh:mm", "15:04"},
		{"HH:mm", "15:04"},
	}
	for _, 项 := range 短语 {
		结果 = strings.ReplaceAll(结果, 项.通配符, 项.布局)
	}

	// 标记替换：月优先于分（mm月 / MM）；短语未覆盖时剩余 mm 视为分钟。
	标记 := []struct {
		通配符 string
		布局  string
	}{
		{"yyyy", "2006"},
		{"YYYY", "2006"},
		{"SSS", "000"},
		{"sss", "000"},
		{"mm月", "01月"},
		{"MM", "01"},
		{"dd", "02"},
		{"DD", "02"},
		{"HH", "15"},
		{"hh", "15"},
		{"mm分", "04分"},
		{"mm", "04"},
		{"ss", "05"},
		{"SS", "05"},
	}
	for _, 项 := range 标记 {
		结果 = strings.ReplaceAll(结果, 项.通配符, 项.布局)
	}
	return 结果
}

// 解析时间文本 自动尝试常见时间格式，成功返回本地时区时间。
func 解析时间文本(文本 string) (time.Time, bool) {
	文本 = strings.TrimSpace(文本)
	if 文本 == "" {
		return time.Time{}, false
	}
	for _, 布局 := range 解析候选布局 {
		时刻, 错误 := time.ParseInLocation(布局, 文本, time.Local)
		if 错误 == nil {
			return 时刻, true
		}
	}
	return time.Time{}, false
}

// D取现行时间戳 取十位为 true 返回秒级（约 10 位），为 false 返回毫秒级（约 13 位）。
// 调用时传 false 即「默认」13 位。
func 取现行时间戳(取十位 bool) int64 {
	if 取十位 {
		return time.Now().Unix()
	}
	return time.Now().UnixMilli()
}

// D时间戳转时间 将时间戳按本地时区格式化为文本。
// 自动识别：绝对值大于等于 1e12 按毫秒，否则按秒。
// 格式使用通配符，例如 yyyy-mm-dd hh:mm:ss。
func 时间戳转时间(戳 int64, 格式 string) string {
	return 时间戳到时间(戳).Format(通配符转布局(格式))
}

// D时间转时间戳 自动识别时间文本格式并转为时间戳。
// 取十位为 true 返回秒，为 false 返回毫秒；解析失败返回 0。
func 时间转时间戳(时间文本 string, 取十位 bool) int64 {
	时刻, 成功 := 解析时间文本(时间文本)
	if !成功 {
		return 0
	}
	if 取十位 {
		return 时刻.Unix()
	}
	return 时刻.UnixMilli()
}

// D时间格式化 将时间文本格式化为目标通配符格式。
// 输入格式自动识别；时间文本为空时使用现行时间；解析失败返回空字符串。
func 时间格式化(时间文本, 目标格式 string) string {
	var 时刻 time.Time
	if strings.TrimSpace(时间文本) == "" {
		时刻 = time.Now()
	} else {
		解析结果, 成功 := 解析时间文本(时间文本)
		if !成功 {
			return ""
		}
		时刻 = 解析结果
	}
	return 时刻.Format(通配符转布局(目标格式))
}

// D延时 阻塞当前协程指定毫秒；毫秒 <= 0 时立即返回。
func 延时(毫秒 int) {
	if 毫秒 <= 0 {
		return
	}
	time.Sleep(time.Duration(毫秒) * time.Millisecond)
}
