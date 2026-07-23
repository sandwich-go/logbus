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
//
// 注意：当环境变量 sys_conf_path_env 存在时（PMT 部署环境），Init 末尾会启用动态日志级别：
// 读取 $sys_conf_path_env/ops_config.json 并 watch，json 中的 log_level 会覆盖 conf.LogLevel。
// 即在 PMT 环境下，业务代码里 WithLogLevel(...) 传入的值会被 yaml 覆盖。本地开发环境
// （未设置 sys_conf_path_env）不受影响，行为不变。
//
// Init 应在启动阶段调用，不支持与 Init 或 NewScopeLogger 并发执行。
// LogLevel 和 EnableTraceLevel 是进程级运行时开关，不随单个 Logger 隔离。
func Init(conf *Conf) {
	config := snapshotConf(conf)

	// 环境变量控制 EnableTraceLevel
	if v := xos.EnvGetCaseInsensitive("logbus_enable_trace_level"); v != "" {
		config.EnableTraceLevel = v == "1" || v == "true"
	}
	// 本地启动多个服务时可以方便地屏蔽 monitor。
	if xos.EnvGetCaseInsensitive("logbus_disable_monitor") != "" {
		config.MonitorOutput = Noop
	}
	ignoreLogicalError.Store(config.IgnoreLogicalError)
	config.prepareWriter()
	initBasics(config)
	activeConfig.Store(config)
	Setting = config.detached()

	initGlobalStdLoggers(config)

	// set logger used in glog
	SetGlobalGLogger(gStdLogger.Load(), config.DefaultChannel, config.PrintAsError, 0)

	// init monitor
	monitorConfig := cloneConf(&config.Conf)
	setDefaultMetricsReporter(monitorConfig.MonitorOutput,
		monitorConfig.DefaultPrometheusListenAddress,
		monitorConfig.DefaultPrometheusPath,
		monitorConfig.DefaultPercentiles,
		monitorConfig.DefaultLabel,
		monitorConfig.MonitorTimingMaxAge)

	if refresh {
		refresh = false
		refreshEarlyLogger()
	}

	// 自动启用 PMT 下发的动态日志级别：若 sys_conf_path_env 环境变量不存在则静默跳过。
	enableDynamicLogLevel()
}

func resetLogBus() {
	//resetLoggerMap()
}
