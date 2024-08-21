package logbus

import (
	"github.com/sandwich-go/logbus/thinkingdata"
	"os"
	"strconv"
	"time"
)

const (
	envTGATimeZoneOffset = "TGA_TIMEZONE_OFFSET"
)

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
	initBasics(conf)

	initGlobalStdLoggers()

	// set logger used in glog
	SetGlobalGLogger(gStdLogger, conf.DefaultChannel, conf.PrintAsError, 0)

	// init monitor
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

	if os.Getenv(envTGATimeZoneOffset) != "" {
		offset, ok := strconv.ParseInt(os.Getenv(envTGATimeZoneOffset), 10, 64)
		if ok == nil {
			thinkingdata.TgaLocation = time.FixedZone("", int(offset))
		}
	}
}

// Close 程序结束时打印缓存中的所有日志 并清理资源
func Close() {
	_ = gStdLogger.Sync()
	newGlobalGLogger.syncDepthLogger()
}

func resetLogBus() {
	//resetLoggerMap()
}
