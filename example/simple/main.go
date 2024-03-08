package main

import (
	"github.com/sandwich-go/logbus"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// 主线程中使用 非线程安全
	logbus.EncodeConfig.CallerKey = "caller"
	logbus.Init(logbus.NewConf(logbus.WithDev(false), logbus.WithDefaultChannel("Simple")))

	// Print server debug log, dd_meta_channel=setting.DefaultChannel
	logbus.Debug("", logbus.Int("int", 123))

	// Print server info log, dd_meta_channel=setting.DefaultChannel
	logbus.Info("", logbus.Int("money", 648))

	// User defined channel, dd_meta_channel=setting.UserDefine
	logbus.InfoWithChannel("UserDefine", "", logbus.Strings("str1", []string{"hello", "world"}))
}
