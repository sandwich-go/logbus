package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type StdLogger struct {
	fetch FetchLogContext
	z     *zap.Logger
}

func (s *StdLogger) WithOptions(opts ...zap.Option) *zap.Logger {
	return s.z.WithOptions(opts...)
}

func (s *StdLogger) SetZLogger(zl *zap.Logger) {
	if zl != nil {
		s.z = zl
	}
}

func (s *StdLogger) Sync() error {
	return s.z.Sync()
}

func (s *StdLogger) getZapLogger() *zap.Logger {
	return s.z
}

func newStdLogger(z *zap.Logger, fetch FetchLogContext) *StdLogger {
	s := &StdLogger{
		fetch: fetch,
		z:     z,
	}
	return s
}

// Enabled 报告该级别的日志是否会被输出。
//
// 判定直接问 zapcore.LevelEnabler，走的是 zap.AtomicLevel 与 atomic.Bool 的原子读
// （见 trackLevelEnabler.Enabled），没有锁，可以放在热路径上。这里不用 GetLogLevel：
// 它只看全局级别，看不到 trackLevelEnabler 对 TrackLevel 的放行，也看不到 core 上的其它过滤。
//
// 用途是让调用方在**构造字段之前**短路。zap 自身的级别判定发生在 z.Debug/z.Info 内部，
// 而 StdLogger.fields 会先调用 FetchLogContext 钩子并分配新切片、GLogger 的各级别方法会先
// append msg 字段，这些代价在级别关闭时同样要付。
func (s *StdLogger) Enabled(lvl zapcore.Level) bool {
	if s == nil || s.z == nil {
		return false
	}
	return s.z.Core().Enabled(lvl)
}
