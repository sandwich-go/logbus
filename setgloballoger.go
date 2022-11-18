package logbus

import (
	"go.uber.org/zap"
)

// SetGlobalGLogger 改变glog使用的globalLogger. oldName used in docs: SetGlogLogger
// logger 为自定义的 logbus.logger(), 传nil则使用默认的
// channelKey MessageKey对应的值，默认 "dd_meta_channel":"server"
// printAsError true 检测到field里有errorType，则把日志级别提升到error
func SetGlobalGLogger(logger *StdLogger, channelKey string, printAsError bool, callerSkip int) {
	if logger == nil {
		logger = gStdLogger
	}
	if channelKey == "" {
		channelKey = Setting.DefaultChannel
	}

	if callerSkip <= 0 {
		callerSkip = Setting.CallerSkip
	}

	newGlobalGLogger = &GLogger{
		channelKey: channelKey,
		stdLogger: &StdLogger{
			z: logger.WithOptions(zap.AddCallerSkip(callerSkip)),
		},
		printAsError: printAsError,
	}
}
