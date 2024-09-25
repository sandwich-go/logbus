package logbus

import (
	"github.com/sandwich-go/boost/xpanic"
	"go.uber.org/zap/buffer"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var EncodeConfig zapcore.EncoderConfig
var ZapConf zap.Config
var BufferPool buffer.Pool
var modulePath string

func init() {
	initZapSetting()
	BufferPool = buffer.NewPool()
	modulePath = func() string {
		_, currentFile, _, _ := runtime.Caller(0)
		currentDir := filepath.Dir(currentFile)
		modulePath, err := findModuleRoot(currentDir)
		xpanic.WhenError(err)
		return modulePath
	}()
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
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	ZapConf = zap.Config{
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    EncodeConfig,
		OutputPaths:      []string{},
		ErrorOutputPaths: []string{},
	}

}

func CustomCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(func(ec zapcore.EntryCaller) string {
		buf := BufferPool.Get()
		buf.AppendString(getCallerPath(ec.File))
		buf.AppendByte(':')
		buf.AppendInt(int64(ec.Line))
		caller := buf.String()
		buf.Free()
		return caller
	}(caller))
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
