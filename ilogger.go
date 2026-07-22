package logbus

import "context"

type NewILogger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	DPanic(ctx context.Context, msg string, fields ...Field)
	Panic(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)

	DebugWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	InfoWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	WarnWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	ErrorWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	DPanicWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	PanicWithChannel(ctx context.Context, c string, msg string, fields ...Field)
	FatalWithChannel(ctx context.Context, c string, msg string, fields ...Field)

	GDebugDepth(ctx context.Context, depth int, msg string, v ...Field)
	GInfoDepth(ctx context.Context, depth int, msg string, v ...Field)
	GWarnDepth(ctx context.Context, depth int, msg string, v ...Field)
	GErrorDepth(ctx context.Context, depth int, msg string, v ...Field)
	GFatalDepth(ctx context.Context, depth int, msg string, v ...Field)

	Sync() error

	syncDepthLogger()
}
