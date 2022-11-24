package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// GLogger 全局对象的类型定义
type GLogger struct {
	channelKey   string
	printAsError bool
	stdLogger    *StdLogger // 对外隐藏StdLogger
	depthLogger  sync.Map
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
	fields = append(fields, String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Info(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.InfoWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Warn(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.WarnWithChannel(s.channelKey, fields...)
}
func (s *GLogger) Error(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
}

func (s *GLogger) DPanic(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.DPanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Panic(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.PanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Fatal(msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.FatalWithChannel(s.channelKey, fields...)
}

func (s *GLogger) getDepthLogger(depth int) *zap.Logger {
	if lg, ok := s.depthLogger.Load(depth); ok {
		return lg.(*zap.Logger)
	}
	cloneLogger := s.stdLogger.WithOptions(zap.AddCallerSkip(depth))
	s.depthLogger.Store(depth, cloneLogger)
	return cloneLogger
}

func (s *GLogger) GDebugDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Debug(s.channelKey, fields...)
}

func (s *GLogger) GInfoDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Info(s.channelKey, fields...)
}
func (s *GLogger) GWarnDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Warn(s.channelKey, fields...)
}
func (s *GLogger) GErrorDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Error(s.channelKey, fields...)
}
func (s *GLogger) GFatalDepth(depth int, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	lg.Fatal(s.channelKey, fields...)
}

// WithChannel
func (s *GLogger) DebugWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	gStdLogger.DebugWithChannel(c, fields...)
}

func (s *GLogger) InfoWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.InfoWithChannel(c, fields...)
}

func (s *GLogger) WarnWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.WarnWithChannel(c, fields...)
}

func (s *GLogger) ErrorWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.ErrorWithChannel(c, fields...)
}

func (s *GLogger) DPanicWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.DPanicWithChannel(c, fields...)
}

func (s *GLogger) PanicWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.PanicWithChannel(c, fields...)
}

func (s *GLogger) FatalWithChannel(c string, msg string, fields ...Field) {
	fields = append(fields, String("glog-msg", msg))
	s.stdLogger.FatalWithChannel(c, fields...)
}
