package logbus

import (
	"context"
	"sync/atomic"

	"go.uber.org/zap"
)

// default logger provided
var globalGLogger atomic.Pointer[GLogger]

func defaultGLogger() *GLogger {
	return globalGLogger.Load()
}

func Debug(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Debug(ctx, msg, v...)
}

func Info(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Info(ctx, msg, v...)
}

func Warn(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Warn(ctx, msg, v...)
}

func Error(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Error(ctx, msg, v...)
}

func DPanic(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().DPanic(ctx, msg, v...)
}

func Panic(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Panic(ctx, msg, v...)
}

func Fatal(ctx context.Context, msg string, v ...Field) {
	defaultGLogger().Fatal(ctx, msg, v...)
}

// WithChannel
func DebugWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().DebugWithChannel(ctx, c, msg, fields...)
}
func InfoWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().InfoWithChannel(ctx, c, msg, fields...)
}
func WarnWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().WarnWithChannel(ctx, c, msg, fields...)
}
func ErrorWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().ErrorWithChannel(ctx, c, msg, fields...)
}
func DPanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().DPanicWithChannel(ctx, c, msg, fields...)
}
func PanicWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().PanicWithChannel(ctx, c, msg, fields...)
}
func FatalWithChannel(ctx context.Context, c string, msg string, fields ...Field) {
	defaultGLogger().FatalWithChannel(ctx, c, msg, fields...)
}

// DebugDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func DebugDepth(ctx context.Context, depth int, msg string, v ...Field) {
	defaultGLogger().GDebugDepth(ctx, depth, msg, v...)
}

// InfoDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func InfoDepth(ctx context.Context, depth int, msg string, v ...Field) {
	defaultGLogger().GInfoDepth(ctx, depth, msg, v...)
}

// WarnDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func WarnDepth(ctx context.Context, depth int, msg string, v ...Field) {
	defaultGLogger().GWarnDepth(ctx, depth, msg, v...)
}

// ErrorDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func ErrorDepth(ctx context.Context, depth int, msg string, v ...Field) {
	defaultGLogger().GErrorDepth(ctx, depth, msg, v...)
}

// FatalDepth 用于 glog 被再封装，并从 ctx 提取日志元数据。
func FatalDepth(ctx context.Context, depth int, msg string, v ...Field) {
	defaultGLogger().GFatalDepth(ctx, depth, msg, v...)
}

func DefaultLogger() NewILogger {
	return defaultGLogger()
}

// GetZapLogger 返回zap.Logger
func GetZapLogger() *zap.Logger {
	logger := defaultGLogger()
	if logger == nil {
		return nil
	}
	return logger.GetStdLogger().getZapLogger()
}
