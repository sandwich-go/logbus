package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	gStdLogger     *StdLogger
	gMonitorLogger *StdLogger
)

func initGlobalStdLoggers() {
	initGStdLogger()
	initGMonitorLogger()
}

func initGStdLogger() {
	var cores []zapcore.Core
	encoder := newJSONEncoder(EncodeConfig)
	if Setting.Dev {
		encoder = newConsoleEncoder(EncodeConfig)
	}

	// stdout 只能输出到stdout
	var writer zapcore.WriteSyncer
	writer = os.Stdout
	if Setting.BufferedStdout {
		writer = BufferedWriteSyncer
	}
	stdCore := zapcore.NewCore(encoder, writer, Setting.LogLevel).With([]zap.Field{zap.String(Tags, Setting.DefaultTag)})
	cores = append(cores, stdCore)

	gStdLogger = &StdLogger{
		z: gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		})),
		//tags: []string{Setting.DefaultTag},
	}
}

func initGMonitorLogger() {
	var cores []zapcore.Core
	encoder := newJSONEncoder(EncodeConfig)
	if Setting.Dev {
		encoder = newConsoleEncoder(EncodeConfig)
	}

	// stdout 只能输出到stdout
	var writer zapcore.WriteSyncer
	writer = os.Stdout
	if Setting.BufferedStdout {
		writer = BufferedWriteSyncer
	}
	stdCore := zapcore.NewCore(encoder, writer, Setting.LogLevel).With([]zap.Field{zap.String(Tags, MonitorTag)})
	cores = append(cores, stdCore)

	gMonitorLogger = &StdLogger{
		z: gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		})),
		//tags: []string{MonitorTag},
	}
}
