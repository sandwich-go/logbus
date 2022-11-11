package logbus

import "go.uber.org/zap"

type ILogger interface {
	Debug(fields ...zap.Field)
	DebugWithChannel(c string, fields ...zap.Field)
	Info(fields ...zap.Field)
	InfoWithChannel(c string, fields ...zap.Field)
	Warn(fields ...zap.Field)
	WarnWithChannel(c string, fields ...zap.Field)
	Error(fields ...zap.Field)
	ErrorWithChannel(c string, fields ...zap.Field)
	DPanic(fields ...zap.Field)
	DPanicWithChannel(c string, fields ...zap.Field)
	Panic(fields ...zap.Field)
	PanicWithChannel(c string, fields ...zap.Field)
	Fatal(fields ...zap.Field)
	FatalWithChannel(c string, fields ...zap.Field)
}

type ITracker interface {
	Track(...zap.Field) error
}
