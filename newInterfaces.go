package logbus

import (
	"go.uber.org/zap"
)

type NewILogger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	DebugWithChannel(c string, msg string, fields ...zap.Field)
	InfoWithChannel(c string, msg string, fields ...zap.Field)
	WarnWithChannel(c string, msg string, fields ...zap.Field)
	ErrorWithChannel(c string, msg string, fields ...zap.Field)
	DPanicWithChannel(c string, msg string, fields ...zap.Field)
	PanicWithChannel(c string, msg string, fields ...zap.Field)
	FatalWithChannel(c string, msg string, fields ...zap.Field)

	GDebugDepth(depth int, msg string, v ...zap.Field)
	GInfoDepth(depth int, msg string, v ...zap.Field)
	GWarnDepth(depth int, msg string, v ...zap.Field)
	GErrorDepth(depth int, msg string, v ...zap.Field)
	GFatalDepth(depth int, msg string, v ...zap.Field)
}
