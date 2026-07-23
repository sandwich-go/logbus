package logbus

import (
	"sync/atomic"

	"github.com/sandwich-go/boost/xos"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gStdLogger     atomic.Pointer[StdLogger]
	gMonitorLogger *StdLogger
	cacheGLogger   []GLoggerCache
)

type GLoggerCache struct {
	tagName     string
	fields      []zap.Field
	gLogger     *GLogger
	printAsErr  bool
	inheritMode bool
	fetch       FetchLogContext
}

func initGlobalStdLoggers(config *configSnapshot) {
	gStdLogger.Store(newNLoggerInstance(config, config.DefaultTag))
	gMonitorLogger = newNLoggerInstance(config, MonitorTag)
}

// newScopeLoggerInternal 内部实现，printAsErrorWhenRefresh 为 refresh 时的 printAsError；
// usePrintAsErrorFromVisitor 为 true 时在非 refresh 分支使用 from.GetPrintAsError()，否则使用 printAsErrorWhenRefresh。
func newScopeLoggerInternal(tagName string, printAsErrorWhenRefresh bool, usePrintAsErrorFromVisitor bool, fetch FetchLogContext, fields ...zap.Field) NewILogger {
	config := currentConfig()
	if refresh {
		stdLogger := newNLoggerInstance(config, tagName, fields...)
		if fetch != nil {
			stdLogger.fetch = fetch
		}
		lg := NewGLogger(stdLogger, "", printAsErrorWhenRefresh)
		cacheGLogger = append(cacheGLogger, GLoggerCache{
			tagName:     tagName,
			fields:      fields,
			gLogger:     lg,
			printAsErr:  printAsErrorWhenRefresh,
			inheritMode: usePrintAsErrorFromVisitor,
			fetch:       fetch,
		})
		return lg
	}
	from := defaultGLogger()
	if from == nil {
		return nil
	}

	newStdLogger := newNLoggerInstance(config, tagName, fields...)
	// 比全局 package function 少一层调用
	newStdLogger.z = newStdLogger.z.WithOptions(zap.AddCallerSkip(-2))
	if fetch != nil {
		newStdLogger.fetch = fetch
	}

	printAsError := printAsErrorWhenRefresh
	if usePrintAsErrorFromVisitor {
		printAsError = from.GetPrintAsError()
	}
	return NewGLogger(newStdLogger, from.GetChannelKey(), printAsError)
}

func NewScopeLogger(tagName string, fields ...zap.Field) NewILogger {
	return newScopeLoggerInternal(tagName, false, true, nil, fields...)
}

func NewScopeLoggerPrintAsError(tagName string, fields ...zap.Field) NewILogger {
	return newScopeLoggerInternal(tagName, true, false, nil, fields...)
}

// NewScopeLoggerWithFetchFunc 用于创建自定义 fetch FetchLogContext 的logger
func NewScopeLoggerWithFetchFunc(tagName string, fetch FetchLogContext, fields ...zap.Field) NewILogger {
	return newScopeLoggerInternal(tagName, false, true, fetch, fields...)
}

// CoreWrapper 根据tag name对core进行一层封装
var CoreWrapper = func(tagName string, core zapcore.Core) zapcore.Core {
	return core
}

func newNLoggerInstance(config *configSnapshot, tagName string, fields ...zap.Field) *StdLogger {
	if tagName == "" {
		tagName = config.DefaultTag
	}

	var cores []zapcore.Core
	encoder := newJSONEncoder(config.encoder)
	if config.Dev {
		encoder = newConsoleEncoder(config.encoder)
	}

	stdCore := zapcore.NewCore(encoder, config.writer, newTrackLevelEnabler()).With(append([]zap.Field{zap.String(Tags, tagName)}, fields...))
	stdCore = CoreWrapper(tagName, stdCore)
	if xos.EnvGetCaseInsensitive("logbus_core_dup") != "" {
		stdCore = NewDupCore(stdCore)
	}
	cores = append(cores, stdCore)

	s := newStdLogger(gBasicZLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})), config, config.FetchLogContext)
	return s
}

// refreshEarlyLogger 用于刷新在Init() 之前创建的logger
func refreshEarlyLogger() {
	for _, v := range cacheGLogger {
		newStdLogger := newScopeLoggerInternal(v.tagName, v.printAsErr, v.inheritMode, v.fetch, v.fields...)
		v.gLogger.channelKey = newStdLogger.(*GLogger).channelKey
		v.gLogger.printAsError = newStdLogger.(*GLogger).printAsError
		_ = v.gLogger.stdLogger.Sync()
		v.gLogger.depthLogger.Range(func(key, _ interface{}) bool {
			v.gLogger.depthLogger.Delete(key)
			return true
		})
		v.gLogger.stdLogger = newStdLogger.(*GLogger).stdLogger
	}
	cacheGLogger = nil
}
