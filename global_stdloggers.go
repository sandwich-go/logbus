package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gStdLogger     *StdLogger
	gMonitorLogger *StdLogger
)

func initGlobalStdLoggers() {
	gStdLogger = newNLoggerInstance(Setting.DefaultTag)
	gMonitorLogger = newNLoggerInstance(MonitorTag)
}

func NewScopeLogger(tagName string, fields ...zap.Field) NewILogger {
	from, ok := newGlobalGLogger.(GLoggerVisitor)
	if !ok {
		return nil
	}
	return &GLogger{
		channelKey:   from.GetChannelKey(),
		printAsError: from.GetPrintAsError(),
		stdLogger:    newNLoggerInstance(tagName, fields...),
	}
}

func newNLoggerInstance(tagName string, fields ...zap.Field) *StdLogger {
	if tagName == "" {
		tagName = Setting.DefaultTag
	}

	var cores []zapcore.Core
	encoder := newJSONEncoder(EncodeConfig)
	if Setting.Dev {
		encoder = newConsoleEncoder(EncodeConfig)
	}

	// stdout 只能输出到stdout
	//var writer zapcore.WriteSyncer
	//writer = os.Stdout
	var writer = Setting.WriteSyncer
	if Setting.BufferedStdout {
		BufferedWriteSyncer.WS = Setting.WriteSyncer
		writer = BufferedWriteSyncer
	}

	stdCore := zapcore.NewCore(encoder, writer, Setting.LogLevel).With(append([]zap.Field{zap.String(Tags, tagName)}, fields...))
	cores = append(cores, stdCore)

	return &StdLogger{
		z: gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		})),
	}
}
