package main

import (
	"github.com/sandwich-go/logbus"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// Init with conf
	logbus.Init(logbus.NewConf(logbus.WithDev(true), logbus.WithDefaultChannel("Simple"), logbus.WithCallerSkip(1)))

	// Print server debug log, dd_meta_channel=setting.DefaultChannel
	logbus.Debug("", logbus.Int("int", 123))
	// Print server debug log with specific channel, dd_meta_channel=logbus.THINKINGDATA
	logbus.DebugWithChannel(logbus.THINKINGDATA, "", logbus.Int("int", 123))

	// Print server info log, dd_meta_channel=setting.DefaultChannel
	logbus.Info("", logbus.Int("money", 648))
	// Print server info log with specific channel, dd_meta_channel=logbus.BI
	logbus.InfoWithChannel(logbus.BI, "", logbus.Int("money", 1296))

	// Print server warning log, dd_meta_channel=setting.DefaultChannel
	logbus.Warn("", logbus.String("str", "warning"))
	// Print server warning log with specific channel, dd_meta_channel=logbus.BIGQUERY
	logbus.WarnWithChannel(logbus.BIGQUERY, "", logbus.String("str", "warning"))

	// Print bi log, dd_meta_channel=bi
	//logbus.logger().L().Info("bi", zap.Int("money", 648)) // L()来输入dd_meta_channel的方式废弃
	logbus.InfoWithChannel(logbus.BI, "", logbus.Int("money", 648))
}
