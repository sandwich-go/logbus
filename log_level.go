package logbus

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

var runtimeLogLevelMu sync.Mutex

// SetLogLevel 在运行中动态调整日志级别，并立即作用于全局 logger 与已创建的 scope logger。
func SetLogLevel(level zapcore.Level) {
	runtimeLogLevelMu.Lock()
	defer runtimeLogLevelMu.Unlock()
	Setting.LogLevel = level
	runtimeLogLevel.SetLevel(level)
}
