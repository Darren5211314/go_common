package darren

import "github.com/Darren5211314/go_common/internal/dlog"

// Debug 输出 Debug 日志（消息 + 交替键值对）。
func Debug(消息 string, 键值 ...any) {
	dlog.Sugar().Debugw(消息, 键值...)
}

// Info 输出 Info 日志（消息 + 交替键值对）。
func Info(消息 string, 键值 ...any) {
	dlog.Sugar().Infow(消息, 键值...)
}

// Warn 输出 Warn 日志（消息 + 交替键值对）。
func Warn(消息 string, 键值 ...any) {
	dlog.Sugar().Warnw(消息, 键值...)
}

// Error 输出 Error 日志（消息 + 交替键值对）。
func Error(消息 string, 键值 ...any) {
	dlog.Sugar().Errorw(消息, 键值...)
}

// DPanic 输出 DPanic 日志（开发环境随后 panic；生产等同 Error）。
func DPanic(消息 string, 键值 ...any) {
	dlog.Sugar().DPanicw(消息, 键值...)
}

// Panic 输出 Panic 日志后 panic。
func Panic(消息 string, 键值 ...any) {
	dlog.Sugar().Panicw(消息, 键值...)
}

// Fatal 输出 Fatal 日志后 os.Exit(1)。
func Fatal(消息 string, 键值 ...any) {
	dlog.Sugar().Fatalw(消息, 键值...)
}

// Debugf 以 printf 风格输出 Debug。
func Debugf(格式 string, 参数 ...any) {
	dlog.Sugar().Debugf(格式, 参数...)
}

// Infof 以 printf 风格输出 Info。
func Infof(格式 string, 参数 ...any) {
	dlog.Sugar().Infof(格式, 参数...)
}

// Warnf 以 printf 风格输出 Warn。
func Warnf(格式 string, 参数 ...any) {
	dlog.Sugar().Warnf(格式, 参数...)
}

// Errorf 以 printf 风格输出 Error。
func Errorf(格式 string, 参数 ...any) {
	dlog.Sugar().Errorf(格式, 参数...)
}

// DPanicf 以 printf 风格输出 DPanic。
func DPanicf(格式 string, 参数 ...any) {
	dlog.Sugar().DPanicf(格式, 参数...)
}

// Panicf 以 printf 风格输出 Panic 后 panic。
func Panicf(格式 string, 参数 ...any) {
	dlog.Sugar().Panicf(格式, 参数...)
}

// Fatalf 以 printf 风格输出 Fatal 后 os.Exit(1)。
func Fatalf(格式 string, 参数 ...any) {
	dlog.Sugar().Fatalf(格式, 参数...)
}
