package stdl

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/basics"
	"bitbucket.org/funplus/sandwich/pkg/logbus/config"
	"bitbucket.org/funplus/sandwich/pkg/logbus/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sort"
	"strings"
)

// getLogger get StdLogger by tagKey. If StdLogger doesn't exist, add a new one
// to be changed to private function: getLogger(...string)
func getLogger(tags ...string) *StdLogger {
	if len(tags) == 0 {
		tags = []string{config.DefaultTag}
	}
	for k, v := range tags {
		tags[k] = strings.TrimSpace(v)
	}

	tagKey := getTagKey(tags)
	if logger, ok := global.LoggerMap.Load(tagKey); ok {
		return logger.(*StdLogger)
	}

	var cores []zapcore.Core
	encoder := basics.NewJSONEncoder(config.EncodeConfig)
	if basics.Setting.Dev {
		encoder = basics.NewConsoleEncoder(config.EncodeConfig)
	}
	if basics.Setting.OutputStdout {
		var writer zapcore.WriteSyncer
		writer = os.Stdout
		if basics.Setting.BufferedStdout {
			writer = config.BufferedWriteSyncer
		}
		stdCore := zapcore.NewCore(encoder, writer, basics.Setting.LogLevel).With([]zap.Field{zap.Strings(config.Tags, tags)})
		cores = append(cores, stdCore)
	}
	if basics.Setting.OutputLocalFile {
		var writers []zapcore.WriteSyncer
		for _, tag := range tags {
			fileWriter := basics.NewFileWriter(tag)
			writers = append(writers, zapcore.AddSync(fileWriter))
		}
		localFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), basics.Setting.LogLevel)
		cores = append(cores, localFileCore)
	}
	if basics.Setting.OutputFluentd {
		fluentdCore := basics.NewFluentdCore(config.EncodeConfig, tags, basics.Setting.LogLevel)
		cores = append(cores, fluentdCore)
	}

	newCore := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
	newZapLogger := basics.BasicLogger.WithOptions(newCore)
	logger := NewDefaultStdLogger(newZapLogger, tags)

	global.LoggerMap.Store(tagKey, logger)

	return logger
}

func getTagKey(tags []string) string {
	sort.Strings(tags)
	return strings.Join(tags, ".")
}
