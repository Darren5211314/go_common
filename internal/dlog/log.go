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
	once   sync.Once
	sugar *zap.SugaredLogger
)

// 初始化 显式触发初始化（可重复调用，安全）。日常不必调用，首次打日志会自动初始化。
func init() {
	ensureInit()
}

func ensureInit() {
	once.Do(func() {
		encodeCfg := zap.NewDevelopmentEncoderConfig()
		encodeCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encodeCfg.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		}
		encodeCfg.EncodeCaller = zapcore.ShortCallerEncoder

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encodeCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		// CallerSkip(1)：跳过 darren.Info / dlog.Info 这一层封装，行号落在业务调用处。
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		zap.ReplaceGlobals(logger)
		sugar = logger.Sugar()
	})
}

// Sugar 返回已初始化的 SugaredLogger（带 CallerSkip）。
func Sugar() *zap.SugaredLogger {
	ensureInit()
	return sugar
}
