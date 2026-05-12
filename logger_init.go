package logbus

import "github.com/sandwich-go/boost/xos"

// refresh 用于刷新在Init() 之前创建的logger
var refresh = true

func init() {
	refresh = false // 默认不执行 业务调用时执行且只执行一次
	Init(NewConf())
	refresh = true
}

// Init logBus初始化 会有两次调用。一次是init()，一次是手动调用Init的时候
// 允许不手动Init的情况下使用默认配置调用logBus
func Init(conf *Conf) {
	// 环境变量控制 EnableTraceLevel
	if v := xos.EnvGetCaseInsensitive("logbus_enable_trace_level"); v != "" {
		conf.EnableTraceLevel = v == "1" || v == "true"
	}
	initBasics(conf)

	// 必须先管理 track file writer 生命周期，再创建 logger 实例，
	// 这样 newNLoggerInstance 才能引用到最新的全局单例 gTrackFileWriteSyncer。
	initTrackFileWriteSyncer()

	initGlobalStdLoggers()

	// set logger used in glog
	SetGlobalGLogger(gStdLogger, conf.DefaultChannel, conf.PrintAsError, 0)

	// init monitor
	// 本地启动多个服务时可以方便的屏蔽monitor
	if xos.EnvGetCaseInsensitive("logbus_disable_monitor") != "" {
		conf.MonitorOutput = Noop
	}
	setDefaultMetricsReporter(conf.MonitorOutput,
		conf.DefaultPrometheusListenAddress,
		conf.DefaultPrometheusPath,
		conf.DefaultPercentiles,
		conf.DefaultLabel,
		conf.MonitorTimingMaxAge)

	if refresh {
		refresh = false
		refreshEarlyLogger()
	}
}

func resetLogBus() {
	//resetLoggerMap()
}
