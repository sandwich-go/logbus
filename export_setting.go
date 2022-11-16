package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var EncodeConfig = zapcore.EncoderConfig{
	//CallerKey:      "caller",
	LevelKey:      LevelKey,
	MessageKey:    Meta, //zap's sampling algorithm uses the message to identify duplicate entries.
	TimeKey:       TimeKey,
	NameKey:       "logger",
	StacktraceKey: "stacktrace",
	LineEnding:    zapcore.DefaultLineEnding,
	EncodeLevel:   zapcore.LowercaseLevelEncoder,
	EncodeTime:    zapcore.ISO8601TimeEncoder,
	//EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller: zapcore.ShortCallerEncoder,
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

var ZapConf = zap.Config{
	Development:      false,
	Encoding:         "json",
	EncoderConfig:    EncodeConfig,
	OutputPaths:      []string{},
	ErrorOutputPaths: []string{},
}
