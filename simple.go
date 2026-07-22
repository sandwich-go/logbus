package logbus

import (
	"context"

	"go.uber.org/zap"
)

// default logger provided
var newGlobalGLogger NewILogger

func Debug(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Debug(ctx, msg, v...)
}

func Info(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Info(ctx, msg, v...)
}

func Warn(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Warn(ctx, msg, v...)
}

func Error(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Error(ctx, msg, v...)
}

func DPanic(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.DPanic(ctx, msg, v...)
}

func Panic(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Panic(ctx, msg, v...)
}

func Fatal(ctx context.Context, msg string, v ...Field) {
	newGlobalGLogger.Fatal(ctx, msg, v...)
}

// WithChannel
func DebugWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.DebugWithChannel(ctx, c, msg, fields...)
}
func InfoWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.InfoWithChannel(ctx, c, msg, fields...)
}
func WarnWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.WarnWithChannel(ctx, c, msg, fields...)
}
func ErrorWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.ErrorWithChannel(ctx, c, msg, fields...)
}
func DPanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.DPanicWithChannel(ctx, c, msg, fields...)
}
func PanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.PanicWithChannel(ctx, c, msg, fields...)
}
func FatalWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	newGlobalGLogger.FatalWithChannel(ctx, c, msg, fields...)
}

// DebugDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func DebugDepth(ctx context.Context, depth int, msg string, v ...Field) {
	newGlobalGLogger.GDebugDepth(ctx, depth, msg, v...)
}

// InfoDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func InfoDepth(ctx context.Context, depth int, msg string, v ...Field) {
	newGlobalGLogger.GInfoDepth(ctx, depth, msg, v...)
}

// WarnDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func WarnDepth(ctx context.Context, depth int, msg string, v ...Field) {
	newGlobalGLogger.GWarnDepth(ctx, depth, msg, v...)
}

// ErrorDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func ErrorDepth(ctx context.Context, depth int, msg string, v ...Field) {
	newGlobalGLogger.GErrorDepth(ctx, depth, msg, v...)
}

// FatalDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func FatalDepth(ctx context.Context, depth int, msg string, v ...Field) {
	newGlobalGLogger.GFatalDepth(ctx, depth, msg, v...)
}

func DefaultLogger() NewILogger {
	return newGlobalGLogger
}

// GetZapLogger 返回zap.Logger
func GetZapLogger() *zap.Logger {
	from, ok := newGlobalGLogger.(GLoggerVisitor)
	if !ok {
		return nil
	}
	return from.GetStdLogger().getZapLogger()
}
