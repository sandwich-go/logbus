package logbus

import (
	"go.uber.org/zap"
)

// Debug log
func Debug(fields ...zap.Field) {
	Logger().Debug(fields...)
}

func DebugWithChannel(c string, fields ...zap.Field) {
	Logger().DebugWithChannel(c, fields...)
}

// Info log
func Info(fields ...zap.Field) {
	Logger().Info(fields...)
}

func InfoWithChannel(c string, fields ...zap.Field) {
	Logger().InfoWithChannel(c, fields...)
}

// Warn log
func Warn(fields ...zap.Field) {
	Logger().Warn(fields...)
}

func WarnWithChannel(c string, fields ...zap.Field) {
	Logger().WarnWithChannel(c, fields...)
}

// Error log
func Error(fields ...zap.Field) {
	Logger().Error(fields...)
}

func ErrorWithChannel(c string, fields ...zap.Field) {
	Logger().ErrorWithChannel(c, fields...)
}

// DPanic log
func DPanic(fields ...zap.Field) {
	Logger().DPanic(fields...)
}

func DPanicWithChannel(c string, fields ...zap.Field) {
	Logger().DPanicWithChannel(c, fields...)
}

// Panic log
func Panic(fields ...zap.Field) {
	Logger().Panic(fields...)
}

func PanicWithChannel(c string, fields ...zap.Field) {
	Logger().PanicWithChannel(c, fields...)
}

// Fatal log
func Fatal(fields ...zap.Field) {
	Logger().Fatal(fields...)
}

func FatalWithChannel(c string, fields ...zap.Field) {
	Logger().FatalWithChannel(c, fields...)
}
