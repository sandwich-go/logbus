package basics

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/config"
	"bitbucket.org/funplus/sandwich/pkg/logbus/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var BasicLogger *zap.Logger
var BasicSugarLogger *zap.SugaredLogger

var Setting = config.NewDefaultConf()

func InitBasic(c *config.Conf) {
	ResetLogBus()
	var err error
	Setting = c
	config.ZapConf.Level = zap.NewAtomicLevelAt(c.LogLevel)
	if c.Dev {
		config.EncodeConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		config.EncodeConfig.CallerKey = "caller"
		config.EncodeConfig.EncodeDuration = zapcore.StringDurationEncoder
		config.ZapConf.Development = true
	} else {
		config.EncodeConfig.EncodeDuration = config.DurationEncoder
	}
	BasicLogger, err = config.ZapConf.Build(
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddStacktrace(c.StackLogLevel),
		zap.WithClock(&localClock{}),
	)
	if err != nil {
		panic(err)
	}
	BasicSugarLogger = BasicLogger.Sugar()
}

func ResetLogBus() {
	global.LoggerMap.Range(func(key, value interface{}) bool {
		global.LoggerMap.Delete(key)
		return true
	})
}
