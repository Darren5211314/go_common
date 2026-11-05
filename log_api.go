package darren

import "github.com/Darren5211314/go_common/internal/dlog"

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
