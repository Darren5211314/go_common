// Package dlog 基于 zap 的彩色控制台日志，首次使用时自动初始化。
package dlog

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	一次   sync.Once
	糖化日志 *zap.SugaredLogger
)

// 初始化 显式触发初始化（可重复调用，安全）。日常不必调用，首次打日志会自动初始化。
func 初始化() {
	确保已初始化()
}

func 确保已初始化() {
	一次.Do(func() {
		编码配置 := zap.NewDevelopmentEncoderConfig()
		编码配置.EncodeLevel = zapcore.CapitalColorLevelEncoder
		编码配置.EncodeTime = func(时刻 time.Time, 编码器 zapcore.PrimitiveArrayEncoder) {
			编码器.AppendString(时刻.Format("2006-01-02 15:04:05"))
		}
		编码配置.EncodeCaller = zapcore.ShortCallerEncoder

		核心 := zapcore.NewCore(
			zapcore.NewConsoleEncoder(编码配置),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		// CallerSkip(1)：跳过 darren.Info / dlog.Info 这一层封装，行号落在业务调用处。
		日志器 := zap.New(核心, zap.AddCaller(), zap.AddCallerSkip(1))
		zap.ReplaceGlobals(日志器)
		糖化日志 = 日志器.Sugar()
	})
}

// Sugar 返回已初始化的 SugaredLogger（带 CallerSkip）。
func Sugar() *zap.SugaredLogger {
	确保已初始化()
	return 糖化日志
}
