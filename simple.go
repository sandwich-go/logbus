package logbus

import "go.uber.org/zap"

// default logger provided
var newGlobalGLogger NewILogger

func Debug(msg string, v ...Field) {
	newGlobalGLogger.Debug(msg, v...)
}

func Info(msg string, v ...Field) {
	newGlobalGLogger.Info(msg, v...)
}

func Warn(msg string, v ...Field) {
	newGlobalGLogger.Warn(msg, v...)
}

func Error(msg string, v ...Field) {
	newGlobalGLogger.Error(msg, v...)
}

func Panic(msg string, v ...Field) {
	newGlobalGLogger.Panic(msg, v...)
}

func Fatal(msg string, v ...Field) {
	newGlobalGLogger.Fatal(msg, v...)
}

// WithChannel
func DebugWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.DebugWithChannel(c, msg, fields...)
}
func InfoWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.InfoWithChannel(c, msg, fields...)
}
func WarnWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.WarnWithChannel(c, msg, fields...)
}
func ErrorWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.ErrorWithChannel(c, msg, fields...)
}
func DPanicWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.DPanicWithChannel(c, msg, fields...)
}
func PanicWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.PanicWithChannel(c, msg, fields...)
}
func FatalWithChannel(c string, msg string, fields ...Field) {
	newGlobalGLogger.FatalWithChannel(c, msg, fields...)
}

// DebugDepth 用于glog被再封装
func DebugDepth(depth int, msg string, v ...Field) {
	newGlobalGLogger.GDebugDepth(depth, msg, v...)
}

// InfoDepth 用于glog被再封装
func InfoDepth(depth int, msg string, v ...Field) {
	newGlobalGLogger.GInfoDepth(depth, msg, v...)
}

// WarnDepth 用于glog被再封装
func WarnDepth(depth int, msg string, v ...Field) {
	newGlobalGLogger.GWarnDepth(depth, msg, v...)
}

// ErrorDepth 用于glog被再封装
func ErrorDepth(depth int, msg string, v ...Field) {
	newGlobalGLogger.GErrorDepth(depth, msg, v...)
}

// FatalDepth 用于glog被再封装
func FatalDepth(depth int, msg string, v ...Field) {
	newGlobalGLogger.GFatalDepth(depth, msg, v...)
}

// GetZapLogger
func GetZapLogger() *zap.Logger {
	from, ok := newGlobalGLogger.(GLoggerVisitor)
	if !ok {
		return nil
	}
	return from.GetStdLogger().getZapLogger()
}
