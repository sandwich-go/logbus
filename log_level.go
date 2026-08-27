package logbus

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

var runtimeLogLevelMu sync.Mutex

// SetLogLevel 在运行中动态调整日志级别，并立即作用于全局 logger 与已创建的 scope logger。
//
// runtimeLogLevelMu 只负责让多个写者的两步写（Setting.LogLevel 镜像 + runtimeLogLevel）
// 不互相交错，读侧不需要持锁。Setting.LogLevel 是给外部 Conf.GetLogLevel() 读的镜像，
// 级别判定唯一依据是 runtimeLogLevel。
func SetLogLevel(level zapcore.Level) {
	runtimeLogLevelMu.Lock()
	defer runtimeLogLevelMu.Unlock()
	Setting.LogLevel = level
	runtimeLogLevel.SetLevel(level)
}

// GetLogLevel 返回当前运行时日志级别。
//
// 不加 runtimeLogLevelMu：zap.AtomicLevel 内部是 *atomic.Int32，Level() 本身就是原子读，
// 且这里没有第二份状态需要与它成对观察（Setting.LogLevel 只是给外部读的镜像，
// 不参与级别判定）。加锁不增加任何保证，只会让本函数无法用在热路径上。
func GetLogLevel() zapcore.Level {
	return runtimeLogLevel.Level()
}
