package dtime

// TimeNS 命名空间类型（根包以未导出别名挂到 var Time）。
type TimeNS struct{}

func (TimeNS) D取现行时间戳(tenDigits bool) int64 {
	return nowTimestamp(tenDigits)
}

func (TimeNS) D时间戳转时间(ts int64, format string) string {
	return timestampFormat(ts, format)
}

func (TimeNS) D时间转时间戳(timeText string, tenDigits bool) int64 {
	return timeToTimestamp(timeText, tenDigits)
}

func (TimeNS) D时间格式化(timeText, targetLayout string) string {
	return formatTime(timeText, targetLayout)
}

func (TimeNS) D延时(ms int) {
	sleepMs(ms)
}

