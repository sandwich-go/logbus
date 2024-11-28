package logbus

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var EncodeConfig zapcore.EncoderConfig
var ZapConf zap.Config

func init() {
	initZapSetting()
}

// initZapSetting 更新配置数据
func initZapSetting() {
	EncodeConfig = zapcore.EncoderConfig{
		//CallerKey:      "caller",
		LevelKey:      LevelKey,
		MessageKey:    Meta, //zap's sampling algorithm uses the message to identify duplicate entries.
		TimeKey:       TimeKey,
		NameKey:       "logger",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   traceLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	ZapConf = zap.Config{
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    EncodeConfig,
		OutputPaths:      []string{},
		ErrorOutputPaths: []string{},
	}

}

// DurationEncoder serializes a time.Duration to a floating-point number of milliseconds elapsed.
var DurationEncoder = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

var BufferedWriteSyncer = &zapcore.BufferedWriteSyncer{
	WS:            os.Stdout,
	Size:          256 * 1024, // 256 kB
	FlushInterval: 30 * time.Second,
}
