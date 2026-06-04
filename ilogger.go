package logbus

import "context"

type NewILogger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	DebugWithContext(ctx context.Context, msg string, fields ...Field)
	InfoWithContext(ctx context.Context, msg string, fields ...Field)
	WarnWithContext(ctx context.Context, msg string, fields ...Field)
	ErrorWithContext(ctx context.Context, msg string, fields ...Field)
	DPanicWithContext(ctx context.Context, msg string, fields ...Field)
	PanicWithContext(ctx context.Context, msg string, fields ...Field)
	FatalWithContext(ctx context.Context, msg string, fields ...Field)

	DebugWithChannel(c string, msg string, fields ...Field)
	InfoWithChannel(c string, msg string, fields ...Field)
	WarnWithChannel(c string, msg string, fields ...Field)
	ErrorWithChannel(c string, msg string, fields ...Field)
	DPanicWithChannel(c string, msg string, fields ...Field)
	PanicWithChannel(c string, msg string, fields ...Field)
	FatalWithChannel(c string, msg string, fields ...Field)

	GDebugDepth(depth int, msg string, v ...Field)
	GInfoDepth(depth int, msg string, v ...Field)
	GWarnDepth(depth int, msg string, v ...Field)
	GErrorDepth(depth int, msg string, v ...Field)
	GFatalDepth(depth int, msg string, v ...Field)

	Sync() error

	syncDepthLogger()
}
