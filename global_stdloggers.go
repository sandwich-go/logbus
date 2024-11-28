package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gStdLogger     *StdLogger
	gMonitorLogger *StdLogger
	cacheGLogger   []GLoggerCache
)

type GLoggerCache struct {
	tagName string
	fields  []zap.Field
	gLogger *GLogger
}

func initGlobalStdLoggers() {
	gStdLogger = newNLoggerInstance(Setting.DefaultTag)
	gMonitorLogger = newNLoggerInstance(MonitorTag)
}

func NewScopeLogger(tagName string, fields ...zap.Field) NewILogger {
	if refresh {
		lg := &GLogger{
			stdLogger: newNLoggerInstance(tagName, fields...),
		}
		cacheGLogger = append(cacheGLogger, GLoggerCache{
			tagName: tagName,
			fields:  fields,
			gLogger: lg,
		})
		return lg
	}
	from, ok := newGlobalGLogger.(GLoggerVisitor)
	if !ok {
		return nil
	}

	newStdLogger := newNLoggerInstance(tagName, fields...)
	// 比 newGlobalGLogger 少一层调用
	newStdLogger.z = newStdLogger.z.WithOptions(zap.AddCallerSkip(-1))

	return &GLogger{
		channelKey:   from.GetChannelKey(),
		printAsError: from.GetPrintAsError(),
		stdLogger:    newStdLogger,
	}
}

// NewScopeLoggerWithFetchFunc 用于创建自定义 fetch FetchLogContext 的logger
func NewScopeLoggerWithFetchFunc(tagName string, fetch FetchLogContext, fields ...zap.Field) NewILogger {
	lg := NewScopeLogger(tagName, fields...)
	lg.(*GLogger).stdLogger.fetch = fetch
	return lg
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

	stdCore := zapcore.NewCore(encoder, writer, NewTraceLevelEnabler(Setting.LogLevel)).With(append([]zap.Field{zap.String(Tags, tagName)}, fields...))
	cores = append(cores, stdCore)

	return &StdLogger{
		z: gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		})),
	}
}

// refreshEarlyLogger 用于刷新在Init() 之前创建的logger
func refreshEarlyLogger() {
	for _, v := range cacheGLogger {
		newStdLogger := NewScopeLogger(v.tagName, v.fields...)
		v.gLogger.channelKey = newStdLogger.(*GLogger).channelKey
		v.gLogger.printAsError = newStdLogger.(*GLogger).printAsError
		oldFetch := v.gLogger.stdLogger.fetch
		v.gLogger.stdLogger = newStdLogger.(*GLogger).stdLogger
		if oldFetch != nil {
			v.gLogger.stdLogger.fetch = oldFetch
		}
	}
	cacheGLogger = nil
}
