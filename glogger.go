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

func (s *GLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Info(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.InfoWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.WarnWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Error(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
}

func (s *GLogger) DPanic(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.DPanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Panic(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.PanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.FatalWithChannel(s.channelKey, fields...)
}

func (s *GLogger) getDepthLogger(depth int) *StdLogger {
	if lg, ok := s.depthLogger.Load(depth); ok {
		return lg.(*StdLogger)
	}
	cloneLogger := newStdLogger(s.stdLogger.WithOptions(zap.AddCallerSkip(depth)), s.stdLogger.config(), s.stdLogger.fetch)
	s.depthLogger.Store(depth, cloneLogger)
	return cloneLogger
}

func (s *GLogger) GDebugDepth(ctx context.Context, depth int, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) GInfoDepth(ctx context.Context, depth int, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.InfoWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GWarnDepth(ctx context.Context, depth int, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.WarnWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GErrorDepth(ctx context.Context, depth int, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	lg.ErrorWithChannel(s.channelKey, fields...)
}
func (s *GLogger) GFatalDepth(ctx context.Context, depth int, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	lg := s.getDepthLogger(depth)
	lg.FatalWithChannel(s.channelKey, fields...)
}

// WithChannel
func (s *GLogger) DebugWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.DebugWithChannel(c, fields...)
}

func (s *GLogger) InfoWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.InfoWithChannel(c, fields...)
}

func (s *GLogger) WarnWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.WarnWithChannel(c, fields...)
}

func (s *GLogger) ErrorWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.ErrorWithChannel(c, fields...)
}

func (s *GLogger) DPanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.DPanicWithChannel(c, fields...)
}

func (s *GLogger) PanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.PanicWithChannel(c, fields...)
}

func (s *GLogger) FatalWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	fields = append(fields, FromContext(ctx), String(MsgBody, msg))
	s.stdLogger.FatalWithChannel(c, fields...)
}
func (s *GLogger) Sync() error {
	return s.stdLogger.Sync()
}
