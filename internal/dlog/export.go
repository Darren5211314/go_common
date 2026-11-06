package dlog

// Info 输出 Info 日志（消息 + 交替键值对）。
func Info(msg string, kv ...any) {
	Sugar().Infow(msg, kv...)
}

// Warn 输出 Warn 日志（消息 + 交替键值对）。
func Warn(msg string, kv ...any) {
	Sugar().Warnw(msg, kv...)
}

// Error 输出 Error 日志（消息 + 交替键值对）。
func Error(msg string, kv ...any) {
	Sugar().Errorw(msg, kv...)
}

// Infof 以 printf 风格输出 Info。
func Infof(format string, arg ...any) {
	Sugar().Infof(format, arg...)
}

// Warnf 以 printf 风格输出 Warn。
func Warnf(format string, arg ...any) {
	Sugar().Warnf(format, arg...)
}

// Errorf 以 printf 风格输出 Error。
func Errorf(format string, arg ...any) {
	Sugar().Errorf(format, arg...)
}
