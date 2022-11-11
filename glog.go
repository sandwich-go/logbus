package logbus

import (
	"github.com/sandwich-go/logbus/basics"
	"github.com/sandwich-go/logbus/glog"
	"github.com/sandwich-go/logbus/stdl"
	"go.uber.org/zap"
)

func ImplementGLog(printAsError bool) {
	setGlobalLogger(nil, "", printAsError)
}

// SetGlobalLogger 改变glog使用的globalLogger
// logger 为自定义的 logbus.Logger(), 传nil则使用默认的
// channelKey MessageKey对应的值，默认 "dd_meta_channel":"server"
// printAsError true 检测到field里有errorType，则把日志级别提升到error
func SetGlobalLogger(logger *stdl.StdLogger, channelKey string, printAsError bool) {
	setGlobalLogger(logger, channelKey, printAsError)
}

func setGlobalLogger(logger *stdl.StdLogger, channelKey string, printAsError bool) {
	if logger == nil {
		logger = GetStdLogger()
	}
	if channelKey == "" {
		channelKey = basics.Setting.DefaultChannel
	}
	cloneLogger := logger.CloneLogger(zap.AddCallerSkip(1))

	glog.SetGlobalGLogger(NewGLogger(channelKey, cloneLogger, printAsError))
}
