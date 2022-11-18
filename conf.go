package logbus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate optionGen  --option_return_previous=false
func _ConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// log
		"LogLevel":       (zapcore.Level)(zap.DebugLevel), //@MethodComment(日志级别，默认 zap.DebugLevel)
		"Dev":            false,                           // false 输出json格式， true 则输出带颜色的易读log @MethodComment(是否输出带颜色的易读log，默认关闭)
		"DefaultChannel": string(SERVERLOG),               // 默认的dd_meta_channel @MethodComment(设置默认的dd_meta_channel)
		"DefaultTag":     string(DefaultTag),              // 默认打印的tag @MethodComment(设置默认的tag)
		"CallerSkip":     1,                               // zap logger callerSkip @MethodComment(等于zap.CallerSkip)
		//"LogId":         true,                            // 输出 log id @MethodComment(是否输出log_xid，默认开启) // 日志规范要求必须要有xid 不作为配置放出
		"StackLogLevel": (zapcore.Level)(zap.ErrorLevel), //@MethodComment(打印stack的最低级别，默认ErrorLevel stack if level >= StackLogLevel)
		// stdout
		"BufferedStdout": false, // @MethodComment(输出stdout时使用 logbus.BufferedWriteSyncer)

		// monitor
		"MonitorOutput": MonitorOutput(Noop), // [Logbus, Noop, Prometheus] @MethodComment(监控输出 Logbus, Noop, Prometheus)
		// The Prometheus metrics will be made available on this port: @MethodComment(prometheus监控输出端口，k8s集群保持默认9158端口)
		"DefaultPrometheusListenAddress": ":9158",
		// This is the endpoint where the Prometheus metrics will be made available ("/metrics" is the default with Prometheus):
		"DefaultPrometheusPath": "/metrics", // @MethodComment(prometheus监控输出接口path)
		// DefaultPercentiles is the default spread of percentiles/quantiles we maintain for timings / histogram metrics:
		"DefaultPercentiles":  []float64{0.5, 0.75, 0.99, 1},          //@MethodComment(监控统计耗时的分位值，默认统计耗时的 50%, 75%, 99%, 100% 的分位数)
		"DefaultLabel":        prometheus.Labels(map[string]string{}), //@MethodComment(监控额外添加的全局label，会在监控指标中显示)
		"MonitorTimingMaxAge": time.Duration(time.Minute),             // @MethodComment(monitor.Timing数据的最大生命周期)

		// glog
		"PrintAsError": true, //@MethodComment(glog输出field带error时，将日志级别提升到error)
	}
}

func init() {
	InstallConfWatchDog(func(cc *Conf) {
		if cc.DefaultLabel == nil {
			panic("DefaultLabel is nil")
		}
		if cc.MonitorOutput != Prometheus && cc.MonitorOutput != Logbus && cc.MonitorOutput != Noop {
			panic("MonitorOutput not match")
		}
	})
}
