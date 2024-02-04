package logbus

import (
	"time"

	prometheusClient "github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/logbus/monitor"
	"github.com/sandwich-go/logbus/monitor/noop"
	"github.com/sandwich-go/logbus/monitor/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gBasicZLogger *zap.Logger        // 用来产生新stdLogger, 更多用在自定义tag的时候
	Setting       = newDefaultConf() // logBus的全局配置
)

func initBasics(c *Conf) {
	resetLogBus()
	// init logBus global setting
	Setting = c

	// init EncodeConfig
	if c.Dev {
		EncodeConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		EncodeConfig.CallerKey = "caller"
		EncodeConfig.EncodeDuration = zapcore.StringDurationEncoder
	} else {
		EncodeConfig.EncodeDuration = DurationEncoder
	}

	// init gBasicZLogger
	var err error
	ZapConf.Level = zap.NewAtomicLevelAt(c.LogLevel)
	ZapConf.EncoderConfig = EncodeConfig
	if c.Dev {
		ZapConf.Development = true
	}
	gBasicZLogger, err = ZapConf.Build(
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddStacktrace(c.StackLogLevel),
		zap.WithClock(&localClock{}),
		zap.WithCaller(ZapConf.EncoderConfig.CallerKey != ""),
	)
	if err != nil {
		panic(err)
	}
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
		DebugWithChannel(Monitor, "", String("prometheus [http] listening on", defaultPrometheusListenAddress), String("path", defaultPrometheusPath))
	}
}
