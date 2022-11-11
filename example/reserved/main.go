package main

import (
	"github.com/sandwich-go/logbus"
	"github.com/sandwich-go/logbus/config"
	"go.uber.org/zap"
)

// 使用"自定义tag"的用法 作为保留
// 以下情况下使用：写入指定文件、使用fluentd socket做日志收集
func main() {
	// close logger before exit
	defer logbus.Close()

	// Init with conf
	logbus.Init(config.NewConf(config.WithOutputStdout(true), config.WithDev(false), config.WithMonitorOutput(logbus.Prometheus), config.WithDefaultChannel("Game")))

	// reason: 不建议使用tag
	// tag=diandian.market  dd_meta_channel=logbus.BI
	//logbus.Logger("diandian.market", "server").L().Info(logbus.BI, zap.Int("dau", 999)) // Deprecated 不兼容升级
	logbus.Logger("diandian.market", "server").InfoWithChannel(config.BI, zap.Int("dau", 9992))

	// reason: 减少暴露tag给使用者，要求tag的时候明确使用LoggerWithTags 这里自定义tag的方式不建议使用
	// tag=payment  dd_meta_channel=setting.DefaultChannel
	logbus.Logger("payment").Warn(zap.Int("money", 6481))

	// 如果要指定channel，必须使用InfoWithChannel
	// tag=[diandian.market, server]  dd_meta_channel=logbus.BI
	//logbus.Logger("diandian.market", "server").L().Info(logbus.BI, zap.Int("dau", 999)) // Deprecated 不兼容升级
	logbus.Logger("diandian.market", "server").InfoWithChannel(config.BI, zap.Int("dau", 999))
}
