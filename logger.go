package logbus

func init() {
	Init(NewConf())
}

// Init logBus初始化 会有两次调用。一次是init()，一次是手动调用Init的时候
// 允许不手动Init的情况下使用默认配置调用logBus
func Init(conf *Conf) {
	initBasics(conf)

	initGlobalStdLoggers()

	// set logger used in glog
	SetGlobalGLogger(gStdLogger, conf.DefaultChannel, conf.PrintAsError, conf.CallerSkip)
	// init monitor
	setDefaultMetricsReporter(conf.MonitorOutput,
		conf.DefaultPrometheusListenAddress,
		conf.DefaultPrometheusPath,
		conf.DefaultPercentiles,
		conf.DefaultLabel,
		conf.MonitorTimingMaxAge)
}

// Close 程序结束时打印缓存中的所有日志 并清理资源
func Close() {
	_ = gStdLogger.Sync()
	newGlobalGLogger.syncDepthLogger()
}

func resetLogBus() {
	//resetLoggerMap()
}
