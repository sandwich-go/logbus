package main

import (
	"context"

	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/logbus"
)

func main() {
	ctx := context.Background()

	// close logger before exit
	defer logbus.Close()

	// 主线程中使用 非线程安全
	logbus.EncodeConfig.CallerKey = "caller"
	logbus.Init(logbus.NewConf(
		logbus.WithDev(false),
		logbus.WithDefaultChannel("Simple"),
		logbus.WithPrintAsError(true),       // 检测到field里有errorType，则把日志级别提升到error
		logbus.WithIgnoreLogicalError(true), // 逻辑错误则不提升到error级别
	))

	// Print server debug log, dd_meta_channel=setting.DefaultChannel
	logbus.Debug(ctx, "", logbus.Int("int", 123))

	// Print server info log, dd_meta_channel=setting.DefaultChannel
	logbus.Info(ctx, "", logbus.Int("money", 648))

	// User defined channel, dd_meta_channel=setting.UserDefine
	logbus.InfoWithChannel(ctx, "UserDefine", "", logbus.Strings("str1", []string{"hello", "world"}))

	// 错误（非逻辑）自动升级为error级别日志
	logbus.Warn(ctx, "自动升级为error日志", logbus.E(xerror.NewText("some error")))
	// 逻辑错误不升级为error级别日志
	logbus.Warn(ctx, "逻辑错误不升级为error日志", logbus.E(xerror.NewText("some logical error").SetLogic()))
}
