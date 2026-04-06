package logbus

import (
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xpanic"
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

// newScopeLoggerInternal 内部实现，printAsErrorWhenRefresh 为 refresh 时的 printAsError；
// usePrintAsErrorFromVisitor 为 true 时在非 refresh 分支使用 from.GetPrintAsError()，否则使用 printAsErrorWhenRefresh。
func newScopeLoggerInternal(tagName string, printAsErrorWhenRefresh bool, usePrintAsErrorFromVisitor bool, fields ...zap.Field) NewILogger {
	if refresh {
		lg := NewGLogger(newNLoggerInstance(tagName, fields...), "", printAsErrorWhenRefresh)
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
	newStdLogger.z = newStdLogger.z.WithOptions(zap.AddCallerSkip(-2))

	printAsError := printAsErrorWhenRefresh
	if usePrintAsErrorFromVisitor {
		printAsError = from.GetPrintAsError()
	}
	return NewGLogger(newStdLogger, from.GetChannelKey(), printAsError)
}

func NewScopeLogger(tagName string, fields ...zap.Field) NewILogger {
	return newScopeLoggerInternal(tagName, false, true, fields...)
}

func NewScopeLoggerPrintAsError(tagName string, fields ...zap.Field) NewILogger {
	return newScopeLoggerInternal(tagName, true, false, fields...)
}

// NewScopeLoggerWithFetchFunc 用于创建自定义 fetch FetchLogContext 的logger
func NewScopeLoggerWithFetchFunc(tagName string, fetch FetchLogContext, fields ...zap.Field) NewILogger {
	lg := NewScopeLogger(tagName, fields...)
	lg.(*GLogger).stdLogger.fetch = fetch
	return lg
}

// CoreWrapper 根据tag name对core进行一层封装
var CoreWrapper = func(tagName string, core zapcore.Core) zapcore.Core {
	return core
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
	if !Setting.DisableTruncateWriteSyncer {
		xpanic.WhenTrue(Setting.TruncateWriteSyncerOption == nil, "TruncateWriteSyncerOption cannot be nil when TruncateWriteSyncer is enabled")
		writer = NewTruncateWriteSyncer(writer, Setting.TruncateWriteSyncerOption)
	}

	stdCore := zapcore.NewCore(encoder, writer, newTrackLevelEnabler()).With(append([]zap.Field{zap.String(Tags, tagName)}, fields...))
	stdCore = CoreWrapper(tagName, stdCore)
	if xos.EnvGetCaseInsensitive("logbus_core_dup") != "" {
		stdCore = NewDupCore(stdCore)
	}
	cores = append(cores, stdCore)

	s := newStdLogger(gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})), nil)
	return s
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
