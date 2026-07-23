package logbus

import (
	"context"
	"time"

	prometheusClient "github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/logbus/debug"
	"github.com/sandwich-go/logbus/monitor"
	"github.com/sandwich-go/logbus/monitor/noop"
	"github.com/sandwich-go/logbus/monitor/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gBasicZLogger *zap.Logger // 用来产生新stdLogger, 更多用在自定义tag的时候

	// Setting 保留用于源码兼容，仅表示最近一次 Init 的配置副本。
	// 替换它的字段不会影响已经创建或之后创建的 Logger；函数闭包和 Writer 等依赖仍保持其对象身份。
	//
	// Deprecated: 使用 NewConf 和 Init 配置 logbus；运行时日志级别使用 SetLogLevel/GetLogLevel。
	Setting = newDefaultConf()
)

func initBasics(c *configSnapshot) {
	initZapSetting()
	resetLogBus()
	runtimeLogLevel.SetLevel(c.LogLevel)
	runtimeEnableTrackLevel.Store(c.EnableTraceLevel)

	// init EncodeConfig
	if c.Dev {
		EncodeConfig.EncodeLevel = trackLevelEncoderWithColor
		EncodeConfig.CallerKey = "caller"
		EncodeConfig.EncodeDuration = zapcore.StringDurationEncoder
		if err := debug.InitAppModuleFilePath(); err == nil {
			EncodeConfig.EncodeCaller = debug.RelativePathCallerEncoder
		} else {
			// fixme 无法使用相对路径的caller模式，退化到short模式
			EncodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
		}
	} else {
		EncodeConfig.EncodeDuration = DurationEncoder
	}

	if c.EncodeCaller != nil {
		EncodeConfig.EncodeCaller = c.EncodeCaller
	}
	c.encoder = EncodeConfig

	// init gBasicZLogger
	var err error
	zapConf := ZapConf
	zapConf.Level = runtimeLogLevel
	zapConf.EncoderConfig = c.encoder
	if c.Dev {
		zapConf.Development = true
	}
	var clock zapcore.Clock
	clock = localClock{}
	if c.UseSystemClock {
		clock = zapcore.DefaultClock
	}
	gBasicZLogger, err = zapConf.Build(
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddStacktrace(c.StackLogLevel),
		zap.WithClock(clock),
		zap.WithCaller(zapConf.EncoderConfig.CallerKey != ""),
	)
	if err != nil {
		panic(err)
	}
	ZapConf = zapConf
}

func setDefaultMetricsReporter(
	monitorOutput MonitorOutput,
	defaultPrometheusListenAddress string,
	defaultPrometheusPath string,
	defaultPercentiles []float64,
	defaultLabel prometheusClient.Labels,
	timingMaxAge time.Duration) {
	switch monitorOutput {
	case Noop:
		monitor.DefaultMetricsReporter = noop.New()
	case Logbus:
		monitor.DefaultMetricsReporter = newLogReporter()
	case Prometheus:
		var err error
		monitor.DefaultMetricsReporter, err = prometheus.New(defaultPrometheusListenAddress, defaultPrometheusPath, defaultPercentiles, defaultLabel, timingMaxAge)
		if err != nil {
			panic(err)
		}
		DebugWithChannel(context.Background(), Monitor, "", String("prometheus [http] listening on", defaultPrometheusListenAddress), String("path", defaultPrometheusPath))
	}
}
