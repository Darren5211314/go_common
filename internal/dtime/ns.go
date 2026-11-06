package dtime

// D时间类 对外命名空间类型（根包通过类型别名导出）。
type D时间类 struct{}

func (D时间类) D取现行时间戳(tenDigits bool) int64 {
	return nowTimestamp(tenDigits)
}

func (D时间类) D时间戳转时间(ts int64, format string) string {
	return timestampFormat(ts, format)
}

func (D时间类) D时间转时间戳(timeText string, tenDigits bool) int64 {
	return timeToTimestamp(timeText, tenDigits)
}

func (D时间类) D时间格式化(timeText, targetLayout string) string {
	return formatTime(timeText, targetLayout)
}

func (D时间类) D延时(ms int) {
	sleepMs(ms)
}

