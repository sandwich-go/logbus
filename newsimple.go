package logbus

import (
	"go.uber.org/zap"
)

// default logger provided
var newGlobalGLogger NewILogger //= NewDefaultLogger(log.New(os.Stdout, "", log.LstdFlags|log.Llongfile))

func Debug(msg string, v ...zap.Field) {
	newGlobalGLogger.Debug(msg, v...)
}

func Info(msg string, v ...zap.Field) {
	newGlobalGLogger.Info(msg, v...)
}

func Warn(msg string, v ...zap.Field) {
	newGlobalGLogger.Warn(msg, v...)
}

func Error(msg string, v ...zap.Field) {
	newGlobalGLogger.Error(msg, v...)
}

func Panic(msg string, v ...zap.Field) {
	newGlobalGLogger.Panic(msg, v...)
}

func Fatal(msg string, v ...zap.Field) {
	newGlobalGLogger.Fatal(msg, v...)
}

// WithChannel
func DebugWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.DebugWithChannel(c, msg, fields...)
}
func InfoWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.InfoWithChannel(c, msg, fields...)
}
func WarnWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.WarnWithChannel(c, msg, fields...)
}
func ErrorWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.ErrorWithChannel(c, msg, fields...)
}
func DPanicWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.DPanicWithChannel(c, msg, fields...)
}
func PanicWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.PanicWithChannel(c, msg, fields...)
}
func FatalWithChannel(c string, msg string, fields ...zap.Field) {
	newGlobalGLogger.FatalWithChannel(c, msg, fields...)
}

// DebugDepth 用于glog被再封装
func DebugDepth(depth int, msg string, v ...zap.Field) {
	newGlobalGLogger.GDebugDepth(depth, msg, v...)
}

// InfoDepth 用于glog被再封装
func InfoDepth(depth int, msg string, v ...zap.Field) {
	newGlobalGLogger.GInfoDepth(depth, msg, v...)
}

// WarnDepth 用于glog被再封装
func WarnDepth(depth int, msg string, v ...zap.Field) {
	newGlobalGLogger.GWarnDepth(depth, msg, v...)
}

// ErrorDepth 用于glog被再封装
func ErrorDepth(depth int, msg string, v ...zap.Field) {
	newGlobalGLogger.GErrorDepth(depth, msg, v...)
}

// FatalDepth 用于glog被再封装
func FatalDepth(depth int, msg string, v ...zap.Field) {
	newGlobalGLogger.GFatalDepth(depth, msg, v...)
}
