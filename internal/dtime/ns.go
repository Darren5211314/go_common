package dtime

// TimeNS 命名空间类型（根包以未导出别名挂到 var Time）。
type TimeNS struct{}

func (TimeNS) Now(tenDigits bool) int64 {
	return nowTimestamp(tenDigits)
}

func (TimeNS) FormatUnix(ts int64, format string) string {
	return timestampFormat(ts, format)
}

func (TimeNS) ParseUnix(timeText string, tenDigits bool) int64 {
	return timeToTimestamp(timeText, tenDigits)
}

func (TimeNS) Format(timeText, targetLayout string) string {
	return formatTime(timeText, targetLayout)
}

func (TimeNS) Sleep(ms int) {
	sleepMs(ms)
}

