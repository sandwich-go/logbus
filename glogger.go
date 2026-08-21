package logbus

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewGLogger(stdLogger *StdLogger, channelKey string, printAsError bool) *GLogger {
	g := &GLogger{
		stdLogger:    stdLogger,
		channelKey:   channelKey,
		printAsError: printAsError,
	}
	PushGLogger(g)
	return g
}

// GLogger 全局对象的类型定义
type GLogger struct {
	channelKey   string
	printAsError bool
	stdLogger    *StdLogger // 对外隐藏StdLogger
	depthLogger  sync.Map
}

type GLoggerVisitor interface {
	GetChannelKey() string
	GetPrintAsError() bool
	GetStdLogger() *StdLogger
}

func (s *GLogger) GetChannelKey() string {
	return s.channelKey
}
func (s *GLogger) GetPrintAsError() bool {
	return s.printAsError
}
func (s *GLogger) GetStdLogger() *StdLogger {
	return s.stdLogger
}

func (s *GLogger) printAsErr(fields ...Field) bool {
	hasErr := false
	for _, v := range fields {
		if v.Type == zapcore.ErrorType {
			hasErr = true
			break
		}
	}
	return hasErr
}

func (s *GLogger) syncDepthLogger() {
	s.depthLogger.Range(func(key, value interface{}) bool {
		if l, ok := value.(*zap.Logger); ok {
			_ = l.Sync()
		}
		return true
	})
}

func (s *GLogger) Debug(msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.DebugLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) DebugWithContext(ctx context.Context, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.DebugLevel) {
		return
	}
	fields = append(fields, FromContext(ctx))
	s.Debug(msg, fields...)
}

func (s *GLogger) Info(msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.InfoLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.InfoWithChannel(s.channelKey, fields...)
}

func (s *GLogger) InfoWithContext(ctx context.Context, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.InfoLevel) {
		return
	}
	fields = append(fields, FromContext(ctx))
	s.Info(msg, fields...)
}

func (s *GLogger) Warn(msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.WarnLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.WarnWithChannel(s.channelKey, fields...)
}

func (s *GLogger) WarnWithContext(ctx context.Context, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会让调用方传入的
	// 切片扩容重分配，且该代价在级别关闭时同样要付（zap 的级别判定在更下层）。
	if !s.printAsError && !s.stdLogger.Enabled(zap.WarnLevel) {
		return
	}
	fields = append(fields, FromContext(ctx))
	s.Warn(msg, fields...)
}

func (s *GLogger) Error(msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
}

func (s *GLogger) ErrorWithContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx))
	s.Error(msg, fields...)
}

func (s *GLogger) DPanic(msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.DPanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) DPanicWithContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx))
	s.DPanic(msg, fields...)
}

func (s *GLogger) Panic(msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.PanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) PanicWithContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx))
	s.Panic(msg, fields...)
}

func (s *GLogger) Fatal(msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.FatalWithChannel(s.channelKey, fields...)
}

func (s *GLogger) FatalWithContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx))
	s.Fatal(msg, fields...)
}

func (s *GLogger) getDepthLogger(depth int) *StdLogger {
	if lg, ok := s.depthLogger.Load(depth); ok {
		return lg.(*StdLogger)
	}
	cloneLogger := newStdLogger(s.stdLogger.WithOptions(zap.AddCallerSkip(depth)), s.stdLogger.fetch)
	s.depthLogger.Store(depth, cloneLogger)
	return cloneLogger
}

func (s *GLogger) GDebugDepth(depth int, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会触发扩容重分配，
	// getDepthLogger 未命中时还会 WithOptions 克隆一个 zap logger。
	// depth logger 由 s.stdLogger 派生，core 相同，故按 s.stdLogger 判级别即可。
	if !s.printAsError && !s.stdLogger.Enabled(zap.DebugLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) GInfoDepth(depth int, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会触发扩容重分配，
	// getDepthLogger 未命中时还会 WithOptions 克隆一个 zap logger。
	// depth logger 由 s.stdLogger 派生，core 相同，故按 s.stdLogger 判级别即可。
	if !s.printAsError && !s.stdLogger.Enabled(zap.InfoLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.InfoWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GWarnDepth(depth int, msg string, fields ...Field) {
	// 级别关闭且不会被 printAsError 改道时直接返回：下面的 append 会触发扩容重分配，
	// getDepthLogger 未命中时还会 WithOptions 克隆一个 zap logger。
	// depth logger 由 s.stdLogger 派生，core 相同，故按 s.stdLogger 判级别即可。
	if !s.printAsError && !s.stdLogger.Enabled(zap.WarnLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.WarnWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GErrorDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.ErrorWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GFatalDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	lg.FatalWithChannel(s.channelKey, fields...)
}

// WithChannel
func (s *GLogger) DebugWithChannel(c string, msg string, fields ...Field) {
	// 级别判定与写入都必须落在本实例的 stdLogger 上：SetGlobalGLogger 可以注入一个级别与
	// 全局默认 logger 不同的自定义 StdLogger，用 gStdLogger 判定会在两个方向上都出错——
	// 实例开启而全局关闭时丢掉本该输出的日志，实例关闭而全局开启时既绕过短路又写错目标。
	// 同级的 InfoWithChannel / WarnWithChannel / ErrorWithChannel 用的都是 s.stdLogger。
	if !s.stdLogger.Enabled(zap.DebugLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.DebugWithChannel(c, fields...)
}

func (s *GLogger) InfoWithChannel(c string, msg string, fields ...Field) {
	// 级别关闭时直接返回，避免下面的 append 触发扩容重分配。
	if !s.stdLogger.Enabled(zap.InfoLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.InfoWithChannel(c, fields...)
}

func (s *GLogger) WarnWithChannel(c string, msg string, fields ...Field) {
	// 级别关闭时直接返回，避免下面的 append 触发扩容重分配。
	if !s.stdLogger.Enabled(zap.WarnLevel) {
		return
	}
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.WarnWithChannel(c, fields...)
}

func (s *GLogger) ErrorWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.ErrorWithChannel(c, fields...)
}

func (s *GLogger) DPanicWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.DPanicWithChannel(c, fields...)
}

func (s *GLogger) PanicWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.PanicWithChannel(c, fields...)
}

func (s *GLogger) FatalWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String(MsgBody, msg))
	s.stdLogger.FatalWithChannel(c, fields...)
}
func (s *GLogger) Sync() error {
	return s.stdLogger.Sync()
}
