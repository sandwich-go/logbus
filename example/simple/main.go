package main

import (
	"github.com/sandwich-go/logbus"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// 主线程中使用 非线程安全
	logbus.Init(logbus.NewConf(logbus.WithDev(true), logbus.WithDefaultChannel("Simple")))

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
	logbus.InfoWithChannel(logbus.BI, "", logbus.Int("money", 648))
}
